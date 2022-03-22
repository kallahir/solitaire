package renderwindow

import (
	"embed"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type RenderWindow struct {
	Window   *sdl.Window
	Renderer *sdl.Renderer
}

func New(title string, w, h int32) (*RenderWindow, error) {
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

func (rw *RenderWindow) LoadTextureFromEmbedFS(fs embed.FS, path string) (*sdl.Texture, error) {
	file, err := fs.ReadFile(path)
	if err != nil {
		return nil, err
	}
	fileRW, err := sdl.RWFromMem(file)
	if err != nil {
		return nil, err
	}
	texture, err := img.LoadTextureRW(rw.Renderer, fileRW, true)
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
