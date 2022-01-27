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
	// deckCard := card.New(-1, "-1", 6*card.Width, 0, nil)
	// var shouldMove bool
	// var pcCards []*card.Card
	for game.IsRunning {
		for event := sdl.WaitEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				fmt.Println("Closing Solitaire!", t)
				game.IsRunning = false
				// case *sdl.MouseMotionEvent:
				// 	// gameBoard.CheckPosition(t.X, t.Y)
				// 	if shouldMove {
				// 		pc := pcCards[len(pcCards)-1]
				// 		pc.Frame.X, pc.Frame.Y = t.X-card.Width/2, t.Y-card.Height/2
				// 	}
				// case *sdl.MouseButtonEvent:
				// 	// Make the CheckPosition method return an array of cards instead of a single card
				// 	boardPosition, cards := gameBoard.CheckPosition(t.X, t.Y)
				// 	if t.State == sdl.RELEASED {
				// 		if boardPosition == board.DrawPosition && !shouldMove {
				// 			gameBoard.DrawCard()
				// 		}
				// 		if shouldMove && len(cards) > 0 {
				// 			fmt.Println("CARD: ", pcCards[len(cards)-1], " DROPPED AT: ", boardPosition, " OVER CARD:", cards[0])
				// 			shouldMove = false
				// 			if !gameBoard.MoveCard(pcCards[len(pcCards)-1], cards[0], boardPosition) {
				// 				pc := pcCards[len(pcCards)-1]
				// 				pc.Frame.X, pc.Frame.Y = pc.OriginalX, pc.OriginalY
				// 				pc.IsBeingUsed = false
				// 				pc = nil
				// 				pcCards = pcCards[:len(pcCards)-1]
				// 			}
				// 		} else if shouldMove {
				// 			shouldMove = false
				// 			pc := pcCards[len(pcCards)-1]
				// 			pc.Frame.X, pc.Frame.Y = pc.OriginalX, pc.OriginalY
				// 			pc.IsBeingUsed = false
				// 			pc = nil
				// 			pcCards = pcCards[:len(pcCards)-1]
				// 		}
				// 	}
				// 	if t.State == sdl.PRESSED && boardPosition != "" {
				// 		if boardPosition != board.DrawPosition && !shouldMove && len(cards) > 0 && cards[0].Rank != -1 && cards[0].Suit != "-1" {
				// 			shouldMove = true
				// 			cards[0].IsBeingUsed = true
				// 			cards[0].OriginalPile = boardPosition
				// 			cards[0].OriginalX = cards[0].Frame.X
				// 			cards[0].OriginalY = cards[0].Frame.Y
				// 			pcCards = append(pcCards, cards[0])
				// 		}
				// 	}
			}
		}
		rw.Clear()
		game.Render(rw)
		// Render Suite Pile
		// for _, cards := range gameBoard.SuitPile {
		// 	for _, c := range cards {
		// 		rw.Render(c)
		// 	}
		// }
		// // Render Columns
		// for _, cards := range gameBoard.Columns {
		// 	for _, c := range cards {
		// 		if c.IsFlippedDown {
		// 			originalTexture := c.Texture
		// 			c.Texture = gameBoard.BackCardTexture
		// 			rw.Render(c)
		// 			c.Texture = originalTexture
		// 		} else {
		// 			rw.Render(c)
		// 		}
		// 	}
		// }
		// // Render Draw Pile
		// if len(gameBoard.DrawPile) > 0 {
		// 	deckCard.Texture = backTexture
		// } else {
		// 	deckCard.Texture = emptyTexture
		// }
		// rw.Render(deckCard)
		// // Render Discard Pile
		// if len(gameBoard.DiscardPile) > 1 && shouldMove {
		// 	rw.Render(gameBoard.DiscardPile[1])
		// 	rw.Render(gameBoard.DiscardPile[0])
		// } else if len(gameBoard.DiscardPile) > 0 {
		// 	rw.Render(gameBoard.DiscardPile[0])
		// }
		// // Render PlayingCard
		// if len(pcCards) > 0 {
		// 	for i := range pcCards {
		// 		rw.Render(pcCards[i])
		// 	}
		// }
		rw.Display()
		sdl.Delay(16)
	}
}

func RemoveFileExtension(file fs.FileInfo) (string, error) {
	result := strings.Split(file.Name(), ".")
	if len(result) < 2 || len(result) > 2 {
		return "", fmt.Errorf("Can't remove file extension from %s", file.Name())
	}
	return result[0], nil
}
