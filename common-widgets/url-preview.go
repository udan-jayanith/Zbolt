package CommonWidgets

import (
	"github.com/udan-jayanith/Zbolt/icons"
	"image"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type URLPreview struct {
	gui.DefaultWidget

	url_preview widget.TextInput
	copy_button widget.Button
}

func (up *URLPreview) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	up.url_preview.SetAutoWrap(true)
	up.url_preview.SetEditable(false)
	up.url_preview.SetHorizontalAlign(widget.HorizontalAlignLeft)
	up.url_preview.SetVerticalAlign(widget.VerticalAlignTop)
	adder.AddWidget(&up.url_preview)

	up.copy_button.SetIcon(icons.Store.Open("copy-all"))
	adder.AddWidget(&up.copy_button)
	return nil
}

func (up *URLPreview) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Gap:       widget.UnitSize(ctx) / 4,
		Items: []gui.LinearLayoutItem{
			{
				Widget: &up.url_preview,
				Size:   gui.FlexibleSize(1),
			},
			{
				Widget: &up.copy_button,
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (up *URLPreview) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	var point image.Point
	u := widget.UnitSize(ctx)

	if w, ok := constraints.FixedWidth(); ok {
		point.X = w
	} else {
		point.X += u / 4
		point.X += up.copy_button.Measure(ctx, gui.Constraints{}).X
		point.X += up.url_preview.Measure(ctx, gui.Constraints{}).X
	}

	if h, ok := constraints.FixedHeight(); ok {
		point.Y = h
	} else {
		point.Y = u * 2
	}
	return point
}

func (up *URLPreview) SetURL(url string) {
	up.url_preview.SetValue(url)
}

func (up *URLPreview) URL() string {
	return up.url_preview.Value()
}
