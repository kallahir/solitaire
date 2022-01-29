package main

import (
	"fmt"
	"io/ioutil"

	"github.com/kallahir/solitaire/board"
	"github.com/kallahir/solitaire/card"
	"github.com/kallahir/solitaire/renderwindow"
	"github.com/kallahir/solitaire/utils"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	fmt.Println("Welcome to Solitaire!")

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	rw, err := renderwindow.New("Solitaire", board.NumberOfColumns*card.Width+(board.NumberOfColumns+1)*card.HSpacing, 1024)
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
		fileName, err := utils.RemoveFileExtension(file)
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
				fmt.Println("Closing Solitaire!")
				game.IsRunning = false
			case *sdl.MouseMotionEvent:
				x, y = t.X, t.Y
			case *sdl.MouseButtonEvent:
				game.HandleClick(t.X, t.Y, t.State)
			}
		}
		rw.Clear()
		game.Render(rw, x, y)
		if game.IsOver() {
			switch GetUserInput(rw.Window) {
			case -1:
				fmt.Println("No selection...")
			case 1:
				fmt.Println("Restarting game...")
				game = board.New(rw, textures)
			case 2:
				fmt.Println("Closing Solitaire!")
				game.IsRunning = false
			}
		}
		rw.Display()
		sdl.Delay(16)
	}
}

func GetUserInput(w *sdl.Window) int32 {
	buttons := []sdl.MessageBoxButtonData{
		{Flags: sdl.MESSAGEBOX_BUTTON_RETURNKEY_DEFAULT, ButtonID: 1, Text: "Yes"},
		{Flags: sdl.MESSAGEBOX_BUTTON_ESCAPEKEY_DEFAULT, ButtonID: 2, Text: "Quit"},
	}

	colorScheme := sdl.MessageBoxColorScheme{
		Colors: [5]sdl.MessageBoxColor{
			{R: 255, G: 0, B: 0},
			{R: 0, G: 255, B: 0},
			{R: 255, G: 255, B: 0},
			{R: 0, G: 0, B: 255},
			{R: 255, G: 0, B: 255},
		},
	}

	messageboxdata := sdl.MessageBoxData{
		Flags:       sdl.MESSAGEBOX_INFORMATION,
		Window:      w,
		Title:       "Congratulations!\nYou Won!",
		Message:     "Do you want to play another match?",
		Buttons:     buttons,
		ColorScheme: &colorScheme,
	}

	var buttonid int32
	var err error
	if buttonid, err = sdl.ShowMessageBox(&messageboxdata); err != nil {
		panic(err)
	}

	return buttonid
}
