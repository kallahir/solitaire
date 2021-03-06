package board

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/kallahir/solitaire/card"
	"github.com/kallahir/solitaire/renderwindow"
	"github.com/kallahir/solitaire/utils"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	NumberOfColumns    = 7
	NumberOfCardsToWin = 56
	Columns            = "c"
	Discard            = "d"
	Suit               = "s"
)

type Board struct {
	DrawPile    []*card.Card
	DiscardPile []*card.Card
	Columns     [7][]*card.Card
	SuitPile    [4][]*card.Card
	Textures    map[string]*sdl.Texture
	IsRunning   bool
	// Hand Variables
	Hand       []*card.Card
	HandOrigin string
}

func New(rw *renderwindow.RenderWindow, textures map[string]*sdl.Texture) *Board {
	var deck []*card.Card
	for _, rank := range card.Ranks() {
		for _, suit := range card.Suits() {
			deck = append(deck, card.New(rank, suit))
		}
	}

	// Shuffle Cards
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })

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
	drawTexture := card.Back
	if len(b.DrawPile) == 1 {
		drawTexture = card.Empty
	}
	rw.Render(&sdl.Rect{
		X: int32(NumberOfColumns-1)*card.Width + int32(NumberOfColumns)*card.HSpacing,
		Y: card.HSpacing,
		W: card.Width,
		H: card.Height,
	}, b.Textures[drawTexture])

	for _, c := range b.DiscardPile {
		rw.Render(&sdl.Rect{
			X: int32(NumberOfColumns-2)*card.Width + int32(NumberOfColumns-1)*card.HSpacing,
			Y: card.HSpacing,
			W: card.Width,
			H: card.Height,
		}, b.Textures[c.TextureKey])
	}

	for i, suit := range b.SuitPile {
		for _, c := range suit {
			rw.Render(&sdl.Rect{
				X: int32(i)*card.Width + int32(i+1)*card.HSpacing,
				Y: card.HSpacing,
				W: card.Width,
				H: card.Height,
			}, b.Textures[c.TextureKey])
		}
	}

	for i, column := range b.Columns {
		for j, c := range column {
			verticalSpacing := int32(j)*card.VSpacing + card.VSpacing
			if j == 0 {
				verticalSpacing += card.VSpacing
			}
			tk := c.TextureKey
			if c.IsFlippedDown {
				tk = card.Back
			}
			rw.Render(&sdl.Rect{
				X: int32(i)*card.Width + int32(i+1)*card.HSpacing,
				Y: card.Height + verticalSpacing,
				W: card.Width,
				H: card.Height,
			}, b.Textures[tk])
		}
	}

	for i, c := range b.Hand {
		rw.Render(&sdl.Rect{
			X: x - card.Width/2,
			Y: y - card.Height/4 + int32(i)*card.VSpacing,
			W: card.Width,
			H: card.Height,
		}, b.Textures[c.TextureKey])
	}
}

func (b *Board) HandleClick(x, y int32, mouseState uint8) {
	if utils.CheckCollision(x, y, &sdl.Rect{X: int32(NumberOfColumns-1)*card.Width + int32(NumberOfColumns)*card.HSpacing, Y: card.HSpacing, H: card.Height, W: card.Width}) {
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

	if utils.CheckCollision(x, y, &sdl.Rect{X: int32(NumberOfColumns-2)*card.Width + int32(NumberOfColumns-1)*card.HSpacing, Y: card.HSpacing, H: card.Height, W: card.Width}) {
		switch {
		case mouseState == sdl.PRESSED && len(b.Hand) == 0 && len(b.DiscardPile) > 1:
			b.Hand = append(b.Hand, b.DiscardPile[len(b.DiscardPile)-1])
			b.HandOrigin = Discard
			b.DiscardPile = b.DiscardPile[:len(b.DiscardPile)-1]
		}
	}

	for i := range card.Suits() {
		if utils.CheckCollision(x, y, &sdl.Rect{X: int32(i)*card.Width + int32(i+1)*card.HSpacing, Y: card.HSpacing, H: card.Height, W: card.Width}) {
			switch {
			case mouseState == sdl.RELEASED && len(b.Hand) == 1:
				if !b.ValidateMovement(Suit, i, b.Hand[0]) {
					continue
				}
				b.SuitPile[i] = append(b.SuitPile[i], b.Hand...)
				b.Hand = []*card.Card{}
				b.FlipOriginCard(b.HandOrigin)
				b.HandOrigin = ""
			case mouseState == sdl.PRESSED && len(b.Hand) == 0 && len(b.SuitPile[i]) > 1:
				b.Hand = append(b.Hand, b.SuitPile[i][len(b.SuitPile[i])-1])
				b.HandOrigin = fmt.Sprint(Suit, i)
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
			case mouseState == sdl.PRESSED && j == len(b.Columns[i])-1 && utils.CheckCollision(x, y, &sdl.Rect{X: int32(i)*card.Width + int32(i+1)*card.HSpacing, Y: card.Height + (int32(j) * card.VSpacing) + card.VSpacing, W: card.Width, H: card.Height}):
				b.Hand = append(b.Hand, b.Columns[i][len(b.Columns[i])-1])
				b.HandOrigin = fmt.Sprint(Columns, i)
				b.Columns[i] = b.Columns[i][:len(b.Columns[i])-1]
			case mouseState == sdl.RELEASED && len(b.Hand) > 0 && utils.CheckCollision(x, y, &sdl.Rect{X: int32(i)*card.Width + int32(i+1)*card.HSpacing, Y: card.Height + (int32(j) * card.VSpacing) + card.VSpacing, W: card.Width, H: card.Height}):
				if !b.ValidateMovement(Columns, i, b.Hand[0]) {
					continue
				}
				b.Columns[i] = append(b.Columns[i], b.Hand...)
				b.Hand = []*card.Card{}
				b.FlipOriginCard(b.HandOrigin)
				b.HandOrigin = ""
			case mouseState == sdl.PRESSED && utils.CheckCollision(x, y, &sdl.Rect{X: int32(i)*card.Width + int32(i+1)*card.HSpacing, Y: card.Height + (int32(j) * card.VSpacing) + card.VSpacing, W: card.Width, H: card.VSpacing}):
				b.Hand = append(b.Hand, b.Columns[i][j:]...)
				b.HandOrigin = fmt.Sprint(Columns, i)
				b.Columns[i] = b.Columns[i][:j]
			}
		}
	}

	// Returning any cards from the Hand to its Original Position
	if mouseState == sdl.RELEASED && len(b.Hand) > 0 {
		switch string(b.HandOrigin[0]) {
		case Discard:
			b.DiscardPile = append(b.DiscardPile, b.Hand...)
		case Suit:
			idx, _ := strconv.Atoi(string(b.HandOrigin[1]))
			b.SuitPile[idx] = append(b.SuitPile[idx], b.Hand...)
		case Columns:
			idx, _ := strconv.Atoi(string(b.HandOrigin[1]))
			b.Columns[idx] = append(b.Columns[idx], b.Hand...)
		}
		b.Hand = []*card.Card{}
		b.HandOrigin = ""
	}
}

func (b *Board) HandleRightClick(x, y int32) {
	if utils.CheckCollision(x, y, &sdl.Rect{X: int32(NumberOfColumns-2)*card.Width + int32(NumberOfColumns-1)*card.HSpacing, Y: card.HSpacing, H: card.Height, W: card.Width}) {
		c := b.DiscardPile[len(b.DiscardPile)-1]
		if b.AutoFillSuitPile(c) {
			b.DiscardPile = b.DiscardPile[:len(b.DiscardPile)-1]
		}
	}

	for i := range b.Columns {
		j := len(b.Columns[i]) - 1
		c := b.Columns[i][j]
		if c.IsFlippedDown || j == 0 {
			continue
		}
		if utils.CheckCollision(x, y, &sdl.Rect{X: int32(i)*card.Width + int32(i+1)*card.HSpacing, Y: card.Height + (int32(j) * card.VSpacing) + card.VSpacing, W: card.Width, H: card.Height}) {
			if b.AutoFillSuitPile(c) {
				b.Columns[i] = b.Columns[i][:j]
				b.Columns[i][j-1].IsFlippedDown = false
			}
		}
	}
}

func (b *Board) AutoFillSuitPile(c *card.Card) bool {
	for idx := range b.SuitPile {
		if b.ValidateMovement(Suit, idx, c) {
			b.SuitPile[idx] = append(b.SuitPile[idx], c)
			return true
		}
	}
	return false
}

func (b *Board) FlipOriginCard(origin string) {
	if string(origin[0]) == Columns {
		idx, _ := strconv.Atoi(string(b.HandOrigin[1]))
		b.Columns[idx][len(b.Columns[idx])-1].IsFlippedDown = false
	}
}

func (b *Board) ValidateMovement(pile string, idx int, src *card.Card) bool {
	switch pile {
	case Suit:
		dst := b.SuitPile[idx][len(b.SuitPile[idx])-1]
		if (dst.Suit == src.Suit && dst.Rank == src.Rank-1) || (dst.Rank == -1 && src.Rank == 1) {
			return true
		}
	case Columns:
		dst := b.Columns[idx][len(b.Columns[idx])-1]
		if (dst.CompareOverlappingSuit(src) && dst.Rank == src.Rank+1) || (dst.Rank == -1 && src.Rank == 13) {
			return true
		}
	}
	return false
}

func (b *Board) IsOver() bool {
	total := 0
	for i := range b.SuitPile {
		total += len(b.SuitPile[i])
	}
	return total == NumberOfCardsToWin
}
