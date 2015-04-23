// Package main provides ...
package main

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
)

type TrieNode struct {
	Iter     *Iterator
	Children [4]*TrieNode
	Value    byte
}

func (tn *TrieNode) PushIter(v1, v2 int) {
	if tn.Iter == nil {
		tn.Iter = NewIter()
	}
	tn.Iter.PushIter(v1, v2)
}

func NewTrieNode(value byte) *TrieNode {
	tn := &TrieNode{Iter: nil, Children: [4]*TrieNode{}, Value: value}
	return tn
}

func (t *Trie) Insert(path []byte, v1 int, v2 int) {
	lenofp := len(path)
	iter := t.Root
	//get special trie node
	for i := 0; i < lenofp; i++ {
		letter := path[i]
		it := iter.Children[t.PosMap[letter]]
		if it == nil {
			it = NewTrieNode(letter)
			iter.Children[t.PosMap[letter]] = it
		}
		iter = it
	}
	iter.PushIter(v1, v2)
}

func (t Trie) GenIterator(path []byte) (*Iterator, error) {
	iter := t.Root
	for _, v := range path {
		iter := iter.Children[t.PosMap[v]]
		if iter == nil {
			return nil, errors.New(fmt.Sprintf("Can't not find this K-mer %s", path))
		}
	}
	//deep copy
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(iter.Iter); err != nil {
		return nil, err
	}
	itern := NewIter()
	if err := gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(itern); err != nil {
		return nil, err
	}
	return itern, nil
}

type Trie struct {
	Root   *TrieNode
	PosMap map[byte]int
}

func NewTrie() *Trie {
	t := &Trie{&TrieNode{Iter: nil, Children: [4]*TrieNode{}, Value: 0}, make(map[byte]int)}
	t.PosMap['A'] = 0
	t.PosMap['C'] = 1
	t.PosMap['G'] = 2
	t.PosMap['T'] = 3
	return t
}
