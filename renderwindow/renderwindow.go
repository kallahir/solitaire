package renderwindow

import (
	"github.com/kallahir/solitaire/card"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type RenderWindow struct {
	Window   *sdl.Window
	Renderer *sdl.Renderer
}

func New(title string, w, h int32) (*RenderWindow, error) {
	window, err := sdl.CreateWindow(title, 0, 0, w, h, sdl.WINDOW_SHOWN)
	if err != nil {
		return nil, err
	}
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return nil, err
	}
	renderer.SetDrawColor(82, 115, 68, 0)
	renderer.Clear()
	return &RenderWindow{Window: window, Renderer: renderer}, nil
}

func (rw *RenderWindow) CleanUp() {
	rw.Window.Destroy()
}

func (rw *RenderWindow) LoadTexture(path string) (*sdl.Texture, error) {
	texture, err := img.LoadTexture(rw.Renderer, path)
	if err != nil {
		return nil, err
	}
	return texture, nil
}

func (rw *RenderWindow) Clear() {
	rw.Renderer.Clear()
}

func (rw *RenderWindow) Render(card *card.Card) {
	dst := &sdl.Rect{X: card.Frame.X, Y: card.Frame.Y, W: card.Frame.W, H: card.Frame.H}
	rw.Renderer.Copy(card.Texture, nil, dst)
}

func (rw *RenderWindow) Display() {
	rw.Renderer.Present()
}
