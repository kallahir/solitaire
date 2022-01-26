package main

import (
	"fmt"

	"github.com/kallahir/solitaire/board"
	"github.com/kallahir/solitaire/card"
	"github.com/kallahir/solitaire/renderwindow"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	fmt.Println("Welcome to Solitaire!")

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	rw, err := renderwindow.New("Solitaire", 812, 900)
	if err != nil {
		panic(err)
	}
	defer rw.CleanUp()

	var deck []*card.Card
	for _, rank := range card.Ranks() {
		for _, suit := range card.Suits() {
			texture, err := rw.LoadTexture(fmt.Sprintf("../resources/cards/%02d%s.gif", rank, suit))
			if err != nil {
				panic(err)
			}
			deck = append(deck, card.New(rank, suit, 5*card.Width, 0, texture))
		}
	}

	backTexture, err := rw.LoadTexture("../resources/cards/back.gif")
	if err != nil {
		panic(err)
	}

	emptyTexture, err := rw.LoadTexture("../resources/cards/empty.gif")
	if err != nil {
		panic(err)
	}

	running := true
	gameBoard := board.New(emptyTexture, backTexture, deck)
	deckCard := card.New(-1, "-1", 6*card.Width, 0, nil)
	shouldMove := false
	// playingCards := []*card.PlayingCard{}
	pc := new(card.PlayingCard)
	for running {
		for event := sdl.WaitEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				fmt.Println("Closing Solitaire!")
				running = false
			case *sdl.MouseMotionEvent:
				// gameBoard.CheckPosition(t.X, t.Y)
				if shouldMove {
					pc.CardDetails.Frame.X, pc.CardDetails.Frame.Y = t.X-card.Width/2, t.Y-card.Height/2
				}
			case *sdl.MouseButtonEvent:
				// Make the CheckPosition method return an array of cards instead of a single card
				boardPosition, cards := gameBoard.CheckPosition(t.X, t.Y)
				if t.State == sdl.RELEASED {
					if boardPosition == board.DrawPosition && !shouldMove {
						gameBoard.DrawCard()
					}
					if shouldMove && len(cards) > 0 {
						fmt.Println("CARD: ", pc.CardDetails, " DROPPED AT: ", boardPosition, " OVER CARD:", cards[0])
						shouldMove = false
						if !gameBoard.MoveCard(pc, cards[0], boardPosition) {
							pc.CardDetails.Frame.X, pc.CardDetails.Frame.Y = pc.OriginalX, pc.OriginalY
							pc.CardDetails.IsBeingUsed = false
							pc.CardDetails = nil
						}
					}
				}
				if t.State == sdl.PRESSED && boardPosition != "" {
					if boardPosition != board.DrawPosition && !shouldMove && len(cards) > 0 && cards[0].Rank != -1 && cards[0].Suit != "-1" {
						shouldMove = true
						pc.CardDetails = cards[0]
						pc.CardDetails.IsBeingUsed = true
						pc.OriginalPile = boardPosition
						pc.OriginalX, pc.OriginalY = cards[0].Frame.X, cards[0].Frame.Y
					}
				}
			}
		}
		rw.Clear()
		// Render Suite Pile
		for _, cards := range gameBoard.SuitPile {
			for _, c := range cards {
				rw.Render(c)
			}
		}
		// Render Columns
		for _, cards := range gameBoard.Columns {
			for _, c := range cards {
				if c.IsFlippedDown {
					originalTexture := c.Texture
					c.Texture = gameBoard.BackCardTexture
					rw.Render(c)
					c.Texture = originalTexture
				} else {
					rw.Render(c)
				}
			}
		}
		// Render Draw Pile
		if len(gameBoard.DrawPile) > 0 {
			deckCard.Texture = backTexture
		} else {
			deckCard.Texture = emptyTexture
		}
		rw.Render(deckCard)
		// Render Discard Pile
		if len(gameBoard.DiscardPile) > 1 && shouldMove {
			rw.Render(gameBoard.DiscardPile[1])
			rw.Render(gameBoard.DiscardPile[0])
		} else if len(gameBoard.DiscardPile) > 0 {
			rw.Render(gameBoard.DiscardPile[0])
		}
		// Render PlayingCard
		if pc.CardDetails != nil {
			rw.Render(pc.CardDetails)
		}
		rw.Display()
		sdl.Delay(16)
	}
}
