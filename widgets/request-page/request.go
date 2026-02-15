package Requester

import (
	"API-Client/basic"
	CommonWidgets "API-Client/common-widgets"
	"image"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type RequestWidget struct {
	gui.DefaultWidget
	input_bar_widget RequestInputBar
	url_preview      widget.TextInput

	tab         CommonWidgets.Tab[string]
	tab_content struct {
		params, header  CommonWidgets.AttributeTable
		body            widget.TextInput
		selected_widget gui.Widget
	}
}

func (rw *RequestWidget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	adder.AddChild(&rw.input_bar_widget)

	{
		rw.url_preview.SetEditable(false)
		rw.url_preview.SetValue("https://github.com/guigui-gui/guigui/issues?q=is%3Aissue%20state%3Aopen%20milestone%3Av0.1.0&page=2")
		rw.url_preview.SetMultiline(true)
		rw.url_preview.SetAutoWrap(true)
		adder.AddChild(&rw.url_preview)
	}

	{
		rw.tab_content.body.SetAutoWrap(true)
		rw.tab_content.body.SetMultiline(true)
		rw.tab_content.body.SetEditable(true)

		rw.tab.SetTabItems([]CommonWidgets.TabItem[string]{
			{
				Text: "Parameters",
			},
			{
				Text: "Headers",
			},
			{
				Text: "Body",
			},
		})

		switch rw.tab.GetSelectedIndex() {
		case 0:
			rw.tab_content.selected_widget = &rw.tab_content.params
		case 1:
			rw.tab_content.selected_widget = &rw.tab_content.header
		case 2:
			rw.tab_content.selected_widget = &rw.tab_content.body
		default:
			panic("Unknown tab was selected")
		}

		adder.AddChild(rw.tab_content.selected_widget)
		adder.AddChild(&rw.tab)
	}
	return nil
}

func (rw *RequestWidget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	u := widget.UnitSize(ctx)
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Gap:       u / 4,
		Padding:   basic.NewPadding(u / 4),
		Items: []gui.LinearLayoutItem{
			{
				Widget: &rw.input_bar_widget,
				Size:   gui.FixedSize(u),
			},
			{
				Widget: &rw.url_preview,
				Size:   gui.FixedSize(u * 2),
			},
			{
				Widget: &rw.tab,
			},
			{
				Widget: rw.tab_content.selected_widget,
				Size:   gui.FlexibleSize(1),
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (rw *RequestWidget) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	point := rw.input_bar_widget.Measure(ctx, constraints)

	if h, ok := constraints.FixedHeight(); ok {
		point.Y = h
	} else {
		point.Y += rw.url_preview.Measure(ctx, constraints).Y
		point.Y += rw.tab.Measure(ctx, constraints).Y
		point.Y += rw.tab_content.selected_widget.Measure(ctx, constraints).Y
	}
	
	if w, ok := constraints.FixedWidth(); ok {
		point.X = w
	}

	return point
}
