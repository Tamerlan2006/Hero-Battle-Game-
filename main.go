package main

import (
	"log"

	gameui "week67/GameUI"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	game := gameui.NewGame()
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Hero Battle Game")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

