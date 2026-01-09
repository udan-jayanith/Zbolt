package main

import (
	"fmt"
	"os"
	"image"

	gui "github.com/guigui-gui/guigui"
	"github.com/hajimehoshi/ebiten/v2"

	WelcomePage "API-Client/widgets/welcome-page"
)

func main() {
	op := &gui.RunOptions{
		Title: "API Client",
		WindowMinSize: image.Pt(700, 444),
		RunGameOptions: &ebiten.RunGameOptions{
			ApplePressAndHoldEnabled: true,
		},
	}
	if err := gui.Run(&WelcomePage.WelcomePage{}, op); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
