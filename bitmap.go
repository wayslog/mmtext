// Package main provides ...
package main

func NewBitOne() *BitOne {
	b := &BitOne{[15625]int64{}, nil}
	b.link = NewLinkedList()
	return b
}

type BitOne struct {
	Value [15625]int64
	link  *LinkedList
}

func (b *BitOne) Fetch() *LinkedNode {
	return b.link.PopHead()
}

func (b *BitOne) PushOne(v1, v2 int) {
	b.Push(v1)
	b.link.InsertByOrder(v1, v2)
}
func (b *BitOne) Validate(value int) bool {
	offset, pos := value/64, value%64
	var c int64 = 1
	c = c << uint(pos)
	return (b.Value[offset] & c) != 0
}
func (b *BitOne) Push(value int) {
	offset, pos := value/64, value%64
	var c int64 = 1
	c = c << uint(pos)
	b.Value[offset] = b.Value[offset] | c
}

type BitTwo struct {
	Value [13]byte
}

func NewBitTwo() *BitTwo {
	b := &BitTwo{[13]byte{}}
	return b
}

func (b *BitTwo) Push(value int) {
	offset, pos := value/8, value%8
	var c byte = 1
	c = c << uint(pos)
	b.Value[offset] = b.Value[offset] | c
}
func (b *BitTwo) Validate(value int) bool {
	offset, pos := value/8, value%8
	var c byte = 1
	c = c << uint(pos)
	return (b.Value[offset] & c) != 0
}

type LinkedNode struct {
	Bit   BitMap
	Order int
	Next  *LinkedNode
}

func NewLinkedList() *LinkedList {
	l := &LinkedList{}
	l.head = &LinkedNode{NewBitTwo(), -1, nil}
	return l
}

type LinkedList struct {
	head *LinkedNode
}

func (l *LinkedList) PopHead() *LinkedNode {
	head, ln := l.Head(), l.Head().Next
	head.Next = ln.Next
	ln.Next = nil
	return ln
}

func (l *LinkedList) Head() *LinkedNode {
	return l.head
}
func (l *LinkedList) InsertByOrder(order, value int) {
	lp := l.head
	ln := l.head.Next
	for {
		if ln == nil {
			//insert  of back
			now := &LinkedNode{}
			now.Bit = &BitTwo{[13]byte{}}
			now.Bit.Push(value)
			lp.Next = now
			break
		}
		if ln.Order == order {
			ln.Bit.Push(value)
			break
		} else if ln.Order < order {
			now := &LinkedNode{}
			now.Bit = &BitTwo{[13]byte{}}
			now.Bit.Push(value)
			lp.Next = now
			now.Next = ln
			break
		}
		lp = ln
		ln = ln.Next
	}
}

type LinkedOneNode struct {
	Value int
	Next  *LinkedOneNode
}

func NewLinedOneNode(v int) *LinkedOneNode {
	l := &LinkedOneNode{Value: v, Next: nil}
	return l
}

type LinkedOneList struct {
	root *LinkedOneNode
}

func (l *LinkedOneList) PushFront(ln *LinkedOneNode) {
	ln.Next = l.root.Next
	l.root.Next = ln
}

func (l *LinkedOneList) PushOne(v1, v2 int) {
	iv := v1 << 8
	iv = iv & v2
	ln := NewLinedOneNode(iv)
	l.PushFront(ln)
}
func (l *LinkedOneList) FetchNG() *LinkedOneNode {
	return l.root.Next
}
