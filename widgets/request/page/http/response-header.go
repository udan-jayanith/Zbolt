package http_widget

import (
	"API-Client/basic"

	"image"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type response_header_widget struct {
	gui.DefaultWidget

	status, response_time, size, proto widget.Text
}

func (rhw *response_header_widget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	// TODO: solve this
	rhw.status.SetTabular(true)
	rhw.status.SetValue("200 Ok")
	adder.AddWidget(&rhw.status)

	rhw.response_time.SetTabular(true)
	rhw.response_time.SetValue("200 ms")
	adder.AddWidget(&rhw.response_time)

	rhw.size.SetTabular(true)
	rhw.size.SetValue("131 B")
	adder.AddWidget(&rhw.size)

	rhw.proto.SetTabular(true)
	rhw.proto.SetValue("HTTP v1.1")
	adder.AddWidget(&rhw.proto)
	return nil
}

func (rhw *response_header_widget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	u := widget.UnitSize(ctx)
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Gap:       u / 4,
		Padding:   basic.NewPadding(u / 4),
		Items: []gui.LinearLayoutItem{
			{
				Widget: &rhw.status,
			},
			{
				Widget: &rhw.response_time,
			},
			{
				Widget: &rhw.size,
			},
			{
				Size: gui.FlexibleSize(1),
			},
			{
				Widget: &rhw.proto,
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (rhw *response_header_widget) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	u := widget.UnitSize(ctx)
	gap := u / 4
	padding := u / 2

	point := image.Pt(padding+gap*4, padding)

	if w, ok := constraints.FixedWidth(); ok {
		point.X = w
	} else {
		status_measurements := rhw.status.Measure(ctx, constraints)
		point.X += status_measurements.X

		response_time_measurements := rhw.response_time.Measure(ctx, constraints)
		point.X += response_time_measurements.X

		size := rhw.size.Measure(ctx, constraints)
		point.X += size.X

		proto_measurements := rhw.proto.Measure(ctx, constraints)
		point.X += proto_measurements.X
	}
	
	if h, ok := constraints.FixedHeight(); ok {
		point.Y = h
	}else{
		point.Y += widget.LineHeight(ctx)
	}
	
	return point
}
