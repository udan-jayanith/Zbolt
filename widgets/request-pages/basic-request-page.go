package RequestPage

import (
	"API-Client/basic"
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
	adder.AddChild(&rib.method_select_widget)
	
	rib.input_widget.SetEditable(true)	
	adder.AddChild(&rib.input_widget)
	
	rib.request_btn_widget.SetText("Request")
	adder.AddChild(&rib.request_btn_widget)
	return nil
}

func (rib *RequestInputBar) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	u := widget.UnitSize(ctx)
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Gap: u,
		Items: []gui.LinearLayoutItem{
			{
				Widget: &rib.method_select_widget,
			},
			{
				Widget: &rib.input_widget,
				Size: gui.FlexibleSize(1),
			},
			{
				Widget: &rib.request_btn_widget,
			},
		},
	}
	layout = basic.Align(layout, basic.Start, basic.Center)
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}


func (rib *RequestInputBar) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	u := widget.UnitSize(ctx)
	if w, ok := constraints.FixedWidth(); ok {
		return image.Pt(w, u*2)
	}else if h, ok := constraints.FixedHeight(); ok {
		return image.Pt(u*10, h)
	}
	return image.Pt(u*10, u*2)
}

type RequestWidget struct {
	gui.DefaultWidget
}

type ResponseWidget struct {
	gui.DefaultWidget
}

type BasicRequestPage struct {
	gui.DefaultWidget
}
