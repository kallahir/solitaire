package board

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/kallahir/solitaire/card"
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
	fmt.Println("DRAWING CARD...")
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
		if !b.DiscardPile[0].IsBeingUsed {
			return DiscardPosition, b.DiscardPile[0]
		}
	}
	for i := range card.Suits() {
		card := b.SuitPile[i][len(b.SuitPile[i])-1]
		if !card.IsBeingUsed && checkCollision(x, y, card.Frame) {
			fmt.Println("SUIT PILE #", i, " | CARD: ", card)
			return fmt.Sprintf("s%d", i), card
		}
	}
	for i := range b.Columns {
		card := b.Columns[i][len(b.Columns[i])-1]
		if !card.IsBeingUsed && checkCollision(x, y, card.Frame) {
			fmt.Println("COLUMN #", i, " | CARD: ", card)
			return fmt.Sprintf("c%d", i), card
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

func (b *Board) MoveCard(pc *card.PlayingCard, destinationCard *card.Card, position string) bool {
	from := strings.Split(pc.OriginalPile, "")
	if len(from) == 0 {
		return false
	}

	to := strings.Split(position, "")
	if len(to) > 2 || len(to) == 0 {
		fmt.Println("TRYING TO MOVE CARD TO DRAW PILE, DISCARD PILE OR SAME PLACE")
		return false
	}

	fmt.Println("VALIDATING IF: ", pc, " CAN BE PLACED ON TOP OF: ", destinationCard, " AT POSITION: ", to)
	switch to[0] {
	case "c":
		idx, _ := strconv.Atoi(to[1])
		if len(b.Columns[idx]) == 1 {
			break
		}
		if pc.CardDetails.Rank != destinationCard.Rank-1 || !pc.CardDetails.ValidOverlappingSuit(destinationCard) {
			fmt.Println("CARD CAN'T BE PLACED INTO COLUMN ", idx)
			return false
		}
	case "s":
		// FIXME: Handle converstion error
		idx, _ := strconv.Atoi(to[1])
		if len(b.SuitPile[idx]) == 1 {
			if pc.CardDetails.Rank != card.Ranks()[0] {
				fmt.Println("FIRST CARD FROM THE SUIT PILE MUST BE A(1)")
				return false
			}
			break
		}
		if pc.CardDetails.Rank != destinationCard.Rank+1 || pc.CardDetails.Suit != destinationCard.Suit {
			fmt.Println("CARD CAN'T BE PLACED INTO SUIT PILE ", idx)
			return false
		}
	}

	fmt.Println("MOVING FROM: ", from, " TO: ", to)
	var c *card.Card
	switch {
	case from[0] == "c":
		// FIXME: Handle converstion error
		idx, _ := strconv.Atoi(from[1])
		c = b.Columns[idx][len(b.Columns[idx])-1]
		c.IsBeingUsed = false
		b.Columns[idx] = b.Columns[idx][:len(b.Columns[idx])-1]
		b.Columns[idx][len(b.Columns[idx])-1].IsFlippedDown = false
	case from[0] == "s":
		// FIXME: Handle converstion error
		idx, _ := strconv.Atoi(from[1])
		c = b.SuitPile[idx][len(b.SuitPile[idx])-1]
		c.IsBeingUsed = false
		b.SuitPile[idx] = b.SuitPile[idx][:len(b.SuitPile[idx])-1]
	case pc.OriginalPile == DiscardPosition:
		if len(b.DiscardPile) > 1 {
			c = b.DiscardPile[0]
			c.IsBeingUsed = false
			b.DiscardPile = b.DiscardPile[1:]
		}
	}
	switch {
	case to[0] == "c":
		// FIXME: Handle converstion error
		idx, _ := strconv.Atoi(to[1])
		last := b.Columns[idx][len(b.Columns[idx])-1]
		c.Frame.X, c.Frame.Y = last.Frame.X, last.Frame.Y
		if len(b.Columns[idx]) > 1 {
			c.Frame.Y += card.Spacing
		}
		b.Columns[idx] = append(b.Columns[idx], c)
	case to[0] == "s":
		// FIXME: Handle converstion error
		idx, _ := strconv.Atoi(to[1])
		last := b.SuitPile[idx][len(b.SuitPile[idx])-1]
		c.Frame.X, c.Frame.Y = last.Frame.X, last.Frame.Y
		b.SuitPile[idx] = append(b.SuitPile[idx], c)
	}
	return true
}
