package websocket_widget

import (
	CommonWidgets "API-Client/common-widgets"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type request_widget struct {
	gui.DefaultWidget

	url_input_bar struct {
		url_input      widget.TextInput
		connect_button widget.Button
	}
	
	url_preview CommonWidgets.URLPreview

	tab_widget CommonWidgets.Tab[struct{}]
	content    struct {
		params, headers CommonWidgets.AttributeTable
		body            widget.TextInput
		selected        gui.Widget
	}
}

func (rw *request_widget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	adder.AddWidget(&rw.url_input_bar.url_input)

	rw.url_input_bar.connect_button.SetText("Connect")
	rw.url_input_bar.connect_button.SetType(widget.ButtonTypePrimary)
	adder.AddWidget(&rw.url_input_bar.connect_button)

	adder.AddWidget(&rw.url_preview)
	
	rw.tab_widget.SetTabItems([]CommonWidgets.TabItem[struct{}]{
		{
			Text: "Params",
		},
		{
			Text: "Body",
		},
		{
			Text: "Header",
		},
	})

	if rw.content.selected == nil {
		rw.content.selected = &rw.content.params
	}

	rw.tab_widget.OnSelect(func(tab_item *CommonWidgets.TabItem[struct{}], index int) {
		switch tab_item.Text {
		case "Params":
			rw.content.selected = &rw.content.params
		case "Body":
			rw.content.selected = &rw.content.body
		case "Header":
			rw.content.selected = &rw.content.headers
		}
	})

	adder.AddWidget(&rw.tab_widget)
	adder.AddWidget(rw.content.selected)
	return nil
}

func (rw *request_widget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
}
