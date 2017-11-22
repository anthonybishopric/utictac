package main

import (
	"github.com/anthonybishopric/utictac/game"
)

func main() {
	runner := game.NewRunner()
	runner.Run(make(chan struct{}))
}
