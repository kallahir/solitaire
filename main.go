package main

import (
	"fmt"

	"github.com/kallahir/solitaire/renderwindow"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	fmt.Println("Welcome to Solitaire!")

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	rw, err := renderwindow.NewRenderWindow("Solitaire", 1280, 720)
	if err != nil {
		panic(err)
	}
	defer rw.CleanUp()

	grassTexture, err := rw.LoadTexture("../resources/gfx/ground_grass_1.png")
	if err != nil {
		panic(err)
	}

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				fmt.Println("Closing Solitaire!")
				running = false
			}
		}
		rw.Clear()
		rw.Render(grassTexture)
		rw.Display()
	}
}
