package main

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"gorgonia.org/tensor"
	"gorgonia.org/tensor/native"
)

// World represents the world in which the cells live and die in.
type World struct {
	W *tensor.Dense
	A [][]float64

	buf  *tensor.Dense
	b    [][]float64
	view *tensor.Dense // buf, but same size as W
	Rule
}

func NewWorld(s tensor.Shape, r Rule) (*World, error) {
	if s.Dims() != 2 {
		return nil, errors.Errorf("Only two dimensional shapes are allowed. Got %v instead", s)
	}
	w := tensor.New(tensor.WithShape(s...), tensor.Of(tensor.Float64))
	a, err := native.MatrixF64(w)
	if err != nil {
		return nil, err
	}
	s2 := s.Clone()
	for i := range s2 {
		s2[i] += 2
	}
	buf := tensor.New(tensor.WithShape(s2...), tensor.Of(tensor.Float64))
	b, err := native.MatrixF64(buf)
	if err != nil {
		return nil, err
	}
	view, err := buf.Slice(S(1, s2[0]-1), S(1, s2[1]-1))
	if err != nil {
		return nil, err
	}
	return &World{
		W:    w,
		A:    a,
		buf:  buf,
		b:    b,
		view: view.(*tensor.Dense),
		Rule: r,
	}, nil
}

func (w *World) Set(cs ...CV) {
	for _, c := range cs {
		w.A[c.Y][c.X] = c.V
	}
}
func (w *World) Step() error {
	if err := tensor.Copy(w.view, w.W); err != nil {
		return err
	}
	for i := 1; i < w.buf.Shape()[0]-1; i++ {
		for j := 1; j < w.buf.Shape()[1]-1; j++ {
			slice, err := w.buf.Slice(S(i-1, i+2), S(j-1, j+2))
			if err != nil {
				return err
			}
			sum, err := tensor.Sum(slice)
			if err != nil {
				return err
			}
			s := int(sum.Data().(float64) - w.b[i][j])
			for _, b := range w.Rule.B {
				if s == b {
					w.A[i-1][j-1] = 1
					break
				}
			}
			var survives bool
			for _, ss := range w.Rule.S {
				if s == ss {
					survives = true
					break
				}
			}
			if !survives {
				w.A[i-1][j-1] = 0
			}
		}
	}
	// process here
	return nil
}

type CV struct {
	X, Y int
	V    float64
}

// Rule represents the rule that the world runs in, in the B/S format.
type Rule struct {
	B []int
	S []int
}

func main() {
	w, _ := NewWorld(tensor.Shape{9, 9}, Rule{[]int{1, 3, 5, 7}, []int{1, 3, 5, 7}})
	w.Set(CV{1, 1, 1})
	fmt.Printf("%v\n", w.Rule)
	fmt.Printf("%#v\n%v\n", w.W, w.A)

	for i := 0; i < 10; i++ {
		w.Step()
		fmt.Printf("%#v\n", w.W)
		time.Sleep(time.Second)
	}

}
