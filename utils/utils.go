package utils

import (
	"fmt"
	"io/fs"
	"strings"

	"github.com/veandco/go-sdl2/sdl"
)

func RemoveFileExtension(file fs.FileInfo) (string, error) {
	result := strings.Split(file.Name(), ".")
	if len(result) < 2 || len(result) > 2 {
		return "", fmt.Errorf("can't remove file extension from %s", file.Name())
	}
	return result[0], nil
}

func CheckCollision(x, y int32, boudingBox *sdl.Rect) bool {
	if x > boudingBox.X && x < boudingBox.X+boudingBox.W && y > boudingBox.Y && y < boudingBox.Y+boudingBox.H {
		return true
	}
	return false
}
