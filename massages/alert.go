package messages

import (
	"github.com/udan-jayanith/Zbolt/basic"
	"image"

	queue "github.com/golang-ds/queue/slicequeue"
	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type alerts struct {
	queue queue.SliceQueue[string]
}

func (a *alerts) Push(message string) {
	a.queue.Enqueue(message)
}

func (a *alerts) Get() (string, bool) {
	return a.queue.Dequeue()
}

func (a *alerts) Len() int {
	return a.queue.Size()
}

var Alerts alerts = alerts{} 

type AlertBox struct {
	gui.DefaultWidget

	text      widget.Text
	panel     widget.Panel
	ok_button widget.Button
	on_ok func(ctx *gui.Context)
}

func (w *AlertBox) SetAlert(alert string) {
	w.text.SetValue(alert)
}

func (w *AlertBox) OnOk(fn func(ctx *gui.Context)){
	w.on_ok = fn
}

func (w *AlertBox) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	w.text.SetAutoWrap(true)
	w.text.SetMultiline(true)
	w.text.SetHorizontalAlign(widget.HorizontalAlignLeft)
	w.text.SetVerticalAlign(widget.VerticalAlignTop)
	w.text.SetSelectable(true)

	w.panel.SetContentConstraints(widget.PanelContentConstraintsFixedWidth)
	w.panel.SetContent(&w.text)
	adder.AddWidget(&w.panel)

	w.ok_button.SetText("Ok")
	w.ok_button.SetType(widget.ButtonTypePrimary)
	if w.on_ok != nil {
		w.ok_button.OnDown(w.on_ok)
	}
	adder.AddWidget(&w.ok_button)
	return nil
}

func (w *AlertBox) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	u := widget.UnitSize(ctx)
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Gap:       u / 4,
		Padding:   basic.NewPadding(u / 2),
		Items: []gui.LinearLayoutItem{
			{
				Widget: &w.panel,
				Size:   gui.FlexibleSize(1),
			},
			{
				Layout: gui.LinearLayout{
					Direction: gui.LayoutDirectionHorizontal,
					Items: []gui.LinearLayoutItem{
						{
							Size: gui.FlexibleSize(1),
						},
						{
							Widget: &w.ok_button,
						},
					},
				},
			},
		},
	}
	
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (w *AlertBox) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	var point image.Point
	u := widget.UnitSize(ctx)
	point.X = u*14
	point.Y = u*8
	return point
}
