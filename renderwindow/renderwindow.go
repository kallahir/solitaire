package renderwindow

import (
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

func (rw *RenderWindow) Render(frame *sdl.Rect, texture *sdl.Texture) {
	rw.Renderer.Copy(texture, nil, frame)
}

func (rw *RenderWindow) Display() {
	rw.Renderer.Present()
}
