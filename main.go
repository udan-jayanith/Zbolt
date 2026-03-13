package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	_ "image/png"
	"log"
	"os"

	gui "github.com/guigui-gui/guigui"
	"github.com/hajimehoshi/ebiten/v2"

	home "API-Client/widgets/home"
	request_page "API-Client/widgets/request/page"
)

type Root struct {
	gui.DefaultWidget

	welcome_page_widget home.HomePage
	request_page_widget request_page.RequestPage
}

func (r *Root) Build(context *gui.Context, adder *gui.ChildAdder) error {
	adder.AddWidget(&r.welcome_page_widget)
	return nil
}

func (r *Root) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layouter.LayoutWidget(&r.welcome_page_widget, widgetBounds.Bounds())
}

//go:embed icon.png
var zbolt_icon_bytes []byte

func main() {
	zebolt_icon, _, err := image.Decode(bytes.NewReader(zbolt_icon_bytes))
	if err != nil {
		log.Fatal(err.Error())
	}
	ebiten.SetWindowIcon([]image.Image{zebolt_icon})

	op := &gui.RunOptions{
		Title:         "API Client",
		WindowMinSize: image.Pt(800, 444),
		RunGameOptions: &ebiten.RunGameOptions{
			ApplePressAndHoldEnabled: true,
		},
	}
	if err := gui.Run(&Root{}, op); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
