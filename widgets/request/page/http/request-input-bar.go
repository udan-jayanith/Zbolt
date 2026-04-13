package http_widget

import (
	"API-Client/icons"
	"image"
	"log"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/hajimehoshi/ebiten/v2"
)

type request_input_bar_widget struct {
	gui.DefaultWidget

	method_select_widget  widget.Select[string]
	method_list           []string
	on_method_changed_fn  func(method string)
	selected_method_index int

	url_input          widget.TextInput
	url_input_disabled bool

	request_button_text          string
	request_btn_widget           widget.Button
	on_request_button_clicked_fn func(ctx *gui.Context, value string)

	open_in_icon *ebiten.Image
	open_in      widget.Button
	on_request   func(ctx *gui.Context, url, method string)
}

func (rib *request_input_bar_widget) init_methods() {
	if len(rib.method_list) == 0 {
		rib.method_list = []string{
			"Get",
			"Post",
			"Put",
			"Patch",
			"Delete",
			"Options",
			"Head",
		}
	}
}

func (rib *request_input_bar_widget) select_method(method string) {
	rib.init_methods()
	for i, v := range rib.method_list {
		if v != method {
			continue
		}

		rib.selected_method_index = i
		break
	}
}

func (rib *request_input_bar_widget) on_method_changed(fn func(method string)) {
	rib.on_method_changed_fn = fn
}

func (rib *request_input_bar_widget) method() string {
	rib.init_methods()
	return rib.method_list[rib.selected_method_index]
}

func (rib *request_input_bar_widget) url_input_value() string {
	return rib.url_input.Value()
}

func (rib *request_input_bar_widget) set_url_input_value(value string) {
	rib.url_input.ForceSetValue(value)
}

func (rib *request_input_bar_widget) disable_url_input(disabled bool) {
	rib.url_input_disabled = disabled
}

func (rib *request_input_bar_widget) on_url_input_value_changed(fn func(context *gui.Context, text string, committed bool)) {
	rib.url_input.OnValueChanged(fn)
}

func (rib *request_input_bar_widget) init_request_button_text() {
	if rib.request_button_text == "" {
		// TODO: Set a icon
		rib.request_button_text = "Request"
	}
}

func (rib *request_input_bar_widget) on_request_button_clicked(fn func(ctx *gui.Context, value string)) {
	rib.on_request_button_clicked_fn = fn
}

func (rib *request_input_bar_widget) set_request_button_value(value string) {
	if value == "Request" || value == "Cancel" {
		rib.request_button_text = value
	}
	log.Fatalln("Unknown request button text:", value)
}

func (rib *request_input_bar_widget) on_open_in_clicked(fn func(context *gui.Context)) {
	rib.open_in.OnUp(fn)
}

func (rib *request_input_bar_widget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	rib.init_methods()
	rib.method_select_widget.SetItemsByStrings(rib.method_list)
	rib.method_select_widget.OnItemSelected(func(_ *gui.Context, index int) {
		rib.selected_method_index = index
		if rib.on_method_changed_fn != nil {
			rib.on_method_changed_fn(rib.method_list[rib.selected_method_index])
		}
	})
	rib.method_select_widget.SelectItemByIndex(rib.selected_method_index)
	adder.AddWidget(&rib.method_select_widget)

	ctx.SetEnabled(&rib.url_input, !rib.url_input_disabled)
	adder.AddWidget(&rib.url_input)

	if rib.open_in_icon == nil {
		rib.open_in_icon = icons.Store.Open("open-in")
	}
	rib.open_in.SetIcon(rib.open_in_icon)
	adder.AddWidget(&rib.open_in)

	rib.init_request_button_text()
	if rib.request_button_text == "Request" {
		rib.request_btn_widget.SetType(widget.ButtonTypePrimary)
	} else {
		rib.request_btn_widget.SetType(widget.ButtonTypeNormal)
	}
	rib.request_btn_widget.SetText(rib.request_button_text)
	rib.request_btn_widget.OnDown(func(ctx *gui.Context) {
		if rib.on_request_button_clicked_fn == nil {
			return
		}
		rib.on_request_button_clicked_fn(ctx, rib.request_button_text)
	})
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
							Widget: &rib.url_input,
							Size:   gui.FlexibleSize(1),
						},
						{
							Widget: &rib.open_in,
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
