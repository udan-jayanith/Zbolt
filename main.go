package main

import (
	"fmt"
	"image"
	"os"

	gui "github.com/guigui-gui/guigui"
	"github.com/hajimehoshi/ebiten/v2"

	 "API-Client/widgets/requester"
)

func main() {
	op := &gui.RunOptions{
		Title: "API Client",
		WindowMinSize: image.Pt(700, 444),
		RunGameOptions: &ebiten.RunGameOptions{
			ApplePressAndHoldEnabled: true,
		},
	}
	if err := gui.Run(&Requester.Requester{}, op); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}