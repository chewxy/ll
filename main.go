package main

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"gorgonia.org/tensor"
	"gorgonia.org/tensor/native"
)

type Topology int

const (
	Plane Topology = iota
	CylinderH
	CylinderV
	Torus
)

// World represents the world in which the cells live and die in.
type World struct {
	Rule
	G *tensor.Dense
	a [][]float64

	// processing
	buf  *tensor.Dense
	b    [][]float64
	view *tensor.Dense // buf, but same size as W

	// config
	isOuterIsotropic bool
	topology         Topology
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
		G:    w,
		a:    a,
		buf:  buf,
		b:    b,
		view: view.(*tensor.Dense),
		Rule: r,
	}, nil
}

func (w *World) Set(cs ...CV) {
	for _, c := range cs {
		w.a[c.Y][c.X] = c.V
	}
}
func (w *World) Step() error {
	w.buf.Zero()
	if err := tensor.Copy(w.view, w.G); err != nil {
		return err
	}
	switch w.topology {
	case Plane:
	case Torus:
		fallthrough
	case CylinderH:
		for i := range w.b {
			w.b[i][0] = w.b[i][len(w.b[i])-2]
			w.b[i][len(w.b)-1] = w.b[i][1]
		}
		if w.topology != Torus {
			break
		}
		fallthrough
	case CylinderV:
		copy(w.b[0], w.b[len(w.b)-2])
		copy(w.b[len(w.b)-1], w.b[1])
	default:
		return errors.Errorf("Topology %v not supported", w.topology)
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
			s := int(sum.Data().(float64))
			if w.isOuterIsotropic {
				s -= int(w.b[i][j])
			}

			for _, b := range w.Rule.B {
				if s == b {
					w.a[i-1][j-1] = 1
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
				w.a[i-1][j-1] = 0
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
	w, _ := NewWorld(tensor.Shape{3, 3}, Rule{[]int{2, 3}, []int{3}})
	w.topology = Torus
	w.Set(CV{0, 2, 1}, CV{1, 2, 1}, CV{2, 2, 1})
	fmt.Printf("%#v\n", w.G)

	for i := 0; i < 5; i++ {
		w.Step()
		fmt.Printf("%#v\n", w.G)
		time.Sleep(time.Second)
	}

}
