package trie

import (
	"fmt"
	"github.com/FrancisMcN/lib-mmdb/field"
	"github.com/FrancisMcN/lib-mmdb/node"
	"math/big"
	"net"
)

type Trie struct {
	totalId     **big.Int
	Root        *node.Node
	dataMap     map[string]int
	data        []byte
	recordSize  int
	Size        uint32
	ShouldPrune bool
}

func NewTrie() *Trie {
	id := big.NewInt(0)
	return &Trie{
		totalId:     &id,
		Root:        node.NewNode(),
		dataMap:     make(map[string]int),
		data:        make([]byte, 0),
		recordSize:  28,
		ShouldPrune: true,
	}
}

func (t *Trie) Insert(cidr *net.IPNet, data field.Field) {

	currentNode := t.Root
	ones, bits := cidr.Mask.Size()

	if bits == 32 {
		ones += 96
	}

	var foundExisting bool
	var existingData field.Field

	for i := 0; i < ones; i++ {

		if currentNode.Data() != nil {
			foundExisting = true
			existingData = currentNode.Data()
			currentNode.SetData(nil)
		}

		if !isSet(cidr.IP, i) {
			left := currentNode.Left
			if left == nil {
				currentNode.SetLeft(node.NewNode())
			}

			if foundExisting {
				if currentNode.Right == nil {
					currentNode.Right = node.NewNode()
				}
				currentNode.Right.SetData(existingData)
			}

			currentNode = currentNode.Left

		} else {

			right := currentNode.Right
			if right == nil {
				currentNode.SetRight(node.NewNode())
			}

			if foundExisting {
				if currentNode.Left == nil {
					currentNode.Left = node.NewNode()
				}
				currentNode.Left.SetData(existingData)
			}

			currentNode = currentNode.Right

		}
	}

	data = t.AddData(data)

	id := big.NewInt(int64(uint32(t.dataMap[fmt.Sprintf("%x", data.String())])))
	currentNode.SetData(data)
	currentNode.SetId(&id)

}

func (t *Trie) AddData(data field.Field) field.Field {

	if data.Type() == field.MapField {
		data = t.PointerifyMap(data.(*field.Map))
	} else if data.Type() == field.ArrayField {
		data = t.PointerifyArray(data.(field.Array))
	} else if _, f := t.dataMap[fmt.Sprintf("%x", data)]; !f {
		t.dataMap[fmt.Sprintf("%x", data)] = len(t.data)
		t.data = append(t.data, data.Bytes()...)
	}

	return data
}

func (t *Trie) Finalise() {
	nid := int64(0)
	if t.ShouldPrune {
		t._finalise(&t.Root, &nid)
	}
	t._finalise2(t.Root, &nid)
	t._finalise3(t.Root, nid)
	(*t.totalId).Set(big.NewInt(nid))
	t.Size = uint32(nid)
}

func (t *Trie) _finalise(parent **node.Node, nid *int64) {

	n := *parent

	if n != nil {

		left := n.Left
		right := n.Right
		// Prune where two child nodes are the same
		if left != nil && right != nil && left.Data() != nil && right.Data() != nil && left.Data().String() == right.Data().String() {
			*parent = node.NewNode()
			(*parent).SetData(left.Data())
			(*parent).SetId(t.totalId)
			return
		}

		if n.Left != nil {
			t._finalise(&n.Left, nid)
		}
		if n.Right != nil {
			t._finalise(&n.Right, nid)
		}
	}
}

func (t *Trie) _finalise2(n *node.Node, nid *int64) {

	if n != nil {

		if n.Data() == nil {
			id := big.NewInt(*nid)
			n.SetId(&id)
		} else {
			d := fmt.Sprintf("%x", n.Data())
			dataOffset, _ := t.dataMap[d]
			id := big.NewInt(int64(uint32(dataOffset)+16) * -1)
			n.SetId(&id)
		}
		if n.Left != nil || n.Right != nil {
			*nid++

			t._finalise2(n.Left, nid)
			t._finalise2(n.Right, nid)
		}
	}
}

func (t *Trie) _finalise3(n *node.Node, nid int64) {

	// Adds the total node count to the nodes that point into the data section
	if n != nil {

		if (*n.Id()).Cmp(big.NewInt(0)) < 0 {
			newId := big.NewInt(0).Mul(big.NewInt(-1), *n.Id())
			newId.Add(newId, big.NewInt(nid))
			n.SetId(&newId)
		}

		t._finalise3(n.Left, nid)
		t._finalise3(n.Right, nid)
	}
}

func (t *Trie) SetTotalId(id *big.Int) {
	*t.totalId = id
}

func (t *Trie) Print() {
	t._print(t.Root)
}

func (t *Trie) _print(n *node.Node) {
	if n != nil {

		fmt.Println("n", n)
		t._print(n.Left)
		t._print(n.Right)
	}
}

func (t *Trie) Serialise(n *node.Node, bytes *[]byte) {

	if n == nil {
		return
	}

	if n.Left == nil && n.Right == nil {
		return
	}

	*bytes = append(*bytes, n.Bytes(t.recordSize, (*t.totalId).Uint64())...)
	t.Serialise(n.Left, bytes)
	t.Serialise(n.Right, bytes)
}

func (t Trie) Bytes() []byte {
	bytes := make([]byte, 0)

	t.Serialise(t.Root, &bytes)

	bytes = append(bytes, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)
	bytes = append(bytes, t.data...)
	bytes = append(bytes, 0xAB, 0xCD, 0xEF, 'M', 'a', 'x', 'M', 'i', 'n', 'd', '.', 'c', 'o', 'm')

	return bytes
}

func (t *Trie) PointerifyMap(m *field.Map) *field.Map {

	m2 := field.NewMap()
	offset := 1

	for _, k := range m.OrderedKeys {
		var kf field.Field
		if key, f := t.dataMap[fmt.Sprintf("%x", k.Bytes())]; f {
			kf = field.Pointer(key)
		} else {
			kf = k
			t.dataMap[fmt.Sprintf("%x", k.Bytes())] = len(t.data) + offset
		}

		offset += len(kf.Bytes())

		v := m.InternalMap[k]
		var vf field.Field
		if val, f := t.dataMap[fmt.Sprintf("%x", v.Bytes())]; f {
			vf = field.Pointer(val)
		} else {
			vf = v
			t.dataMap[fmt.Sprintf("%x", v.Bytes())] = len(t.data) + offset
		}
		offset += len(vf.Bytes())

		m2.Put(kf, vf)

	}

	t.dataMap[fmt.Sprintf("%x", m2)] = len(t.data)
	t.data = append(t.data, m2.Bytes()...)

	return m2
}

func (t *Trie) PointerifyArray(a []field.Field) field.Array {
	arr := make([]field.Field, 0)

	for _, v := range a {

		var vf field.Field
		if _, f := t.dataMap[fmt.Sprintf("%x", v.Bytes())]; !f {
			t.dataMap[fmt.Sprintf("%x", v.Bytes())] = len(t.data)
			t.data = append(t.data, v.Bytes()...)
		}
		pointerVal := t.dataMap[fmt.Sprintf("%x", v.Bytes())]
		v = field.Pointer(pointerVal)
		arr = append(arr, vf)
	}
	t.dataMap[fmt.Sprintf("%x", field.Array(arr))] = len(t.data)
	t.data = append(t.data, field.Array(arr).Bytes()...)
	return arr
}

func (t *Trie) GetDataMap() map[string]int {
	return t.dataMap
}

// Determines if the 'bit' in the IP is set
// 'bit' is calculated from the most significant byte first
func isSet(ip net.IP, bit int) bool {
	whichByte := bit / 8
	ipByte := ip[whichByte]
	return ((ipByte >> (7 - (bit % 8))) & 1) > 0
}
