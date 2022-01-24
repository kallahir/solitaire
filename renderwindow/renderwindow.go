package renderwindow

import (
	"github.com/kallahir/solitaire/entity"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type RenderWindow struct {
	Window   *sdl.Window
	Renderer *sdl.Renderer
}

func NewRenderWindow(title string, w, h int32) (*RenderWindow, error) {
	window, err := sdl.CreateWindow(title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, w, h, sdl.WINDOW_SHOWN)
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

func (rw *RenderWindow) Render(entity *entity.Entity) {
	dst := &sdl.Rect{X: entity.X, Y: entity.Y, W: entity.CurrentFrame.W, H: entity.CurrentFrame.H}
	rw.Renderer.Copy(entity.Texture, nil, dst)
}

func (rw *RenderWindow) Display() {
	rw.Renderer.Present()
}
