// Package main provides ...
package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

const (
	MAXN int64 = 10450
	//MAXN int64 = 3
	MAXM int64 = 99 * 10000
)

var (
	EOI error = errors.New("Iterator get End, no more data")
)

func NewTrieNode() *trieNode {
	tn := &trieNode{chilren: [4]*trieNode{}, settler: nil, count: 0}
	return tn
}

type trieNode struct {
	chilren [4]*trieNode
	settler Settler
	count   int64
}

func NewTrie() *trie {
	t := &trie{make(map[byte]int), NewTrieNode()}
	t.posMap['A'] = 0
	t.posMap['C'] = 1
	t.posMap['G'] = 2
	t.posMap['T'] = 3
	return t
}

type trie struct {
	posMap map[byte]int
	root   *trieNode
}

func (t *trie) Find(path []byte) error {
	tn := t.root
	lenOfp := int32(len(path))
	for _, v := range path {
		tn = tn.chilren[t.posMap[v]]
		if tn == nil {
			return errors.New(fmt.Sprintf("Can't find this path(%s)\n", path))
		}
	}
	if tn.settler == nil {
		return errors.New(fmt.Sprintf("Path(%s) is too short\n", path))
	}
	iter := tn.settler.GenIter()
	for {
		v1, v2, err := iter.Fetch()
		if err != nil {
			if err == EOI {
				return nil
			}
			log.Fatal(err)
		}
		fmt.Printf("%d:(%d,%d)\n", v1+1, v2+1, v2+lenOfp+1)
	}
	return nil
}

func (t *trie) Insert(path []byte, v1, v2 int32) {
	//fmt.Printf("%s\n", path)
	tn := t.root
	for _, v := range path {
		if tn.chilren[t.posMap[v]] == nil {
			tn.chilren[t.posMap[v]] = NewTrieNode()
		}
		tn = tn.chilren[t.posMap[v]]
	}
	if tn.settler == nil {
		tn.settler = NewLinkedSettler()
	}
	tn.count = tn.count + 1
	if tn.count == MAXN {
		//fmt.Println("x")
		tn.settler = ConvLinedToDouble(tn.settler)
	}
	tn.settler.Settle(v1, v2)
}
func (t *trie) GetIter(path []byte) (iterator, error) {
	tn := t.root
	for _, v := range path {
		tn = tn.chilren[t.posMap[v]]
		if tn == nil {
			return nil, errors.New("Can't Generate Iterator")
		}
	}
	return tn.settler.GenIter(), nil
}

func ConvLinedToDouble(s Settler) *doubleSettler {
	iter := s.GenIter()
	ds := NewDoubleSettler()
	ln := ds.root
	var max, now int32 = -1, -1
	for {
		v1, v2, err := iter.Fetch()
		//fmt.Println(v1, v2)
		if err != nil {
			if err == EOI {
				break
			} else {
				log.Fatalf("Can't Conv from linkedSettler to NewDoubleSettler")
			}
		}
		if max == -1 {
			max = v1
		}
		if now != v1 {
			ds.SettleV1(v1)
			lp := &doubleLinkedNode{[13]byte{}, nil}
			ln.next = lp
			ln = lp
			now = v1
		}
		pos := v2 / 8
		offset := v2 % 8
		var dc byte = 1 << uint(offset)
		ln.valueTwo[pos] = ln.valueTwo[pos] | dc
	}
	ds.nowCount = max
	return ds
}

type Settler interface {
	Settle(v1, v2 int32)
	GenIter() iterator
}

type doubleIterator struct {
	dbs        *doubleSettler
	nowTwoNode *doubleLinkedNode //default root.next
	nowOne     int32
	nowTwo     int32
}

func (d *doubleIterator) Fetch() (int32, int32, error) {
	for {
		d.nowTwo = d.nowTwoNode.GettleV2(d.nowTwo - 1)
		if d.nowTwo == -1 {
			d.nowTwo = 99
			d.nowOne = d.dbs.GettleV1(d.nowOne - 1)
			d.nowTwoNode = d.nowTwoNode.next
			if d.nowTwoNode == nil {
				return -1, -1, EOI
			}
			if d.nowTwo == -1 {
				return -1, -1, EOI
			}
		} else {
			break
		}
		return -1, -1, EOI
	}
	if d.nowTwoNode == nil || d.nowOne < 0 {
		return -1, -1, EOI
	}
	return d.nowOne, d.nowTwo, nil
}
func NewDoubleSettler() *doubleSettler {
	ds := &doubleSettler{nowCount: -1, valueOne: [15625]int64{}, root: &doubleLinkedNode{}}
	return ds
}

type doubleSettler struct {
	nowCount int32 //default -1
	valueOne [15625]int64
	root     *doubleLinkedNode
}

type doubleLinkedNode struct {
	valueTwo [13]byte
	next     *doubleLinkedNode
}

func (d *doubleSettler) GettleV1(now int32) int32 {
	pos := now / 64
	offset := now % 64
	for {
		if offset < 0 {
			offset = 63
			pos = pos - 1
		}
		if pos < 0 {
			return -1
		}
		var di int64 = 1 << uint(offset)
		if (d.valueOne[pos] & di) > 0 {
			return pos*64 + offset
		}
		offset = offset - 1
	}
}
func (d *doubleLinkedNode) GettleV2(now int32) int32 {
	pos := now / 8
	offset := (now % 8)
	for {
		if offset < 0 {
			offset = 7
			pos = pos - 1
		}
		if pos < 0 {
			return -1
		}
		var di byte = 1 << uint(offset)
		if (d.valueTwo[pos] & di) > 0 {
			return pos*8 + offset
		}
		offset--
	}
}

//order是由小到达的序列序号的出现顺序，从零开始
func (d *doubleSettler) SettleV2(v2 int32) {
	pos := v2 / 8
	offset := v2 % 8
	var di byte = 1 << uint(offset)
	d.root.next.valueTwo[pos] = d.root.next.valueTwo[pos] | di
}
func (d *doubleSettler) SettleV1(v1 int32) {
	pos := v1 / 64
	offset := v1 % 64
	var di int64 = 1 << uint(offset)
	d.valueOne[pos] = d.valueOne[pos] | di
}

func (d *doubleSettler) Settle(v1, v2 int32) {
	d.SettleV1(v1)
	if v1 != d.nowCount {
		tmpnode := &doubleLinkedNode{}
		tmpnode.valueTwo = [13]byte{}
		tmpnode.next = d.root.next
		d.root.next = tmpnode
		d.nowCount = v1
	}
	d.SettleV2(v2)

}
func (d *doubleSettler) GenIter() iterator {
	iter := &doubleIterator{}
	iter.dbs = d
	iter.nowOne = 1000000
	iter.nowTwo = 100
	iter.nowTwoNode = d.root.next
	iter.nowOne = d.GettleV1(iter.nowOne - 1)
	return iter
}

func NewLinkedSettler() *linkedSettler {
	ls := &linkedSettler{}
	ls.root = &linkedNode{}
	return ls
}

type linkedSettler struct {
	root *linkedNode
}

func (l *linkedSettler) Settle(v1, v2 int32) {
	lt := &linkedNode{}
	lt.next = l.root.next
	l.root.next = lt
	lt.Settle(v1, v2)
}
func (l *linkedSettler) GenIter() iterator {
	iter := &linkedIterator{}
	iter.ln = l.root
	return iter
}

func (l *linkedNode) Settle(v1, v2 int32) {
	v1 = v1 << 8
	l.Value = v1 | v2
}
func (l *linkedNode) Gettle() (int32, int32) {
	return l.Value >> 8, (l.Value << 24) >> 24
}

type linkedNode struct {
	Value int32
	next  *linkedNode
}
type linkedIterator struct {
	ln *linkedNode //default = root
}

func (l *linkedIterator) Fetch() (int32, int32, error) {
	lp := l.ln.next
	l.ln = lp
	if lp == nil {
		return -1, -1, EOI
	}
	v1, v2 := lp.Gettle()
	return v1, v2, nil
}

type iterator interface {
	Fetch() (int32, int32, error)
}

func PushLine(k int32, t *trie, strc chan []byte) {
	var order int32 = -1
	max := 100 - k
	for {
		order += 1
		var i int32 = 0
		line := <-strc
		if line[0] == 'e' {
			break
		}
		for ; i <= max; i++ {
			t.Insert(line[i:i+k], order, i)
		}
	}
}

func BuildTrie(k int32, files ...string) *trie {
	strc := make(chan []byte)
	t := NewTrie()
	go PushLine(k, t, strc)
	for _, v := range files {
		fp, err := os.Open(v)
		defer fp.Close()
		if err != nil {
			log.Fatal(err)
		}
		reader := bufio.NewReader(fp)
		for {
			line, _, err := reader.ReadLine()
			if err != nil {
				if err == io.EOF {
					log.Printf("load %s is done", v)
					break
				}
				log.Fatal(err)
			}
			strc <- line
		}
	}
	eol := []byte("e")
	strc <- eol
	return t
}

var antPosMap map[int]byte = make(map[int]byte)

func WalkTrie(tn *trieNode, prefix []byte) {
	for v, k := range tn.chilren {
		letter := antPosMap[v]
		if k == nil {
			return
		}
		p := append(prefix, letter)
		if k.settler == nil {
			WalkTrie(k, p)
			continue
		}
		lenOfp := int32(len(p))
		iter := k.settler.GenIter()
		fmt.Printf("%s:\n", p)
		for {
			v1, v2, err := iter.Fetch()
			if err != nil {
				if err == EOI {
					break
				}
				log.Fatal(err)
			}
			fmt.Printf("%d %d %d\n", v1+1, v2+1, v2+lenOfp+1)
		}
	}
}

func main() {
	antPosMap[0] = 'A'
	antPosMap[1] = 'C'
	antPosMap[2] = 'G'
	antPosMap[3] = 'T'
	klong, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("%s is not avaliable k-mer", os.Args[1])
	}
	//BuildTrie(int32(klong), "../all_gen.data")
	t := BuildTrie(int32(klong), "../all_gen.data")
	//prefix := []byte{}
	//WalkTrie(t.root, prefix)
	for {
		fmt.Printf("Input Your Path(len=%d):", klong)
		p := []byte{}
		fmt.Scanf("%s", &p)
		lenofp := len(p)
		if lenofp != klong {
			fmt.Println("Wrong length")
		}
		t.Find(p)
	}
}
