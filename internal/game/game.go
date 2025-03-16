package game

import (
	"image"
	"math"
	"toy/internal/spine"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	screenWidth  int
	screenHeight int
	title        string
	spine        *spine.SpineData
}

func (g *Game) Run() error {
	ebiten.SetWindowSize(g.screenWidth, g.screenHeight)
	ebiten.SetWindowTitle(g.title)
	return ebiten.RunGame(g)
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Reset()
	op.GeoM.Translate(-84, -86)

	op.GeoM.Rotate(2 * math.Pi * float64(90) / 360)
	op.GeoM.Translate(float64(g.screenWidth)/2, float64(g.screenHeight)/2)

	//"x": 7.8, "y": 71.88, "rotation": 0.29, "width": 172, "height": 173
	screen.DrawImage(g.spine.Img.SubImage(image.Rect(2, 2, 168, 172)).(*ebiten.Image), op)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.screenWidth, g.screenHeight
}

type Option func(*Game) error

func NewGame(opts ...Option) (*Game, error) {
	g := &Game{}
	for _, opt := range opts {
		if err := opt(g); err != nil {
			return nil, err
		}
	}
	return g, nil
}

func WithScreenSize(width, height int) Option {
	return func(g *Game) error {
		g.screenWidth = width
		g.screenHeight = height
		return nil
	}
}

func WithTitle(title string) Option {
	return func(g *Game) error {
		g.title = title
		return nil
	}
}

func WithSpine(s *spine.SpineData) Option {
	return func(g *Game) error {
		g.spine = s
		return nil
	}
}
