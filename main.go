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

	// grassTexture, err := rw.LoadTexture("../resources/gfx/ground_grass_1.png")
	// if err != nil {
	// 	panic(err)
	// }

	// entity := entity.NewEntity(0, 0, grassTexture)

	cardTexture, err := rw.LoadTexture("../resources/gfx/card.gif")
	if err != nil {
		panic(err)
	}

	var cards []*entity.Entity
	for i := 0; i < 7; i++ {
		if i == 4 {
			continue
		}
		cards = append(cards, entity.NewEntity(int32(i)*116, 0, cardTexture))
	}

	running := true
	// shouldMove := false
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			// switch t := event.(type) {
			switch event.(type) {
			case *sdl.QuitEvent:
				fmt.Println("Closing Solitaire!")
				running = false
				// case *sdl.MouseMotionEvent:
				// 	if shouldMove {
				// 		entity.X, entity.Y = t.X, t.Y
				// 	}
				// case *sdl.MouseButtonEvent:
				// 	if t.State == sdl.PRESSED {
				// 		if shouldMove {
				// 			shouldMove = false
				// 		} else {
				// 			shouldMove = true
				// 		}
				// 	}
			}
		}
		rw.Clear()
		for _, card := range cards {
			rw.Render(card)
		}
		rw.Display()
	}
}
