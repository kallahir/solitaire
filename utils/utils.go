package utils

import (
	"fmt"
	"strings"

	"github.com/veandco/go-sdl2/sdl"
)

func RemoveFileExtension(filename string) (string, error) {
	result := strings.Split(filename, ".")
	if len(result) < 2 || len(result) > 2 {
		return "", fmt.Errorf("can't remove file extension from %s", filename)
	}
	return result[0], nil
}

func CheckCollision(x, y int32, boudingBox *sdl.Rect) bool {
	if x > boudingBox.X && x < boudingBox.X+boudingBox.W && y > boudingBox.Y && y < boudingBox.Y+boudingBox.H {
		return true
	}
	return false
}
