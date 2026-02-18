package Requester

import (
	"API-Client/basic"
	"API-Client/icons"
	"image"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type sidebar_item_type_card struct {
	gui.DefaultWidget

	icon_widget icons.Icon
	text_widget widget.Text

	Text, Icon_name string
	Icon_Size       int
}

func (sitc *sidebar_item_type_card) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	sitc.icon_widget.IconName = "large-icons/" + sitc.Icon_name
	icon_size := image.Pt(sitc.Icon_Size, sitc.Icon_Size)
	sitc.icon_widget.Point = &icon_size
	adder.AddChild(&sitc.icon_widget)

	sitc.text_widget.SetValue(sitc.Text)
	sitc.text_widget.SetVerticalAlign(widget.VerticalAlignMiddle)
	sitc.text_widget.SetHorizontalAlign(widget.HorizontalAlignCenter)
	adder.AddChild(&sitc.text_widget)

	return nil
}

func (sitc *sidebar_item_type_card) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	size := widget.UnitSize(ctx) / 4
	layout := gui.LinearLayout{
		Gap:       size,
		Padding:   basic.NewPadding(size),
		Direction: gui.LayoutDirectionVertical,

		Items: []gui.LinearLayoutItem{
			{
				Widget: &sitc.icon_widget,
			},
			{
				Widget: &sitc.text_widget,
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (sitc *sidebar_item_type_card) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	var point image.Point
	size := widget.UnitSize(ctx) / 4

	if h, ok := constraints.FixedHeight(); ok {
		point.Y = h
	} else {
		point.Y = sitc.icon_widget.Measure(ctx, constraints).Y
		point.Y += sitc.text_widget.Measure(ctx, constraints).Y
		point.Y += size * 3
	}
	
	if w, ok := constraints.FixedWidth(); ok {
		point.X = w
	} else {
		point.X = point.Y
	}

	return point
}

type sidebar_item_type_select struct {
	gui.DefaultWidget
}
