package http_widget

import (
	"API-Client/basic"

	"image"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/hajimehoshi/ebiten/v2"
)

type response_header_widget struct {
	gui.DefaultWidget

	status, response_time, size, proto widget.Text

	show_tooltip   bool
	tooltip        widget.TooltipArea
	tooltip_bounds image.Rectangle
}

func (rhw *response_header_widget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	rhw.status.SetTabular(true)
	if rhw.status.Value() == "" {
		rhw.status.SetValue("None")
	}
	adder.AddWidget(&rhw.status)

	rhw.response_time.SetTabular(true)
	if rhw.response_time.Value() == "" {
		rhw.response_time.SetValue("0ms")
	}
	adder.AddWidget(&rhw.response_time)

	rhw.size.SetTabular(true)
	if rhw.size.Value() == "" {
		rhw.size.SetValue("0B")
	}
	adder.AddWidget(&rhw.size)

	rhw.proto.SetTabular(true)
	if rhw.proto.Value() == "" {
		rhw.proto.SetValue("HTTP vX.X")
	}
	adder.AddWidget(&rhw.proto)

	if rhw.show_tooltip {
		adder.AddWidget(&rhw.tooltip)
	}
	return nil
}

func (rhw *response_header_widget) padding(ctx *gui.Context) gui.Padding {
	return basic.NewPadding(widget.UnitSize(ctx) / 4)
}

func (rhw *response_header_widget) gap(ctx *gui.Context) int {
	return widget.UnitSize(ctx) / 4
}

func (rhw *response_header_widget) show_tooltip_widget(show bool, text string, tooltip_bounds image.Rectangle) {
	rhw.show_tooltip = show
	rhw.tooltip.SetText(text)
	rhw.tooltip_bounds = tooltip_bounds
}

func (rhw *response_header_widget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	if rhw.show_tooltip {
		layouter.LayoutWidget(&rhw.tooltip, rhw.tooltip_bounds)
	}

	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Gap:       rhw.gap(ctx),
		Padding:   rhw.padding(ctx),
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

func (wi *response_header_widget) HandlePointingInput(ctx *gui.Context, widgetBounds *gui.WidgetBounds) gui.HandleInputResult {
	if !widgetBounds.IsHitAtCursor() {
		wi.show_tooltip = false
		return gui.HandleInputResult{}
	}

	padding := wi.padding(ctx)
	gap := wi.gap(ctx)
	b := widgetBounds.Bounds()

	cursor_x, _ := ebiten.CursorPosition()
	left_x := padding.Start + b.Min.X

	w := wi.status.Measure(ctx, gui.Constraints{}).X
	if left_x <= cursor_x && cursor_x <= w+left_x {
		wi.show_tooltip_widget(true, "HTTP status", image.Rect(left_x, b.Min.Y, left_x+w, b.Max.Y))
		return gui.HandleInputResult{}
	}

	left_x += w + gap
	w = wi.response_time.Measure(ctx, gui.Constraints{}).X
	if left_x <= cursor_x && cursor_x <= w+left_x {
		wi.show_tooltip_widget(true, "Response time", image.Rect(left_x, b.Min.Y, left_x+w, b.Max.Y))
		return gui.HandleInputResult{}
	}

	left_x += w + gap
	w = wi.size.Measure(ctx, gui.Constraints{}).X
	if left_x <= cursor_x && cursor_x <= w+left_x {
		wi.show_tooltip_widget(true, "Response size", image.Rect(left_x, b.Min.Y, left_x+w, b.Max.Y))
		return gui.HandleInputResult{}
	}

	right_x := b.Max.X - padding.End
	w = wi.proto.Measure(ctx, gui.Constraints{}).X
	if cursor_x <= right_x && cursor_x >= right_x-w {
		wi.show_tooltip_widget(true, "Used HTTP protocol version", image.Rect(right_x-w, b.Min.Y, right_x, b.Max.Y))
		return gui.HandleInputResult{}
	}
	wi.show_tooltip = false
	return gui.HandleInputResult{}
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
	} else {
		point.Y += widget.LineHeight(ctx)
	}

	return point
}
