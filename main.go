package main

import (
	"fmt"
	"image"
	"os"

	gui "github.com/guigui-gui/guigui"
	"github.com/hajimehoshi/ebiten/v2"

	 "API-Client/widgets/request-page"
)

func main() {
	op := &gui.RunOptions{
		Title: "API Client",
		WindowMinSize: image.Pt(800, 444),
		RunGameOptions: &ebiten.RunGameOptions{
			ApplePressAndHoldEnabled: true,
		},
	}
	if err := gui.Run(&Requester.RequestPage{}, op); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}