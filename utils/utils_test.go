package utils

import (
	"fmt"
	"testing"

	"github.com/veandco/go-sdl2/sdl"
)

func TestRemoveFileExtension(t *testing.T) {
	var tests = []struct {
		filename string
		want     string
		wantErr  bool
	}{
		{
			filename: "file01.jpg",
			want:     "file01",
			wantErr:  false,
		},
		{
			filename: "file02.out.jpg",
			want:     "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			name, err := RemoveFileExtension(tt.filename)
			if err != nil && !tt.wantErr {
				t.Errorf("error: %v", err)
			}
			if name != tt.want {
				t.Errorf("got %s, want %s", name, tt.want)
			}
		})
	}
}

func TestCheckCollision(t *testing.T) {
	var tests = []struct {
		x, y int32
		bb   *sdl.Rect
		want bool
	}{
		{
			x:    1,
			y:    1,
			bb:   &sdl.Rect{X: 0, Y: 0, H: 10, W: 10},
			want: true,
		},
		{
			x:    1,
			y:    1,
			bb:   &sdl.Rect{X: 5, Y: 5, H: 10, W: 10},
			want: false,
		},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%d/%d check collision %v", tt.x, tt.y, tt.bb)
		t.Run(testname, func(t *testing.T) {
			result := CheckCollision(tt.x, tt.y, tt.bb)
			if result != tt.want {
				t.Errorf("got %t, want %t", result, tt.want)
			}
		})
	}
}
