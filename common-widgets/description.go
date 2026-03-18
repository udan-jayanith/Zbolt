package CommonWidgets

import (
	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type Description struct {
	gui.DefaultWidget

	text_widget widget.Text
}

func (w *Description) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	w.text_widget.SetEllipsisString("...")
	w.text_widget.SetMultiline(true)
	w.text_widget.SetAutoWrap(true)
	w.text_widget.SetHorizontalAlign(widget.HorizontalAlignLeft)
	w.text_widget.SetVerticalAlign(widget.VerticalAlignMiddle)
	adder.AddWidget(&w.text_widget)
	return nil
}

func (w *Description) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layouter.LayoutWidget(&w.text_widget, widgetBounds.Bounds())
}

func (w *Description) SetDescription(description string){
	w.text_widget.SetValue(description)
}
