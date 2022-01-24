package card

import "github.com/veandco/go-sdl2/sdl"

type Card struct {
	Rank    int32
	Suit    string
	Frame   *sdl.Rect
	Texture *sdl.Texture
}

const (
	Width  int32 = 116
	Height int32 = 176
)

func New(rank int32, suit string, x, y int32, texture *sdl.Texture) *Card {
	return &Card{
		Rank: rank,
		Suit: suit,
		Frame: &sdl.Rect{
			X: x,
			Y: y,
			H: Height,
			W: Width,
		},
		Texture: texture,
	}
}
