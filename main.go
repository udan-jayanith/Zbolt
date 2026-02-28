package main

import (
	"fmt"
	"image"
	"os"

	gui "github.com/guigui-gui/guigui"
	"github.com/hajimehoshi/ebiten/v2"

	"API-Client/widgets/request/page"
	Welcome "API-Client/widgets/welcome"
)

type Root struct {
	gui.DefaultWidget
	
	welcome_page_widget Welcome.Welcome
	request_page_widget request_page.RequestPage
}

func (r *Root) Build(context *gui.Context, adder *gui.ChildAdder) error {
	adder.AddWidget(&r.request_page_widget)
	return nil
}

func (r *Root) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layouter.LayoutWidget(&r.request_page_widget, widgetBounds.Bounds())
}

func main() {
	op := &gui.RunOptions{
		Title: "API Client",
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