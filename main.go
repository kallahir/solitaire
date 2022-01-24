package main

import (
	"fmt"

	"github.com/kallahir/solitaire/entity"
	"github.com/kallahir/solitaire/renderwindow"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	fmt.Println("Welcome to Solitaire!")

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	rw, err := renderwindow.NewRenderWindow("Solitaire", 812, 900)
	if err != nil {
		panic(err)
	}
	defer rw.CleanUp()

	cardTexture, err := rw.LoadTexture("../resources/gfx/card.gif")
	if err != nil {
		panic(err)
	}

	backTexture, err := rw.LoadTexture("../resources/gfx/back.gif")
	if err != nil {
		panic(err)
	}

	placeholderTexture, err := rw.LoadTexture("../resources/gfx/placeholder.gif")
	if err != nil {
		panic(err)
	}

	var cards []*entity.Entity
	var spacing int32 = 30
	for i := 0; i < 7; i++ {
		for j := 0; j <= i; j++ {
			if i == j {
				cards = append(cards, entity.NewEntity(int32(i)*116, 176+(int32(j)*spacing+spacing), cardTexture))
			} else {
				cards = append(cards, entity.NewEntity(int32(i)*116, 176+(int32(j)*spacing+spacing), backTexture))
			}
		}
	}

	for i := 0; i < 7; i++ {
		switch {
		case i == 4:
			continue
		case i < 4:
			cards = append(cards, entity.NewEntity(int32(i)*116, 0, placeholderTexture))
		case i == 5:
			cards = append(cards, entity.NewEntity(int32(i)*116, 0, cardTexture))
		case i == 6:
			cards = append(cards, entity.NewEntity(int32(i)*116, 0, backTexture))
		}
	}

	running := true
	shouldMove := false
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				fmt.Println("Closing Solitaire!")
				running = false
			case *sdl.MouseMotionEvent:
				if shouldMove {
					cards[len(cards)-2].X, cards[len(cards)-2].Y = t.X, t.Y
				}
			case *sdl.MouseButtonEvent:
				if shouldMove {
					shouldMove = false
				} else {
					shouldMove = true
				}
			}
		}
		rw.Clear()
		for _, card := range cards {
			rw.Render(card)
		}
		rw.Display()
	}
}
