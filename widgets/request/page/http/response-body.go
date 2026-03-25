package http_widget

import (
	"image"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type response_body_header struct {
	gui.DefaultWidget

	file_type       widget.Text

	options struct {
		auto_wrap struct {
			text   widget.Text
			toggle widget.Toggle
		}
		format struct {
			text   widget.Text
			toggle widget.Toggle
		}
		open_with widget.Button
	}
}

func (rbh *response_body_header) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	{
		rbh.file_type.SetValue("Json")
		rbh.file_type.SetVerticalAlign(widget.VerticalAlignMiddle)
		adder.AddWidget(&rbh.file_type)
	}
	{
		rbh.options.auto_wrap.text.SetValue("Auto wrap")
		rbh.options.auto_wrap.text.SetVerticalAlign(widget.VerticalAlignMiddle)
		adder.AddWidget(&rbh.options.auto_wrap.text)

		adder.AddWidget(&rbh.options.auto_wrap.toggle)
	}
	{
		rbh.options.format.text.SetValue("Format")
		rbh.options.format.text.SetVerticalAlign(widget.VerticalAlignMiddle)
		adder.AddWidget(&rbh.options.format.text)

		adder.AddWidget(&rbh.options.format.toggle)
	}
	{
		rbh.options.open_with.SetText("Open")
		rbh.options.open_with.OnUp(func(context *gui.Context) {})
		adder.AddWidget(&rbh.options.open_with)
	}
	return nil
}

func (rbh *response_body_header) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	u := widget.UnitSize(ctx)
	toggle_size := gui.FixedSize(u*2 - u/3)
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Gap:       u / 4,
		Items: []gui.LinearLayoutItem{
			{
				Widget: &rbh.options.auto_wrap.text,
			},
			{
				Widget: &rbh.options.auto_wrap.toggle,
				Size:   toggle_size,
			},
			{
				Widget: &rbh.options.format.text,
			},
			{
				Widget: &rbh.options.format.toggle,
				Size:   toggle_size,
			},
			{
				Size: gui.FlexibleSize(1),
			},
			{
				Widget: &rbh.file_type,
			},
			{
				Widget: &rbh.options.open_with,
			},
		},
	}

	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (rbh *response_body_header) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	u := widget.UnitSize(ctx)
	gap := u/4
	var point image.Point
	
	if w, ok := constraints.FixedWidth(); ok {
		point.X = w
	}else{
		point.X += rbh.options.auto_wrap.text.Measure(ctx, constraints).X+gap
		point.X += rbh.options.auto_wrap.toggle.Measure(ctx, constraints).X+gap
		point.X += rbh.options.format.text.Measure(ctx, constraints).X+gap
		point.X += rbh.options.format.toggle.Measure(ctx, constraints).X+gap
		point.X += rbh.file_type.Measure(ctx, constraints).X+gap
		point.X += rbh.options.open_with.Measure(ctx, constraints).X
	}
	
	if h, ok := constraints.FixedHeight(); ok {
		point.Y = h
	}else{
		point.Y = u
	}
	return point
}

type response_body_widget struct {
	gui.DefaultWidget
	header response_body_header
	view widget.TextInput
}

func (rbw *response_body_widget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	adder.AddWidget(&rbw.header)
	
	rbw.view.SetAutoWrap(rbw.header.options.auto_wrap.toggle.Value())
	rbw.view.SetMultiline(true)
	rbw.view.SetEditable(false)
	adder.AddWidget(&rbw.view)

	return nil
}

func (rbw *response_body_widget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	u := widget.UnitSize(ctx)

	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Gap:       u / 4,
		Items: []gui.LinearLayoutItem{
			{
				Widget: &rbw.header,
			},
			{
				Widget: &rbw.view,
				Size:   gui.FlexibleSize(1),
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}


func (rbw *response_body_widget) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	u := widget.UnitSize(ctx)
	gap := u/4
	var point image.Point
	point.Y = gap
	
	if w, ok := constraints.FixedWidth(); ok {
		point.X = w
	}else{
		point.X = rbw.header.Measure(ctx, constraints).X
	}
	
	if h, ok := constraints.FixedHeight(); ok {
		point.Y = h
	}else{
		point.Y += rbw.header.Measure(ctx, constraints).Y
		point.Y += rbw.view.Measure(ctx, constraints).Y
	}
	
	return point
}