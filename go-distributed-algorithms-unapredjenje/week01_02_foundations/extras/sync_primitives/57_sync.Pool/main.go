package main

import (
	"fmt"
	"sync"
)

type MyStruct struct {
	data string
}

func newMyStructPool(initialCap int) *sync.Pool {
	p := &sync.Pool{
		New: func() interface{} {
			return &MyStruct{}
		},
	}
	for i := 0; i < initialCap; i++ {
		p.Put(&MyStruct{})
	}
	return p
}

func main() {
	pool := newMyStructPool(5)
	for i := range 10 {
		item := pool.Get().(*MyStruct)
		item.data = fmt.Sprintf("Item %d", i)
		fmt.Println(item.data)
		pool.Put(item)
	}
}
