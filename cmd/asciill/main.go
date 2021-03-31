package main

import (
	"fmt"
	"github.com/chewxy/ll"
	"gorgonia.org/tensor"
)

func main() {
	w, _ := ll.New(tensor.Shape{3, 3}, ll.Rule{[]int{2, 3}, []int{3}}, Plane, flase)

	w.Set(ll.CV{0, 2, 1}, ll.CV{1, 2, 1}, ll.CV{2, 2, 1})
	fmt.Printf("%#v\n", w.G)

	for i := 0; i < 5; i++ {
		w.Step()
		fmt.Printf("%#v\n", w.G)
		time.Sleep(time.Second)
	}

}
