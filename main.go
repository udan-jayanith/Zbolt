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
	"github.com/guigui-gui/guigui/basicwidget"
	"github.com/hajimehoshi/ebiten/v2"

	"API-Client/basic"
	message_model "API-Client/message-model"
	"API-Client/widgets/request/def"
	http_widget "API-Client/widgets/request/page/http"
)

type Root struct {
	gui.DefaultWidget
	background  basicwidget.Background
	http_widget http_widget.HTTP_Widget
	req         *def.Request
}

func (r *Root) Build(context *gui.Context, adder *gui.ChildAdder) error {
	adder.AddWidget(&r.background)
	if r.req == nil {
		req := def.NewRequest(def.HTTP, "")
		r.req = &req
	}
	r.http_widget.SetReq(r.req)
	adder.AddWidget(&r.http_widget)
	adder.AddWidget(&message_model.MessageModel)
	return nil
}

func (r *Root) HandleButtonInput(ctx *gui.Context, widgetBounds *gui.WidgetBounds) gui.HandleInputResult {
	return gui.HandleInputResult{}
}

func (r *Root) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	b := widgetBounds.Bounds()
	layouter.LayoutWidget(&r.background, b)
	layouter.LayoutWidget(&message_model.MessageModel, widgetBounds.Bounds())

	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Padding: basic.NewPadding(basicwidget.UnitSize(ctx)/4),
		Items: []gui.LinearLayoutItem{
			{
				Widget: &r.http_widget,
				Size:   gui.FlexibleSize(1),
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
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

	message_model.Show("Hello world", message_model.Alert, nil)
	if err := gui.Run(&Root{}, op); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
