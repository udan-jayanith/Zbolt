package Requester

import (
	gui "github.com/guigui-gui/guigui"
//	widget "github.com/guigui-gui/guigui/basicwidget"
//	"image"
)

type Requester struct {
	gui.DefaultWidget
	request_widget RequestWidget
	response_widget ResponseWidget
}

func (brp *Requester) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	ctx.SetColorMode(gui.ColorModeDark)

	adder.AddChild(&brp.request_widget)
	adder.AddChild(&brp.response_widget)
	return nil
}

func (brp *Requester) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Items: []gui.LinearLayoutItem{
			{
				Widget: &brp.request_widget,
				Size:   gui.FlexibleSize(1),
			},
			{
				Widget: &brp.response_widget,
				Size:   gui.FlexibleSize(1),
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}
