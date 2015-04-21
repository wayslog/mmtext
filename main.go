// Package main provides ...
package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
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
	countMap := make(map[string]int)
	goLimit := make(chan int, 200)
	lenMap := [101]chan int{}
	for i := 1; i < 101; i++ {
		lenMap[i] = make(chan int)
	}

	for k, v := range sumLines {
		if k == 0 {
			continue
		} else if v == nil {
			break
		}
		goLimit <- 1
		log.Println("analysing %d", k)
		go func(target []byte) {
			for i := 1; i <= 100; i++ {
				for j := 0; j < 10-i; j++ {
					c := lenMap[i]
					c <- 1
					tmpSlice := target[j : j+i]
					tmpStr := fmt.Sprintf("%s", tmpSlice)
					if val, ok := countMap[tmpStr]; ok {
						countMap[tmpStr] = val + 1
					} else {
						countMap[tmpStr] = 1
					}
					<-c
				}
			}
			<-goLimit
		}(v)
	}
}

func main() {
	runtime.GOMAXPROCS(4)
	Analyse("out.log", "../all_gen.data")
}
