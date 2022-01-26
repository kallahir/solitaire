package board

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/kallahir/solitaire/card"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	NumberOfColumns = 7
	DrawPosition    = "dwp"
	DiscardPosition = "ddp"
	SuitPile1       = "s1"
	SuitPile2       = "s2"
	SuitPile3       = "s3"
	SuitPile4       = "s4"
	Column1         = "c1"
	Column2         = "c2"
	Column3         = "c3"
	Column4         = "c4"
	Column5         = "c5"
	Column6         = "c6"
	Column7         = "c7"
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
		column := []*card.Card{card.New(-1, "-1", int32(i)*card.Width, card.Height+card.Spacing, empty)}
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
	for i := range card.Suits() {
		suitePile[i] = append(suitePile[i], card.New(-1, "-1", int32(i)*card.Width, 0, empty))
	}

	return &Board{
		DrawPile:         deck,
		DiscardPile:      []*card.Card{card.New(-1, "-1", int32(5)*card.Width, 0, empty)},
		Columns:          columns,
		SuitPile:         suitePile,
		EmptyCardTexture: empty,
		BackCardTexture:  back,
	}
}

func (b *Board) DrawCard() {
	if len(b.DrawPile) > 0 {
		b.DiscardPile = append([]*card.Card{b.DrawPile[len(b.DrawPile)-1]}, b.DiscardPile...)
		b.DrawPile = b.DrawPile[:len(b.DrawPile)-1]
	} else {
		b.DrawPile = b.DiscardPile[:len(b.DiscardPile)-1]
		b.DiscardPile = []*card.Card{card.New(-1, "-1", int32(5)*card.Width, 0, b.EmptyCardTexture)}
	}
}

func (b *Board) CheckPosition(x, y int32) (string, *card.Card) {
	if checkCollision(x, y, &sdl.Rect{X: 6 * card.Width, Y: 0, H: card.Height, W: card.Width}) {
		fmt.Println("DRAW PILE")
		return DrawPosition, nil
	}
	if checkCollision(x, y, b.DiscardPile[0].Frame) {
		fmt.Println("DISCARD PILE | CARD: ", b.DiscardPile[0])
		if len(b.DiscardPile) == 1 {
			return "", nil
		}
		return DiscardPosition, b.DiscardPile[0]
	}
	for i := range card.Suits() {
		pile := b.SuitPile[i]
		if len(pile) == 1 {
			break
		}
		if checkCollision(x, y, pile[len(pile)-1].Frame) {
			fmt.Println("SUIT PILE #", i+1, " | CARD: ", pile[len(pile)-1])
			return fmt.Sprintf("s%d", i+1), pile[len(pile)-1]
		}
	}
	for i := range b.Columns {
		column := b.Columns[i]
		if len(column) == 1 {
			break
		}
		if checkCollision(x, y, column[len(column)-1].Frame) {
			fmt.Println("COLUMN #", i+1, " | CARD: ", column[len(column)-1])
			return fmt.Sprintf("c%d", i+1), column[len(column)-1]
		}
	}
	return "", nil
}

func checkCollision(x, y int32, frame *sdl.Rect) bool {
	if x > frame.X && x < frame.X+frame.W && y > frame.Y && y < frame.Y+frame.H {
		return true
	}
	return false
}
