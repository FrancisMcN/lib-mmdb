package trie

import (
	"fmt"
	"github.com/FrancisMcN/lib-mmdb2/field"
	"github.com/FrancisMcN/lib-mmdb2/node"
	"math/big"
	"net"
)

type Trie struct {
	totalId    **big.Int
	root       *node.Node
	dataMap    map[string]int
	data       []byte
	recordSize int
	Size uint32
	ShouldPrune bool
}

func NewTrie() *Trie {
	id := big.NewInt(0)
	return &Trie{
		totalId: &id,
		root:    node.NewNode(),
		dataMap: make(map[string]int),
		data: make([]byte, 0),
		recordSize: 28,
		ShouldPrune: true,
	}
}

func (t *Trie) Insert(cidr *net.IPNet, data field.Field) {

	currentNode := t.root
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
			//left := currentNode.Children()[0]
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
	if _, f := t.dataMap[data.String()]; !f {
		t.dataMap[fmt.Sprintf("%x", data)] = len(t.data)
		t.data = append(t.data, data.Bytes()...)
	}
	id := big.NewInt(int64(uint32(t.dataMap[data.String()])))
	currentNode.SetData(data)
	currentNode.SetId(&id)

}

func (t *Trie) Finalise() {
	nid := int64(0)
	if t.ShouldPrune {
		t._finalise(&t.root, &nid)
	}
	t._finalise2(t.root, &nid)
	t._finalise3(t.root, nid)
	(*t.totalId).Set(big.NewInt(nid))
	t.Size = uint32(nid)
	//fmt.Println(nid)
}

func (t *Trie) _finalise(parent **node.Node, nid *int64) {

	n := *parent

	if n != nil {

		left := n.Left
		right := n.Right
		// Prune where two child nodes are the same
		if left != nil && right != nil && left.Data() != nil && left.Data() == right.Data() {
			//fmt.Println("found a node that can be removed")
			//fmt.Println("parent", *parent)
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
	//fmt.Println(n)
	if n != nil {

		if n.Data() == nil {
			id := big.NewInt(*nid)
			n.SetId(&id)
		} else {
			d := fmt.Sprintf("%x", n.Data())
			dataOffset, _ := t.dataMap[d]
			id := big.NewInt(int64(uint32(dataOffset) + 16) * -1)
			n.SetId(&id)
		}
		if n.Left != nil || n.Right != nil {
			*nid++

			//children := n.Children()
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
	t._print(t.root)
}

func (t *Trie) _print(n *node.Node) {
	if n != nil {

		fmt.Println("n", n)

		//children := n.Children()
		t._print(n.Left)
		t._print(n.Right)
	}
}

func (t *Trie) Serialise(n *node.Node, bytes *[]byte) {

	if n == nil {
		return
	}
	//if n.Left != nil && n.Right != nil {
	//	if n.Left.Id.Cmp(n.Right.Id) == 0 {
	//		return
	//	}
	//}
	//if n.Data() != nil {
	//	if (*n.Id).Uint64() < uint64(t.size) {
	//		*n.Id = (*n.Id).Add((*n.Id), big.NewInt(int64(16+t.size)))
	//	}
	//}
	if n.Left == nil && n.Right == nil {
		return
	}
	//if n.Left != nil || n.Right != nil {
	*bytes = append(*bytes, n.Bytes(t.recordSize, (*t.totalId).Uint64())...)
	//fmt.Println(bytes, n.Id, n.Left, n.Right)
	t.Serialise(n.Left, bytes)
	t.Serialise(n.Right, bytes)
	//}
}

func (t Trie) Bytes() []byte {
	bytes := make([]byte, 0)

	//bytes = append(bytes, t.root.Bytes(t.recordSize, t.size)...)
	//fmt.Println("left", t.root.left.id, "right", t.root.right)
	t.Serialise(t.root, &bytes)

	//queue := make([]*node.Node, 1)
	//visited := make(map[*node.Node]bool)
	//queue[0] = t.root
	//
	//for len(queue) > 0 {
	//
	//	n := queue[0]
	//
	//	if _, f := visited[n]; f {
	//		queue = queue[1:]
	//		continue
	//	}
	//
	//	visited[n] = true
	//	if n.Data != nil {
	//		if n.Id.Uint64() < uint64(t.size) {
	//			n.Id = n.Id.Add(n.Id, big.NewInt(int64(16+t.size)))
	//		}
	//	}
	//
	//	queue = queue[1:]
	//	//if len(queue) > 0 && n.Left == nil && n.Right == nil && n.Depth == queue[0].Depth {
	//	//	//queue[0].Id = n.Id
	//	//	n.Left = queue[0].Left
	//	//	n.Right = queue[0].Right
	//	//	fmt.Println("n", n, "q", queue[0])
	//	//	queue = queue[1:]
	//	//}
	//	bytes = append(bytes, n.Bytes(t.recordSize, t.size)...)
	//
	//	if n.Left != nil {
	//		queue = append(queue, n.Left)
	//	}
	//	if n.Right != nil {
	//		queue = append(queue, n.Right)
	//	}
	//
	//}

	bytes = append(bytes, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)
	bytes = append(bytes, t.data...)
	bytes = append(bytes, 0xAB, 0xCD, 0xEF, 'M', 'a', 'x', 'M', 'i', 'n', 'd', '.', 'c', 'o', 'm')

	//t.Serialise(t.root.left, &bytes)
	//t.Serialise(t.root.right, &bytes)

	return bytes
}

// Determines if the 'bit' in the IP is set
// 'bit' is calculated from the most significant byte first
func isSet(ip net.IP, bit int) bool {
	whichByte := bit / 8
	ipByte := ip[whichByte]
	return ((ipByte >> (7 - (bit % 8))) & 1) > 0
}