package main

import (
	"os"
	"runtime/pprof"

	"github.com/chewxy/ll"
	"github.com/chewxy/ll/render/vis"
	"github.com/hajimehoshi/ebiten"
	"gorgonia.org/tensor"
)

func main() {
	f, err := os.Create("prof.prof")
	if err != nil {
		panic(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	w, _ := ll.New(tensor.Shape{800, 800}, ll.Rule{[]int{1, 3, 5, 7}, []int{1, 3, 5, 7}}, ll.Plane, false)
	g := vis.Make(w)

	ebiten.SetMaxTPS(20)
	ebiten.SetWindowSize(800, 800)
	ebiten.SetWindowTitle("Hello Livestream")
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
