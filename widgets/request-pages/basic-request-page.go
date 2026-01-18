package RequestPage

import (
	"image"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type RequestInputBar struct {
	gui.DefaultWidget
	method_select_widget widget.Select[string]
	input_widget         widget.TextInput
	request_btn_widget   widget.Button
}

func (rib *RequestInputBar) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	rib.method_select_widget.SetItemsByStrings([]string{
		"Get",
		"Post",
		"Put",
		"Patch",
		"Delete",
		"Options",
		"Head",
	})

	selected_item_index := max(rib.method_select_widget.SelectedItemIndex(), 0)
	rib.method_select_widget.SelectItemByIndex(selected_item_index)
	adder.AddChild(&rib.method_select_widget)

	rib.input_widget.SetEditable(true)
	adder.AddChild(&rib.input_widget)

	rib.request_btn_widget.SetText("Request")
	rib.request_btn_widget.SetTextBold(true)
	adder.AddChild(&rib.request_btn_widget)
	return nil
}

func (rib *RequestInputBar) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	u := widget.UnitSize(ctx)
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Items: []gui.LinearLayoutItem{
			{
				Size: gui.FixedSize(u),
				Layout: gui.LinearLayout{
					Direction: gui.LayoutDirectionHorizontal,
					Gap:       u / 4,
					Items: []gui.LinearLayoutItem{
						{
							Widget: &rib.method_select_widget,
						},
						{
							Widget: &rib.input_widget,
							Size:   gui.FlexibleSize(1),
						},
						{
							Widget: &rib.request_btn_widget,
						},
					},
				},
			},
		},
	}

	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (rib *RequestInputBar) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	u := widget.UnitSize(ctx)
	if w, ok := constraints.FixedWidth(); ok {
		return image.Pt(w, u*2)
	} else if h, ok := constraints.FixedHeight(); ok {
		return image.Pt(u*10, h)
	}
	return image.Pt(u*10, u*2)
}

type RequestWidget struct {
	gui.DefaultWidget
	input_bar_widget RequestInputBar
	url_preview widget.TextInput
}

func (rw *RequestWidget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	adder.AddChild(&rw.input_bar_widget)
	
	rw.url_preview.SetEditable(false)
	rw.url_preview.SetValue("https://github.com/guigui-gui/guigui/issues?q=is%3Aissue%20state%3Aopen%20milestone%3Av0.1.0&page=2")
	rw.url_preview.SetMultiline(true)
	rw.url_preview.SetAutoWrap(true)
	adder.AddChild(&rw.url_preview)
	return nil
}

func (rw *RequestWidget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	u := widget.UnitSize(ctx)
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Gap: u/4,
		Items: []gui.LinearLayoutItem{
			{
				Widget: &rw.input_bar_widget,
				Size: gui.FixedSize(u),
			},
			{
				Widget: &rw.url_preview,
				Size: gui.FixedSize(u*2),
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

type ResponseWidget struct {
	gui.DefaultWidget
}

type BasicRequestPage struct {
	gui.DefaultWidget
	background     widget.Background
	request_widget RequestWidget
}

func (brp *BasicRequestPage) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	ctx.SetColorMode(gui.ColorModeDark)
	adder.AddChild(&brp.background)

	adder.AddChild(&brp.request_widget)
	return nil
}

func (brp *BasicRequestPage) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layouter.LayoutWidget(&brp.background, widgetBounds.Bounds())

	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Items: []gui.LinearLayoutItem{
			{
				Widget: &brp.request_widget,
				Size:   gui.FlexibleSize(1),
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}
