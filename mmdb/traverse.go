package mmdb

import (
	"errors"
	"fmt"
	"github.com/FrancisMcN/lib-mmdb/field"
	"github.com/FrancisMcN/lib-mmdb/node"
	"net"
)

// Internal structure used to keep track of nodes we still need to visit.
type netNode struct {
	ip      net.IP
	bit     uint
	pointer uint
}

// Networks represents a set of subnets that we are iterating over.
type Networks struct {
	//reader   *Reader
	mmdb     *MMDB
	nodes    []netNode // Nodes we still have to visit.
	lastNode netNode
	err      error

	skipAliasedNetworks bool
}

var (
	allIPv4 = &net.IPNet{IP: make(net.IP, 4), Mask: net.CIDRMask(0, 32)}
	allIPv6 = &net.IPNet{IP: make(net.IP, 16), Mask: net.CIDRMask(0, 128)}
)

// NetworksOption are options for Networks and NetworksWithin.
type NetworksOption func(*Networks)

// SkipAliasedNetworks is an option for Networks and NetworksWithin that
// makes them not iterate over aliases of the IPv4 subtree in an IPv6
// database, e.g., ::ffff:0:0/96, 2001::/32, and 2002::/16.
//
// You most likely want to set this. The only reason it isn't the default
// behavior is to provide backwards compatibility to existing users.
func SkipAliasedNetworks(networks *Networks) {
	networks.skipAliasedNetworks = true
}

// Networks returns an iterator that can be used to traverse all networks in
// the database.
//
// Please note that a MaxMind DB may map IPv4 networks into several locations
// in an IPv6 database. This iterator will iterate over all of these locations
// separately. To only iterate over the IPv4 networks once, use the
// SkipAliasedNetworks option.
func (m *MMDB) Networks() *Networks {
	var networks *Networks
	if m.metadata.IpVersion == 6 {
		networks = m.NetworksWithin(allIPv6)
	} else {
		networks = m.NetworksWithin(allIPv4)
	}
	return networks
}

func (m *MMDB) traverseTree(ip net.IP, nd, bitCount uint) (uint, int) {
	nodeCount := m.metadata.NodeCount
	recordSize := m.metadata.RecordSize
	recordBytes := recordSize / 8
	nodeBytes := recordBytes * 2
	if recordSize%8 > 0 {
		nodeBytes++
	}

	i := uint(0)

	for ; i < bitCount && uint32(nd) < nodeCount; i++ {
		bit := uint(1) & (uint(ip[i>>3]) >> (7 - (i % 8)))
		offset := nd * uint(nodeBytes)
		n := node.FromBytes(m.Bst[offset:uint32(offset)+uint32(nodeBytes)], recordSize)
		if bit == 0 {
			nd = uint(n[0].Uint64())
		} else {
			nd = uint(n[1].Uint64())
		}
	}

	return nd, int(i)
}

// NetworksWithin returns an iterator that can be used to traverse all networks
// in the database which are contained in a given network.
//
// Please note that a MaxMind DB may map IPv4 networks into several locations
// in an IPv6 database. This iterator will iterate over all of these locations
// separately. To only iterate over the IPv4 networks once, use the
// SkipAliasedNetworks option.
//
// If the provided network is contained within a network in the database, the
// iterator will iterate over exactly one network, the containing network.
func (m *MMDB) NetworksWithin(network *net.IPNet) *Networks {
	if m.metadata.IpVersion == 4 && network.IP.To4() == nil {
		return &Networks{
			err: fmt.Errorf(
				"error getting networks with '%s': you attempted to use an IPv6 network in an IPv4-only database",
				network.String(),
			),
		}
	}

	networks := &Networks{mmdb: m}

	ip := network.IP
	prefixLength, _ := network.Mask.Size()

	if m.metadata.IpVersion == 6 && len(ip) == net.IPv4len {
		if networks.skipAliasedNetworks {
			ip = net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, ip[0], ip[1], ip[2], ip[3]}
		} else {
			ip = ip.To16()
		}
		prefixLength += 96
	}

	pointer, bit := m.traverseTree(ip, 0, uint(prefixLength))
	networks.nodes = []netNode{
		{
			ip:      ip,
			bit:     uint(bit),
			pointer: pointer,
		},
	}
	return networks
}

// Next prepares the next network for reading with the Network method. It
// returns true if there is another network to be processed and false if there
// are no more networks or if there is an error.
func (n *Networks) Next() bool {
	if n.err != nil {
		return false
	}

	nodeCount := n.mmdb.metadata.NodeCount
	recordSize := n.mmdb.metadata.RecordSize
	recordBytes := recordSize / 8
	nodeBytes := recordBytes * 2
	if recordSize%8 > 0 {
		nodeBytes++
	}

	for len(n.nodes) > 0 {
		nd := n.nodes[len(n.nodes)-1]
		n.nodes = n.nodes[:len(n.nodes)-1]
		for uint32(nd.pointer) != nodeCount {

			if uint32(nd.pointer) > nodeCount {
				n.lastNode = nd
				return true
			}
			ipRight := make(net.IP, len(nd.ip))
			copy(ipRight, nd.ip)
			if len(ipRight) <= int(nd.bit>>3) {
				n.err = errors.New(
					fmt.Sprintf("invalid search tree at %v/%v", ipRight, nd.bit))
				return false
			}
			ipRight[nd.bit>>3] |= 1 << (7 - (nd.bit % 8))
			offset := nd.pointer * uint(nodeBytes)
			nodes := node.FromBytes(n.mmdb.Bst[offset:uint32(offset)+uint32(nodeBytes)], recordSize)
			rightPointer := uint(nodes[1].Uint64())
			nd.bit++
			n.nodes = append(n.nodes, netNode{
				pointer: rightPointer,
				ip:      ipRight,
				bit:     nd.bit,
			})

			nd.pointer = uint(nodes[0].Uint64())
		}
	}

	return false
}

// Network returns the current network or an error if there is a problem
// decoding the data for the network. It takes a pointer to a result value to
// decode the network's data into.
func (n *Networks) Network() (*net.IPNet, field.Field, error) {
	if n.err != nil {
		return nil, nil, n.err
	}
	var result field.Field
	fp := field.FieldParserSingleton()
	//dataOffset := n[0].Sub(n[0], big.NewInt(int64(nodeCount)))
	//dataOffset = dataOffset.Sub(dataOffset, big.NewInt(16))
	//fp.SetOffset(uint32(dataOffset.Uint64()))

	dataOffset := uint32(n.lastNode.pointer) - n.mmdb.metadata.NodeCount - 16
	fp.SetOffset(dataOffset)
	result = fp.Parse(n.mmdb.Data)

	//if err := n.reader.retrieveData(n.lastNode.pointer, result); err != nil {
	//	return nil, err
	//}

	ip := n.lastNode.ip
	prefixLength := int(n.lastNode.bit)

	//// We do this because uses of SkipAliasedNetworks expect the IPv4 networks
	//// to be returned as IPv4 networks. If we are not skipping aliased
	//// networks, then the user will get IPv4 networks from the ::FFFF:0:0/96
	//// network as Go automatically converts those.
	//if n.skipAliasedNetworks && isInIPv4Subtree(ip) {
	//	ip = ip[12:]
	//	prefixLength -= 96
	//}

	return &net.IPNet{
		IP:   ip,
		Mask: net.CIDRMask(prefixLength, len(ip)*8),
	}, result, nil
}

// Err returns an error, if any, that was encountered during iteration.
func (n *Networks) Err() error {
	return n.err
}

// isInIPv4Subtree returns true if the IP is an IPv6 address in the database's
// IPv4 subtree.
func isInIPv4Subtree(ip net.IP) bool {
	if len(ip) != 16 {
		return false
	}
	for i := 0; i < 12; i++ {
		if ip[i] != 0 {
			return false
		}
	}
	return true
}
