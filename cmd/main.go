package main

import (
	"fmt"
	"log"
	"toy/internal/game"
	"toy/internal/spine"
)

func main() {
	spineData, err := spine.LoadSpineData("./data/spine/hero-ess.json")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(spineData.Skeleton.Hash)

	spineData.SetAnimation("idle")

	g, err := game.NewGame(game.WithScreenSize(800, 600),
		game.WithTitle("toy"),
		game.WithSpine(spineData),
	)
	if err != nil {
		log.Fatal(err)
	}
	if err := g.Run(); err != nil {
		log.Fatal(err)
	}
}
