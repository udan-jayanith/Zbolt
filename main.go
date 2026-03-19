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
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/hajimehoshi/ebiten/v2"

	messages "API-Client/massages"
	home "API-Client/widgets/home"
	request_page "API-Client/widgets/request/page"
)

type Root struct {
	gui.DefaultWidget

	welcome_page_widget home.HomePage
	request_page_widget request_page.RequestPage

	alert_box messages.AlertBox
	popup     widget.Popup
}

func (r *Root) Build(context *gui.Context, adder *gui.ChildAdder) error {
	adder.AddWidget(&r.request_page_widget)

	if r.popup.IsOpen() {
		r.popup.SetAnimated(true)
		r.popup.SetContent(&r.alert_box)
		adder.AddWidget(&r.popup)
	} else {
		alert, ok := messages.Alerts.Get()
		r.popup.SetOpen(ok)
		r.alert_box.SetAlert(alert)

		r.alert_box.OnOk(func(ctx *gui.Context) {
			alert, ok := messages.Alerts.Get()
			r.popup.SetOpen(ok)
			r.alert_box.SetAlert(alert)
		})
	}

	return nil
}

func (r *Root) HandleButtonInput(ctx *gui.Context, widgetBounds *gui.WidgetBounds) gui.HandleInputResult {
	return gui.HandleInputResult{}
}

func (r *Root) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	b := widgetBounds.Bounds()

	if r.popup.IsOpen() {
		alert_box_point := r.alert_box.Measure(ctx, gui.Constraints{})
		middle := image.Rectangle{
			Min: image.Point{
				X: b.Max.X/2 - alert_box_point.X/2,
				Y: b.Max.Y/2 - alert_box_point.Y/2,
			},
		}
		middle.Max = middle.Min.Add(alert_box_point)

		layouter.LayoutWidget(&r.popup, middle)
	}

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
