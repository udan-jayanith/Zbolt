package http_widget

import (
	"image"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type request_input_bar_widget struct {
	gui.DefaultWidget
	method_select_widget widget.Select[string]
	input_widget         widget.TextInput
	request_btn_widget   widget.Button
	on_request func(ctx *gui.Context, url, method string)
}

func (rib *request_input_bar_widget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
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
	adder.AddWidget(&rib.method_select_widget)

	rib.input_widget.SetEditable(true)
	adder.AddWidget(&rib.input_widget)

	rib.request_btn_widget.SetText("Request")
	rib.request_btn_widget.SetType(widget.ButtonTypePrimary)
	adder.AddWidget(&rib.request_btn_widget)
	return nil
}

func (rib *request_input_bar_widget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
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

func (rib *request_input_bar_widget) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	u := widget.UnitSize(ctx)
	if w, ok := constraints.FixedWidth(); ok {
		return image.Pt(w, u*2)
	} else if h, ok := constraints.FixedHeight(); ok {
		return image.Pt(u*10, h)
	}
	return image.Pt(u*10, u*2)
}

func (rib *request_input_bar_widget) OnRequest(fn func(ctx *gui.Context, url, method string)){
	rib.request_btn_widget.OnDown(func(ctx *gui.Context) {
		method, _ := rib.method_select_widget.SelectedItem()
		fn(ctx, rib.input_widget.Value(), method.Text)
	})
}