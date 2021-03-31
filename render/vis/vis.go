package vis

import (
	"image/color"

	"github.com/chewxy/ll"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type G struct {
	*ll.World
	scaleH, scaleW float64
}

func Make(w *ll.World) G {
	s := w.G.Shape()
	return G{
		World:  w,
		scaleH: float64(800 / s[0]),
		scaleW: float64(800 / s[1]),
	}
}

func (g G) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 800, 800 // TODO make this not const.
}

func (g G) Update(screen *ebiten.Image) error {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		g.handleMouse(x, y)
	}
	return g.Step()
}

func (g G) Draw(screen *ebiten.Image) {
	for i := range g.A {
		for j := range g.A[i] {
			c := color.Black
			if g.A[i][j] == 1 {
				c = color.White
			}
			ebitenutil.DrawRect(screen, float64(i)*g.scaleH, float64(j)*g.scaleW, g.scaleW, g.scaleH, c)
		}
	}
}

func (g G) handleMouse(x, y int) {

	j := y / int(g.scaleH)
	i := x / int(g.scaleW)

	switch g.A[i][j] {
	case 0.0:
		g.A[i][j] = 1
	case 1.0:
		g.A[i][j] = 0
	default:
		panic("Cannot handle non-binary numbers")
	}

}
