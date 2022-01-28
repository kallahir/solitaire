package board

import (
	"math/rand"
	"time"

	"github.com/kallahir/solitaire/card"
	"github.com/kallahir/solitaire/renderwindow"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	NumberOfColumns = 7
	DrawPosition    = "dwp"
	DiscardPosition = "ddp"
)

type Board struct {
	DrawPile    []*card.Card
	DiscardPile []*card.Card
	Columns     [7][]*card.Card
	SuitPile    [4][]*card.Card
	Textures    map[string]*sdl.Texture
	IsRunning   bool
	// Hand Variables
	Hand                 []*card.Card
	HandPreviousLocation string
}

func New(rw *renderwindow.RenderWindow, textures map[string]*sdl.Texture) *Board {
	var deck []*card.Card
	for _, rank := range card.Ranks() {
		for _, suit := range card.Suits() {
			deck = append(deck, card.New(rank, suit))
		}
	}

	// Shuffle Cards random number of times up to card.MaxShuffle
	for i := 0; i < rand.Intn(card.MaxShuffle); i++ {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })
	}

	// Pick Cards for the Columns
	var columns [7][]*card.Card
	for i := 0; i < NumberOfColumns; i++ {
		column := []*card.Card{card.New(-1, card.Empty)}
		for j := 0; j <= i; j++ {
			// Pick Last Card from the deck
			c := deck[len(deck)-1]
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
		suitePile[i] = append(suitePile[i], card.New(-1, card.Empty))
	}

	deck = append([]*card.Card{card.New(-1, card.Empty)}, deck...)

	return &Board{
		DrawPile:    deck,
		DiscardPile: []*card.Card{card.New(-1, card.Empty)},
		Columns:     columns,
		SuitPile:    suitePile,
		Textures:    textures,
		IsRunning:   true,
	}
}

func (b *Board) Render(rw *renderwindow.RenderWindow, x, y int32) {
	for _, c := range b.DrawPile {
		rw.Render(&sdl.Rect{
			X: int32(NumberOfColumns-1) * card.Width,
			Y: 0,
			W: card.Width,
			H: card.Height,
		}, b.Textures[c.TextureKey])
	}

	for _, c := range b.DiscardPile {
		rw.Render(&sdl.Rect{
			X: int32(NumberOfColumns-2) * card.Width,
			Y: 0,
			W: card.Width,
			H: card.Height,
		}, b.Textures[c.TextureKey])
	}

	for i, suit := range b.SuitPile {
		for _, c := range suit {
			rw.Render(&sdl.Rect{
				X: int32(i) * card.Width,
				Y: 0,
				W: card.Width,
				H: card.Height,
			}, b.Textures[c.TextureKey])
		}
	}

	// FIXME: Place first card over empty card
	for i, column := range b.Columns {
		for j, c := range column {
			verticalSpacing := int32(j) * card.Spacing
			if j == 0 {
				verticalSpacing += card.Spacing
			}
			tk := b.Textures[c.TextureKey]
			if c.IsFlippedDown {
				tk = b.Textures[card.Back]
			}
			rw.Render(&sdl.Rect{
				X: int32(i) * card.Width,
				Y: card.Height + verticalSpacing + card.Spacing,
				W: card.Width,
				H: card.Height,
			}, tk)
		}
	}

	for i, c := range b.Hand {
		rw.Render(&sdl.Rect{
			X: x - card.Width/2,
			Y: y - card.Height/2 + int32(i)*card.Spacing,
			W: card.Width,
			H: card.Height,
		}, b.Textures[c.TextureKey])
	}
}

// TODO1: All sdl.PRESSED events that involves the hand, must records previous location and last card
// 		  so I can flip the card up and/or put the cards back on the previous location
// TODO2: All sdl.RELEASED events that involves the hand, must apply the game rules by checking the card
// 		  on the top of the hand and the card on the bottom of the pile or top of the suit
func (b *Board) HandleClick(x, y int32, mouseState uint8) {
	if CheckCollision(x, y, &sdl.Rect{X: 6 * card.Width, Y: 0, H: card.Height, W: card.Width}) {
		switch {
		case mouseState == sdl.PRESSED && len(b.DrawPile) > 1:
			b.DiscardPile = append(b.DiscardPile, b.DrawPile[len(b.DrawPile)-1])
			b.DrawPile = b.DrawPile[:len(b.DrawPile)-1]
		case mouseState == sdl.PRESSED && len(b.DrawPile) == 1 && len(b.DiscardPile) > 1:
			for i := len(b.DiscardPile) - 1; i > 0; i-- {
				b.DrawPile = append(b.DrawPile, b.DiscardPile[i])
			}
			b.DiscardPile = []*card.Card{card.New(-1, card.Empty)}
		}
	}

	if CheckCollision(x, y, &sdl.Rect{X: 5 * card.Width, Y: 0, H: card.Height, W: card.Width}) {
		switch {
		case mouseState == sdl.PRESSED && len(b.Hand) == 0 && len(b.DiscardPile) > 1:
			b.Hand = append(b.Hand, b.DiscardPile[len(b.DiscardPile)-1])
			b.DiscardPile = b.DiscardPile[:len(b.DiscardPile)-1]
			// TODO: A Card can only get back to the Discard Pile if it was the last to be removed from there
			// I think I'll need a flag variable for this, like saving the TextureKey or something like that
			// case mouseState == sdl.RELEASED && len(b.Hand) == 1 && len(b.DiscardPile) > 1 && b.LastDrawn == b.Hand[0].TextureKey:
			// 	b.DiscardPile = append(b.DiscardPile, b.Hand...)
			// 	b.Hand = []*card.Card{}
			// 	b.LastDrawn = ""
		}
	}

	for i := range card.Suits() {
		if CheckCollision(x, y, &sdl.Rect{X: int32(i) * card.Width, Y: 0, H: card.Height, W: card.Width}) {
			switch {
			case mouseState == sdl.RELEASED && len(b.Hand) == 1:
				b.SuitPile[i] = append(b.SuitPile[i], b.Hand...)
				b.Hand = []*card.Card{}
			case mouseState == sdl.PRESSED && len(b.Hand) == 0 && len(b.SuitPile[i]) > 1:
				b.Hand = append(b.Hand, b.SuitPile[i][len(b.SuitPile[i])-1])
				b.SuitPile[i] = b.SuitPile[i][:len(b.SuitPile[i])-1]
			}
			break
		}
	}

	for i := range b.Columns {
		for j, c := range b.Columns[i] {
			switch {
			case c.IsFlippedDown || (j == 0 && mouseState == sdl.PRESSED):
				continue
			case mouseState == sdl.PRESSED && j == len(b.Columns[i])-1 && CheckCollision(x, y, &sdl.Rect{X: int32(i) * card.Width, Y: card.Height + (int32(j) * card.Spacing) + card.Spacing, W: card.Width, H: card.Height}):
				b.Hand = append(b.Hand, b.Columns[i][len(b.Columns[i])-1])
				b.Columns[i] = b.Columns[i][:len(b.Columns[i])-1]
				// FIXME: Flipping the last card from the pile can only happen after the cards are really played
				if len(b.Columns[i]) > 1 {
					b.Columns[i][len(b.Columns[i])-1].IsFlippedDown = false
				}
			case mouseState == sdl.RELEASED && len(b.Hand) > 0 && CheckCollision(x, y, &sdl.Rect{X: int32(i) * card.Width, Y: card.Height + (int32(j) * card.Spacing) + card.Spacing, W: card.Width, H: card.Height}):
				b.Columns[i] = append(b.Columns[i], b.Hand...)
				b.Hand = []*card.Card{}
			case mouseState == sdl.PRESSED && CheckCollision(x, y, &sdl.Rect{X: int32(i) * card.Width, Y: card.Height + (int32(j) * card.Spacing) + card.Spacing, W: card.Width, H: card.Spacing}):
				b.Hand = append(b.Hand, b.Columns[i][j:]...)
				// FIXME: Flipping the last card from the pile can only happen after the cards are really played
				b.Columns[i] = b.Columns[i][:j]
				if len(b.Columns[i]) > 1 {
					b.Columns[i][len(b.Columns[i])-1].IsFlippedDown = false
				}
			}
		}
	}

	// FIXME: This may not apply once every region is handled
	if mouseState == sdl.RELEASED && len(b.Hand) > 0 {
		b.DiscardPile = append(b.DiscardPile, b.Hand...)
		b.Hand = []*card.Card{}
	}
}

func CheckCollision(x, y int32, boudingBox *sdl.Rect) bool {
	if x > boudingBox.X && x < boudingBox.X+boudingBox.W && y > boudingBox.Y && y < boudingBox.Y+boudingBox.H {
		return true
	}
	return false
}
