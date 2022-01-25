package card

import "github.com/veandco/go-sdl2/sdl"

type Card struct {
	Rank          int32
	Suit          string
	Frame         *sdl.Rect
	Texture       *sdl.Texture
	IsFlippedDown bool
}

const (
	Width      int32 = 116
	Height     int32 = 176
	Spacing    int32 = 50
	MaxShuffle int   = 1000
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

func Suits() []string {
	return []string{"s", "h", "d", "c"}
}

func Ranks() []int32 {
	return []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}
}
