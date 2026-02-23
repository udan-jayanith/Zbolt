package http

import (
	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"API-Client/widgets/request/def"
	// "image"
)

type HTTP_Widget struct {
	gui.DefaultWidget
	request_widget  request_widget
	response_widget response_widget
}

func (brp *HTTP_Widget) RequestType() def.RequestType {
	return def.HTTP
}

func (brp *HTTP_Widget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	ctx.SetColorMode(gui.ColorModeDark)

	adder.AddChild(&brp.request_widget)
	adder.AddChild(&brp.response_widget)
	return nil
}

func (brp *HTTP_Widget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Gap:       widget.UnitSize(ctx) / 4,
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