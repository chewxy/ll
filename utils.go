package ll

import "gorgonia.org/tensor"

var _ tensor.Slice = sli{}

type sli struct{ s, e, t int }

func S(s ...int) sli {
	switch len(s) {
	case 1:
		return sli{s[0], s[0] + 1, 1}
	case 2:
		return sli{s[0], s[1], 1}
	case 3:
		return sli{s[0], s[1], s[2]}
	default:
		panic("Foo")
	}

}
func (s sli) Start() int { return s.s }
func (s sli) End() int   { return s.e }
func (s sli) Step() int  { return s.t }

func Sum(a *tensor.Dense, it tensor.Iterator) float64 {
	it.Reset()
	data := a.Data().([]float64)
	var sum float64

	for i, err := it.Start(); err == nil; i, err = it.Next() {
		sum += data[i]
	}
	return sum
}
