package main

import (
	"github.com/chewxy/ll"
	"github.com/chewxy/ll/render/vis"
	"github.com/hajimehoshi/ebiten"
	"gorgonia.org/tensor"
)

func main() {
	w, _ := ll.New(tensor.Shape{80, 80}, ll.Rule{[]int{1, 3, 5, 7}, []int{1, 3, 5, 7}}, ll.Torus, false)
	g := vis.Make(w)

	ebiten.SetMaxTPS(20)
	ebiten.SetWindowSize(800, 800)
	ebiten.SetWindowTitle("Hello Livestream")
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
