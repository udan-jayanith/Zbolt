package CommonWidgets

import (
	"API-Client/basic"
	"image"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type popup_form_content struct {
	gui.DefaultWidget

	field_widget  widget.Text
	input_widget  widget.TextInput
	button_widget widget.Button
}

func (content *popup_form_content) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	if content.field_widget.Value() != "" {
		content.field_widget.SetVerticalAlign(widget.VerticalAlignBottom)
		content.field_widget.SetHorizontalAlign(widget.HorizontalAlignLeft)
		content.field_widget.SetOpacity(0.82)
		adder.AddChild(&content.field_widget)
	}

	adder.AddChild(&content.input_widget)

	content.button_widget.SetType(widget.ButtonTypePrimary)
	adder.AddChild(&content.button_widget)
	return nil
}

func (content *popup_form_content) gap(ctx *gui.Context) int {
	return widget.UnitSize(ctx) / 4
}

func (content *popup_form_content) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	gap := content.gap(ctx)

	horizontal_layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Gap:       gap,
		Items: []gui.LinearLayoutItem{
			{
				Widget: &content.input_widget,
				Size:   gui.FlexibleSize(1),
			},
			{
				Widget: &content.button_widget,
			},
		},
	}

	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Gap: gap/2,
		Items: []gui.LinearLayoutItem{
			{
				Widget: &content.field_widget,
			},
			{
				Layout: horizontal_layout,
			},
		},
	}

	if content.field_widget.Value() == "" {
		layout = horizontal_layout
	}

	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (content *popup_form_content) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	var point image.Point
	u := widget.UnitSize(ctx)
	gap := content.gap(ctx)

	point.X = u * 10
	if content.field_widget.Value() != "" {
		point.Y += content.field_widget.Measure(ctx, gui.Constraints{}).Y + gap/2
	}

	point.Y += max(content.button_widget.Measure(ctx, gui.Constraints{}).Y, content.input_widget.Measure(ctx, gui.Constraints{}).Y)
	return point
}

type SimpleFormPopup struct {
	gui.DefaultWidget

	popup_widget   widget.Popup
	popup_content  popup_form_content
	padding_widget WidgetWithPadding[*popup_form_content]
}

func (sfp *SimpleFormPopup) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	sfp.padding_widget.SetWidget(&sfp.popup_content)
	sfp.padding_widget.SetPadding(basic.NewPadding(widget.UnitSize(ctx) / 3))

	sfp.popup_widget.SetContent(&sfp.padding_widget)
	sfp.popup_widget.SetAnimated(true)
	sfp.popup_widget.SetCloseByClickingOutside(true)
	adder.AddChild(&sfp.popup_widget)
	return nil
}

func (sfp *SimpleFormPopup) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layouter.LayoutWidget(&sfp.popup_widget, widgetBounds.Bounds())
}

func (sfp *SimpleFormPopup) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	return sfp.padding_widget.Measure(ctx, constraints)
}

func (sfp *SimpleFormPopup) SetButtonText(text string) {
	sfp.popup_content.button_widget.SetText(text)
}

func (sfp *SimpleFormPopup) SetFieldValue(text string) {
	sfp.popup_content.field_widget.SetValue(text)
}

func (sfp *SimpleFormPopup) SetOpen(open bool) {
	sfp.popup_widget.SetOpen(open)
}
