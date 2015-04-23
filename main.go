// Package main provides ...
package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"runtime"
	"sync"
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

func Analyse(trie *Trie, filepaths ...string) {
	sumstr := []byte{}
	for _, v := range filepaths {
		str, err := ioutil.ReadFile(v)
		if err != nil {
			log.Fatal(err)
		}
		sumstr = append(sumstr, str...)
	}
	sumLines := [1000001][]byte{}
	order := 1
	begin := 0
	end := 0
	for idx, v := range sumstr {
		if v == 10 {
			end = idx
			sumLines[order] = sumstr[begin:end]
			order += 1
			begin = end + 1
		}
	}
	goLimit := make(chan int, 5)
	wg := &sync.WaitGroup{}

	for k, v := range sumLines {
		if k == 0 {
			continue
		} else if v == nil {
			break
		}
		goLimit <- 1
		//log.Printf("analysing %d", k)
		wg.Add(1)
		go func(w *sync.WaitGroup, glmt chan int, target []byte) {
			defer w.Done()
			for i := 1; i <= 100; i++ {
				for j := 0; j < 100-i; j++ {
					tmpSlice := target[j : j+i]
					trie.Insert(tmpSlice, k, j)
				}
			}
			<-glmt
		}(wg, goLimit, v)
	}
	wg.Wait()
}

func main() {
	runtime.GOMAXPROCS(4)
	trie := NewTrie()
	Analyse(trie, "../all_gen.data")
	fmt.Println("Input Your Sequence:")
	for {
		str := ""
		fmt.Scanf("%s", &str)
		iter, err := trie.GenIterator([]byte(str))
		lenOfStr := len(str)
		if err != nil {
			log.Println(err)
		}
		for {
			one, two, err := iter.FetchNext()
			if err != nil {
				log.Println(err)
			}
			fmt.Printf("%d->(%d,%d)", one, two, two+lenOfStr)
		}
	}
}
