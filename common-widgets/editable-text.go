package CommonWidgets

import (
	"image"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type EditableText struct {
	gui.DefaultWidget
	not_editable bool
	widget widget.Text
}

func (et *EditableText) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	et.widget.SetTabular(true)
	et.widget.SetVerticalAlign(widget.VerticalAlignMiddle)
	et.widget.SetHorizontalAlign(widget.HorizontalAlignLeft)
	et.widget.SetSelectable(true)
	et.widget.SetEditable(!et.not_editable)
	adder.AddWidget(&et.widget)
	return nil
}

func (et *EditableText) SetValue(text string) {
	et.widget.SetValue(text)
}

func (et *EditableText) Value() string {
	return et.widget.Value()
}

func (et *EditableText) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layouter.LayoutWidget(&et.widget, widgetBounds.Bounds())
}

func (et *EditableText) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	return et.widget.Measure(ctx, constraints)
}

func (et *EditableText) SetEditable(editable bool) {
	et.not_editable = !editable
}