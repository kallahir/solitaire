package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"strings"

	"github.com/kallahir/solitaire/board"
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

	resources, err := ioutil.ReadDir("resources/cards")
	if err != nil {
		panic(err)
	}

	textures := make(map[string]*sdl.Texture)
	for _, file := range resources {
		if file.IsDir() {
			continue
		}
		texture, err := rw.LoadTexture(fmt.Sprintf("resources/cards/%s", file.Name()))
		if err != nil {
			panic(err)
		}
		fileName, err := RemoveFileExtension(file)
		if err != nil {
			panic(err)
		}
		textures[fileName] = texture
	}

	game := board.New(rw, textures)
	var x, y int32
	for game.IsRunning {
		for event := sdl.WaitEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				fmt.Println("Closing Solitaire!", t)
				game.IsRunning = false
			case *sdl.MouseMotionEvent:
				x, y = t.X, t.Y
			case *sdl.MouseButtonEvent:
				game.HandleClick(t.X, t.Y, t.State)
			}
		}
		rw.Clear()
		game.Render(rw, x, y)
		rw.Display()
		sdl.Delay(16)
	}
}

func RemoveFileExtension(file fs.FileInfo) (string, error) {
	result := strings.Split(file.Name(), ".")
	if len(result) < 2 || len(result) > 2 {
		return "", fmt.Errorf("can't remove file extension from %s", file.Name())
	}
	return result[0], nil
}
