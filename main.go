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
	"API-Client/widgets/inspect"
	request_page "API-Client/widgets/request/page"
)

type Root struct {
	gui.DefaultWidget

	welcome_page_widget home.HomePage
	request_page_widget request_page.RequestPage
	inspect_widget inspect.InspectWidget
}

func (r *Root) Build(context *gui.Context, adder *gui.ChildAdder) error {
	adder.AddWidget(&r.request_page_widget)
	
	r.inspect_widget.SetOpen(true)
	adder.AddWidget(&r.inspect_widget)
	return nil
}

func (r *Root) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	b := widgetBounds.Bounds()
	
	layouter.LayoutWidget(&r.request_page_widget, b)
	
	layouter.LayoutWidget(&r.inspect_widget, image.Rectangle{
		Min: image.Point{
			X: b.Min.X,
			Y: b.Max.Y/2,
		},
		Max: image.Point{
			X: b.Max.X,
			Y: b.Max.Y,
		},
	})
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
		Title:         "Zbolt",
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