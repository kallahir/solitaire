package board

import (
	"math/rand"
	"time"

	"github.com/kallahir/solitaire/card"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	NumberOfColumns = 7
)

type Board struct {
	DrawPile    []*card.Card
	DiscardPile []*card.Card
	Columns     [7][]*card.Card
	SuitPile    [4][]*card.Card
	// Common Textures
	EmptyCardTexture *sdl.Texture
	BackCardTexture  *sdl.Texture
}

func New(empty, back *sdl.Texture, deck []*card.Card) *Board {
	// Shuffle Cards random number of times up to card.MaxShuffle
	for i := 0; i < rand.Intn(card.MaxShuffle); i++ {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })
	}

	// Pick Cards for the Columns
	var columns [7][]*card.Card
	for i := 0; i < NumberOfColumns; i++ {
		var column []*card.Card
		for j := 0; j <= i; j++ {
			// Pick Last Card from the deck
			c := deck[len(deck)-1]
			// Update Card position
			c.Frame.X, c.Frame.Y = int32(i)*card.Width, card.Height+(int32(j)*card.Spacing)+card.Spacing
			if i != j {
				c.IsFlippedDown = true
			}
			// Add to the current column
			column = append(column, deck[len(deck)-1])
			// Remove used card from the deck
			deck = deck[:len(deck)-1]
		}
		columns[i] = column
	}

	var suitePile [4][]*card.Card
	for i, _ := range card.Suits() {
		suitePile[i] = append(suitePile[i], card.New(-1, "-1", int32(i)*card.Width, 0, empty))
	}

	return &Board{
		DrawPile:         deck,
		DiscardPile:      []*card.Card{},
		Columns:          columns,
		SuitPile:         suitePile,
		EmptyCardTexture: empty,
		BackCardTexture:  back,
	}
}
