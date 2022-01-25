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
			deck = append(deck, card.New(rank, suit, 6*card.Width, 0, texture))
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

	board := board.New(emptyTexture, backTexture, deck)

	running := true
	deckCard := card.New(-1, "-1", 6*card.Width, 0, nil)
	shouldMove := false
	var originalX, originalY int32
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				fmt.Println("Closing Solitaire!")
				running = false
			case *sdl.MouseMotionEvent:
				fmt.Println("Mouse", t.Which, "moved by", t.XRel, t.YRel, "at", t.X, t.Y)
				if shouldMove {
					c := board.DiscardPile[0]
					c.Frame.X, c.Frame.Y = t.X, t.Y
				}
			case *sdl.MouseButtonEvent:
				if t.State == sdl.RELEASED {
					if shouldMove {
						shouldMove = false
						c := board.DiscardPile[0]
						c.Frame.X, c.Frame.Y = originalX, originalY
					}

					if t.X > 6*card.Width && t.X < 7*card.Width && t.Y > 0 && t.Y < card.Height {
						fmt.Println("Mouse", t.Which, "button", t.Button, "released at", t.X, t.Y)
						if len(board.DrawPile) > 0 {
							c := board.DrawPile[len(board.DrawPile)-1]
							c.Frame.X = 5 * card.Width

							board.DrawPile = board.DrawPile[:len(board.DrawPile)-1]
							board.DiscardPile = append([]*card.Card{c}, board.DiscardPile...)
						} else {
							board.DrawPile = board.DiscardPile
							board.DiscardPile = []*card.Card{}
						}
					}
				}
				if t.State == sdl.PRESSED {
					if len(board.DiscardPile) > 0 {
						c := board.DiscardPile[0]
						if t.X > c.Frame.X && t.X < c.Frame.X+c.Frame.W && t.Y > c.Frame.Y && t.Y < c.Frame.Y+c.Frame.H && !shouldMove {
							shouldMove = true
							originalX, originalY = c.Frame.X, c.Frame.Y
						}
					}
				}
			}
		}
		rw.Clear()
		// Render Suite Pile
		for _, cards := range board.SuitPile {
			for _, c := range cards {
				rw.Render(c)
			}
		}
		// Render Columns
		for _, cards := range board.Columns {
			for _, c := range cards {
				if c.IsFlippedDown {
					originalTexture := c.Texture
					c.Texture = board.BackCardTexture
					rw.Render(c)
					c.Texture = originalTexture
				} else {
					rw.Render(c)
				}
			}
		}
		// Render Drwa Pile
		if len(board.DrawPile) > 0 {
			deckCard.Texture = backTexture
		} else {
			deckCard.Texture = emptyTexture
		}
		rw.Render(deckCard)
		// Render Discard Pile
		if len(board.DiscardPile) > 1 && shouldMove {
			rw.Render(board.DiscardPile[1])
			rw.Render(board.DiscardPile[0])
		} else if len(board.DiscardPile) > 0 {
			rw.Render(board.DiscardPile[0])
		}
		rw.Display()
	}
}
