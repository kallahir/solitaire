package board

import (
	"github.com/kallahir/solitaire/card"
	"github.com/veandco/go-sdl2/sdl"
)

type Board struct {
	DrawPile    []*card.Card
	DiscardPile []*card.Card
	Columns     [7]*card.Card
	SuitPile    [4]*card.Card
	// Common Textures
	EmptyCardTexture *sdl.Texture
	BackCardTexture  *sdl.Texture
}

func New(empty, back *sdl.Texture) (*Board, error) {
	return &Board{
		DrawPile:         []*card.Card{},
		DiscardPile:      []*card.Card{},
		Columns:          [7]*card.Card{},
		SuitPile:         [4]*card.Card{},
		EmptyCardTexture: empty,
		BackCardTexture:  back,
	}, nil
}
