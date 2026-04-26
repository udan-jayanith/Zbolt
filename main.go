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

	message_model "API-Client/message-model"
	home "API-Client/widgets/home"
	request_page "API-Client/widgets/request/page"
)

type Root struct {
	gui.DefaultWidget

	message_model_widget gui.Widget
	welcome_page_widget  home.HomePage
	request_page_widget  request_page.RequestPage
}

func (r *Root) Build(context *gui.Context, adder *gui.ChildAdder) error {
	if r.message_model_widget == nil {
		r.message_model_widget = &message_model.MessageModel
	}
	adder.AddWidget(&r.request_page_widget)
	adder.AddWidget(&message_model.MessageModel)
	return nil
}

func (r *Root) HandleButtonInput(ctx *gui.Context, widgetBounds *gui.WidgetBounds) gui.HandleInputResult {
	return gui.HandleInputResult{}
}

func (r *Root) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	b := widgetBounds.Bounds()
	layouter.LayoutWidget(r.message_model_widget, b)
	layouter.LayoutWidget(&r.request_page_widget, b)
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
		WindowMinSize: image.Pt(900, 544),
		RunGameOptions: &ebiten.RunGameOptions{
			ApplePressAndHoldEnabled: true,
		},
	}
	if err := gui.Run(&Root{}, op); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
