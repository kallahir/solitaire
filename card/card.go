package card

import (
	"fmt"
)

const (
	Width    int32  = 116
	Height   int32  = 176
	VSpacing int32  = 50
	HSpacing int32  = 15
	Empty    string = "empty"
	Back     string = "back"
)

type Card struct {
	Rank          int32
	Suit          string
	TextureKey    string
	IsFlippedDown bool
	IsSpecialCard bool
}

func New(rank int32, suit string) *Card {
	var textureKey string
	if suit == Empty || suit == Back {
		textureKey = suit
	} else {
		textureKey = fmt.Sprintf("%02d%s", rank, suit)
	}
	return &Card{Rank: rank, Suit: suit, TextureKey: textureKey}
}

func Suits() []string {
	return []string{"s", "h", "d", "c"}
}

func Ranks() []int32 {
	return []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}
}

func (c *Card) CompareOverlappingSuit(dst *Card) bool {
	switch {
	case c.Suit == "s" || c.Suit == "c":
		if dst.Suit == "s" || dst.Suit == "c" {
			return false
		}
	case c.Suit == "h" || c.Suit == "d":
		if dst.Suit == "h" || dst.Suit == "d" {
			return false
		}
	}
	return true
}
