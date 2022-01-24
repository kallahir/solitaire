package main

import (
	"fmt"
	"math/rand"
	"time"

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

	// cardTexture, err := rw.LoadTexture("../resources/gfx/card.gif")
	// if err != nil {
	// 	panic(err)
	// }

	// var cards []*entity.Entity
	// var spacing int32 = 30
	// for i := 0; i < 7; i++ {
	// 	for j := 0; j <= i; j++ {
	// 		if i == j {
	// 			cards = append(cards, entity.NewEntity(int32(i)*116, 176+(int32(j)*spacing+spacing), cardTexture))
	// 		} else {
	// 			cards = append(cards, entity.NewEntity(int32(i)*116, 176+(int32(j)*spacing+spacing), backTexture))
	// 		}
	// 	}
	// }

	// for i := 0; i < 7; i++ {
	// 	switch {
	// 	case i == 4:
	// 		continue
	// 	case i < 4:
	// 		cards = append(cards, entity.NewEntity(int32(i)*116, 0, placeholderTexture))
	// 	case i == 5:
	// 		cards = append(cards, entity.NewEntity(int32(i)*116, 0, cardTexture))
	// 	case i == 6:
	// 		cards = append(cards, entity.NewEntity(int32(i)*116, 0, backTexture))
	// 	}
	// }

	var drawPile []*card.Card
	for rank := int32(1); rank <= 13; rank++ {
		for _, suit := range []string{"s", "h", "d", "c"} {
			texture, err := rw.LoadTexture(fmt.Sprintf("../resources/cards/%02d%s.gif", rank, suit))
			if err != nil {
				panic(err)
			}
			drawPile = append(drawPile, card.New(rank, suit, 6*card.Width, 0, texture))
		}
	}

	// Shuffle Cards
	for i := 0; i < 100; i++ {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(drawPile), func(i, j int) { drawPile[i], drawPile[j] = drawPile[j], drawPile[i] })
	}

	backTexture, err := rw.LoadTexture("../resources/cards/back.gif")
	if err != nil {
		panic(err)
	}

	placeholderTexture, err := rw.LoadTexture("../resources/cards/placeholder.gif")
	if err != nil {
		panic(err)
	}

	running := true
	var discardPile []*card.Card
	deckCard := card.New(-1, "-1", 6*card.Width, 0, nil)
	// shouldMove := false
	// playingCard := cards[len(cards)-2]
	// originalX, originalY := playingCard.X, playingCard.Y
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				fmt.Println("Closing Solitaire!")
				running = false
			// case *sdl.MouseMotionEvent:
			// 	fmt.Println("Mouse", t.Which, "moved by", t.XRel, t.YRel, "at", t.X, t.Y)
			// 	// if shouldMove {
			// 	// 	playingCard.X, playingCard.Y = t.X, t.Y
			// 	// }
			case *sdl.MouseButtonEvent:
				if t.State == sdl.RELEASED {
					if t.X > 6*card.Width && t.X < 7*card.Width && t.Y > 0 && t.Y < card.Height {
						fmt.Println("Mouse", t.Which, "button", t.Button, "released at", t.X, t.Y)
						if len(drawPile) > 0 {
							c := drawPile[len(drawPile)-1]
							c.Frame.X = 5 * card.Width

							drawPile = drawPile[:len(drawPile)-1]
							discardPile = append([]*card.Card{c}, discardPile...)
						} else {
							drawPile = discardPile
							discardPile = []*card.Card{}
						}
					}
				}
				// if t.X > playingCard.X && t.X < playingCard.X+playingCard.CurrentFrame.W && t.Y > playingCard.Y && t.Y < playingCard.Y+playingCard.CurrentFrame.H && !shouldMove {
				// 	shouldMove = true
				// } else {
				// 	shouldMove = false
				// 	playingCard.X, playingCard.Y = originalX, originalY
				// }
			}
		}
		rw.Clear()
		if len(drawPile) > 0 {
			deckCard.Texture = backTexture
		} else {
			deckCard.Texture = placeholderTexture
		}
		rw.Render(deckCard)
		if len(discardPile) > 0 {
			rw.Render(discardPile[0])
		}
		// for _, card := range deck {
		// 	rw.Render(card)
		// }
		rw.Display()
	}
}
