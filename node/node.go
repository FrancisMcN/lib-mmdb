package node

import (
	"fmt"
	"github.com/FrancisMcN/lib-mmdb/field"
	"math"
	"math/big"
)

type Node struct {
	id    **big.Int
	data  field.Field
	Left  *Node
	Right *Node
}

func NewNode() *Node {
	id := big.NewInt(0)
	return &Node{
		id: &id,
	}
}

func (n *Node) String() string {
	id := big.NewInt(0)
	data := ""

	if n.id != nil {
		id = *n.id
	}
	if n.data != nil {
		data = fmt.Sprintf(" [%s]", n.data.String())
	}
	left := "nil"
	if n.Left != nil {
		left = (*n.Left.id).String()
	}
	right := "nil"
	if n.Right != nil {
		right = (*n.Right.id).String()
	}

	return fmt.Sprintf("id: %s (%s, %s) %s", id, left, right, data)
}

func (n *Node) SetLeft(left *Node) {
	n.Left = left
}

func (n *Node) SetRight(right *Node) {
	n.Right = right
}

func (n *Node) SetId(id **big.Int) {
	n.id = id
}

func (n *Node) Id() **big.Int {
	return n.id
}

func (n *Node) SetData(data field.Field) {
	n.data = data
}

func (n *Node) Data() field.Field {
	return n.data
}

func FromBytes(bytes []byte, recordSize uint16) [2]*big.Int {

	// Is the middle byte shared between the two records?
	var shared byte
	var hasExtraByte bool
	nodeSize := (recordSize / 8) * 2
	if recordSize%8 > 0 {
		nodeSize++
		hasExtraByte = true
	}

	if len(bytes) > (int(recordSize)/8)*2 {
		// recordSize / 8 gives the record size in bytes
		// A node is 2 records, so recordSize / 8 gives the middle byte if
		// the node has an odd number of bytes.
		mid := recordSize / 8
		shared = bytes[mid]
	}
	//fmt.Println(fmt.Sprintf("%08b", shared))
	leftBytes := bytes[:recordSize/8]
	if hasExtraByte {
		//fmt.Println("shared: ", []byte{ shared & 0b1111_0000 })
		//leftBytes = append([]byte{ shared & 0b1111_0000 }, leftBytes...)
		leftBytes = append([]byte{shared & 0b1111_0000 >> 4}, leftBytes...)
		//fmt.Println(fmt.Sprintf("%08b", leftBytes))
	}
	left := big.NewInt(0)
	left.SetBytes(leftBytes)

	var rightBytes []byte
	rightBytes = bytes[nodeSize/2 : nodeSize]
	//fmt.Println("b", bytes[recordSize / 8:])
	if hasExtraByte {
		rightBytes = rightBytes[1:]
		//fmt.Println("shared: ", []byte{ shared & 0b0000_1111 })
		//rightBytes = append([]byte{ shared & 0b0000_1111 }, rightBytes...)
		//rightBytes[0] |= shared & 0b0000_1111
		rightBytes = append([]byte{shared & 0b0000_1111}, rightBytes...)
	}
	//fmt.Println("rb", rightBytes, fmt.Sprintf("%x", rightBytes))
	right := big.NewInt(0)
	right.SetBytes(rightBytes)
	//fmt.Println("rb", rightBytes, fmt.Sprintf("%x", rightBytes))
	//fmt.Println("left", left, leftBytes)
	//fmt.Println("right", right, rightBytes)
	//fmt.Println(fmt.Sprintf("bytes: %08b", bytes))
	//fmt.Println(fmt.Sprintf("shared: %08b", shared))
	//return [2]*big.Int{
	//	left,
	//	right,
	//}
	//fmt.Println(leftBytes, rightBytes)
	l := big.NewInt(0)
	l.SetBytes(left.Bytes())

	r := big.NewInt(0)
	r.SetBytes(right.Bytes())

	b := [2]*big.Int{
		l,
		r,
	}
	//fmt.Println("b", b, left.Bytes(), right.Bytes(), bytes)
	//fmt.Println("b", bytes)
	return b

}

func (n Node) Bytes(recordSize int, initialSize uint64) []byte {

	b := make([]byte, 0)
	var l = big.NewInt(int64(initialSize))
	if n.Left != nil {
		l = *n.Left.Id()
	}
	lBytes := l.Bytes()

	if float64(len(lBytes)) < float64(recordSize)/8 {
		diff := int(math.Ceil(float64(recordSize)/8)) - len(lBytes)
		lBytesNew := make([]byte, diff)
		lBytes = append(lBytesNew, lBytes...)
	}
	var r = big.NewInt(int64(initialSize))
	if n.Right != nil {
		r = *n.Right.Id()
	}
	rBytes := r.Bytes()

	if float64(len(rBytes)) < float64(recordSize)/8 {
		diff := int(math.Ceil(float64(recordSize)/8)) - len(rBytes)
		rBytesNew := make([]byte, diff)
		rBytes = append(rBytesNew, rBytes...)
	}
	// If the record has the extra byte
	// For example, a 28 bit record requires 32 bits to store
	if len(lBytes) > recordSize/8 {
		extraLeftByte := lBytes[0]
		// Only the 4 least significant bits are required
		extraLeftByte &= 0b0000_1111
		extraLeftByte <<= 4
		b = append(b, lBytes[1:]...)
		b = append(b, extraLeftByte)

		extraRightByte := rBytes[0]
		// Only the 4 least significant bits are required
		extraRightByte &= 0b0000_1111
		//fmt.Println("extra", extraRightByte, rBytes)
		//// Need to shift 4 bits to the left, to make room for the 4 bits from the right byte
		//extraLeftByte >>= 4
		// Add 4 bits from the right byte to the extra byte
		b[len(b)-1] |= extraRightByte
		b = append(b, rBytes[1:]...)
	} else {
		b = append(b, lBytes...)
		b = append(b, rBytes...)
	}
	return b
}
