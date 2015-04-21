// Package main provides ...
package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
)

type BitMap interface {
	Push(int)
	Validate(int) bool
}

type Two interface {
	Push(int) int
	Validate(int)
}

type One interface {
	PushOne(int, int)
	//ValidateOne(int, int) bool
	Fetch() *LinkedNode
}

type Iterator struct {
	one           One
	NowOne        int
	NowTwo        int
	NowLinkedNode *LinkedNode
}

func NewIter() *Iterator {
	it := &Iterator{}
	one := NewBitOne()
	it.one = one
	it.NowOne = 1000000 - 1
	it.NowTwo = 99
	return it
}
func (iter *Iterator) PushIter(v1, v2 int) {
	iter.one.PushOne(v1, v2)
}
func (iter *Iterator) FetchNext() (int, int, error) {
	for {
		if iter.NowTwo == -1 {
			iter.NowTwo = 99
			iter.NowOne -= 1
			iter.NowLinkedNode = iter.one.Fetch()
			if iter.NowLinkedNode == nil || iter.NowOne < 0 {
				return -1, -1, errors.New("No More Iterator, All Done")
			}
		} else {
			iter.NowTwo -= 1
		}
		if iter.NowLinkedNode.Bit.Validate(iter.NowTwo) {
			return iter.NowOne, iter.NowTwo, nil
		}
	}
}

func BuildTrie(filepaths ...string) {
	for k, v := range filepaths {
		fmt.Println(k, v)
	}
}

func Analyse(out string, filepaths ...string) {

	for _, v := range filepaths {
		fp, err := os.Open(v)
		if err != nil {
			log.Fatal(err)
		}
		defer fp.Close()
		bufr := bufio.NewReader(fp)
		for {
			line, err := bufr.ReadSlice('\n')
			if err != nil && err == io.EOF {
				break
			}
			fmt.Println(line)
		}
		for {

		}
	}

	log.Println("printing !!!")
	outbytes := []byte{}
	ioutil.WriteFile(out, outbytes, 777)
}

func main() {
	runtime.GOMAXPROCS(4)
	Analyse("out.log", "solexa_100_170_1.fa", "solexa_100_170_2.fa")
}
