package CommonWidgets

import (
	"API-Client/widgets/request/def"
	"image"
	"strings"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type RequestResponse int

const (
	HTTP_Request RequestResponse = iota
	HTTP_Response
)

type http_body_header_widget struct {
	gui.DefaultWidget

	t            RequestResponse
	content_type struct {
		content_type def.ContentType
		text         widget.Text

		input gui.WidgetWithSize[*widget.Combobox]
	}

	options struct {
		auto_wrap struct {
			text   widget.Text
			toggle widget.Toggle
		}
		format struct {
			text   widget.Text
			toggle widget.Toggle
		}
		open_with    widget.Button
		on_open_with func(context *gui.Context, content_type string)
	}
}

func (w *http_body_header_widget) request_build(ctx *gui.Context, adder *gui.ChildAdder) {
	w.content_type.input.SetIntrinsicSize()
	w.content_type.input.SetFixedWidth(widget.UnitSize(ctx) * 5)
	
	input_widget := w.content_type.input.Widget()
	input_widget.SetAllowFreeInput(true)
	input_widget.SetItems([]string{"application/json", "application/octet-stream", "text/html", "text/plain", "image/png", "image/jpeg"})
	w.content_type.content_type = def.ContentType(input_widget.Value())
	adder.AddWidget(&w.content_type.input)
}

func (w *http_body_header_widget) response_build(ctx *gui.Context, adder *gui.ChildAdder) {
	_, sub_t := w.content_type.content_type.Parse()
	w.content_type.text.SetValue(strings.ToUpper(sub_t))
	w.content_type.text.SetVerticalAlign(widget.VerticalAlignMiddle)
	adder.AddWidget(&w.content_type.text)

	w.options.open_with.SetText("Open with")
	if w.options.on_open_with != nil {
		w.options.open_with.OnUp(func(ctx *gui.Context) {
			w.options.on_open_with(ctx, w.content_type.text.Value())
		})
	}

	adder.AddWidget(&w.options.open_with)
}

func (w *http_body_header_widget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	if w.t == HTTP_Request {
		w.request_build(ctx, adder)
	} else {
		w.response_build(ctx, adder)
	}

	{
		w.options.auto_wrap.text.SetValue("Auto wrap")
		w.options.auto_wrap.text.SetVerticalAlign(widget.VerticalAlignMiddle)
		adder.AddWidget(&w.options.auto_wrap.text)

		adder.AddWidget(&w.options.auto_wrap.toggle)
	}
	{
		w.options.format.text.SetValue("Format")
		w.options.format.text.SetVerticalAlign(widget.VerticalAlignMiddle)
		adder.AddWidget(&w.options.format.text)

		adder.AddWidget(&w.options.format.toggle)
	}
	return nil
}

func (rbh *http_body_header_widget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
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
			func() gui.LinearLayoutItem {
				item := gui.LinearLayoutItem{}
				if rbh.t == HTTP_Request {
					item.Widget = &rbh.content_type.input
				} else {
					item.Widget = &rbh.content_type.text
				}

				return item
			}(),
		},
	}

	if rbh.t == HTTP_Response {
		layout.Items = append(layout.Items, gui.LinearLayoutItem{
			Widget: &rbh.options.open_with,
		})
	}

	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (rbh *http_body_header_widget) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	u := widget.UnitSize(ctx)
	gap := u / 4
	var point image.Point

	if w, ok := constraints.FixedWidth(); ok {
		point.X = w
	} else {
		point.X += rbh.options.auto_wrap.text.Measure(ctx, constraints).X + gap
		point.X += rbh.options.auto_wrap.toggle.Measure(ctx, constraints).X + gap
		point.X += rbh.options.format.text.Measure(ctx, constraints).X + gap
		point.X += rbh.options.format.toggle.Measure(ctx, constraints).X + gap

		point.X += func() gui.Widget {
			if rbh.t == HTTP_Request {
				return &rbh.content_type.input
			}

			return &rbh.content_type.text
		}().Measure(ctx, constraints).X + gap

		point.X += rbh.options.open_with.Measure(ctx, constraints).X
	}

	if h, ok := constraints.FixedHeight(); ok {
		point.Y = h
	} else {
		point.Y = u
	}
	return point
}

type BodyWidget struct {
	gui.DefaultWidget
	t RequestResponse

	header http_body_header_widget
	view   widget.TextInput
}

func (w *BodyWidget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	adder.AddWidget(&w.header)

	w.view.SetAutoWrap(w.header.options.auto_wrap.toggle.Value())
	w.view.SetMultiline(true)
	w.view.SetEditable(w.t == HTTP_Request)
	adder.AddWidget(&w.view)

	// make the view handle images and text.
	// Show content type not supported if content type is not jpg, png or text type.
	return nil
}

func (w *BodyWidget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	u := widget.UnitSize(ctx)

	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Gap:       u / 4,
		Items: []gui.LinearLayoutItem{
			{
				Widget: &w.header,
			},
			{
				Widget: &w.view,
				Size:   gui.FlexibleSize(1),
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (body *BodyWidget) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	u := widget.UnitSize(ctx)
	gap := u / 4
	var point image.Point
	point.Y = gap

	if w, ok := constraints.FixedWidth(); ok {
		point.X = w
	} else {
		point.X = body.header.Measure(ctx, constraints).X
	}

	if h, ok := constraints.FixedHeight(); ok {
		point.Y = h
	} else {
		point.Y += body.header.Measure(ctx, constraints).Y
		point.Y += body.view.Measure(ctx, constraints).Y
	}

	return point
}

func (body *BodyWidget) SetType(t RequestResponse) {
	body.t = t
	body.header.t = t
}

func (body *BodyWidget) SetBody(content string, content_type def.ContentType) {
	t, sub_t := content_type.Parse()
	if t == "text" || (t == "application" && sub_t == "json") {
		body.view.SetValue(content)
	}
}

func (body *BodyWidget) ContentType() def.ContentType {
	return body.header.content_type.content_type
}

func (body *BodyWidget) SetContentType(content_type def.ContentType) {
	if body.t == HTTP_Response {
		body.header.content_type.content_type = content_type
	}
}

func (body *BodyWidget) OnAutowrapToggle(fn func(ctx *gui.Context, value bool)) {
	body.header.options.auto_wrap.toggle.OnValueChanged(fn)
}

func (body *BodyWidget) OnFormatToggle(fn func(ctx *gui.Context, value bool)) {
	body.header.options.format.toggle.OnValueChanged(fn)
}

func (body *BodyWidget) SetAutowrap(autowrap bool) {
	body.header.options.auto_wrap.toggle.SetValue(autowrap)
}

func (body *BodyWidget) SetFormat(format bool) {
	body.header.options.format.toggle.SetValue(format)
}
