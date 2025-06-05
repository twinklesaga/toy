package game

import (
	"toy/internal/spine"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	screenWidth  int
	screenHeight int
	title        string
	spine        *spine.Spine
}

func (g *Game) Run() error {
	ebiten.SetWindowSize(g.screenWidth, g.screenHeight)
	ebiten.SetWindowTitle(g.title)
	return ebiten.RunGame(g)
}

func (g *Game) Update() error {
	if g.spine != nil {
		g.spine.Update(1.0 / 60.0)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.spine == nil {
		return
	}
	g.spine.Draw(screen, float64(g.screenWidth)/2, float64(g.screenHeight)/2)
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

func WithSpine(s *spine.Spine) Option {
	return func(g *Game) error {
		g.spine = s
		return nil
	}
}
