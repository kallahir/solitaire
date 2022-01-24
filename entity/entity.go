package entity

import "github.com/veandco/go-sdl2/sdl"

type Entity struct {
	X            int32
	Y            int32
	CurrentFrame *sdl.Rect
	Texture      *sdl.Texture
}

func NewEntity(x, y int32, texture *sdl.Texture) *Entity {
	return &Entity{X: x, Y: y, CurrentFrame: &sdl.Rect{X: x, Y: y, W: 116, H: 176}, Texture: texture}
}
