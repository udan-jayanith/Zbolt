package main

import (
	"fmt"
	"image"
	"os"

	gui "github.com/guigui-gui/guigui"
	"github.com/hajimehoshi/ebiten/v2"

	RequestPage "API-Client/widgets/request-pages"
//	WelcomePage "API-Client/widgets/welcome-page"
)

func main() {
	op := &gui.RunOptions{
		Title: "API Client",
		WindowMinSize: image.Pt(700, 444),
		RunGameOptions: &ebiten.RunGameOptions{
			ApplePressAndHoldEnabled: true,
		},
	}
	if err := gui.Run(&RequestPage.BasicRequestPage{}, op); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
