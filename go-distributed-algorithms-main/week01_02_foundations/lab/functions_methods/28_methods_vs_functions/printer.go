package main

import "fmt"

type Printer struct{}

func (p *Printer) Print(m string) {
	fmt.Println(m)
}
