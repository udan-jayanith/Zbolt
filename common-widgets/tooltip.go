package CommonWidgets

import (
	"image"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type Opener interface {
}

type WidgetWithTooltip[T gui.Widget] struct {
	gui.DefaultWidget

	widget       lazy_widget[T]
	tooltip      widget.TooltipArea
	tooltip_text string
}

func (w *WidgetWithTooltip[T]) SetTooltip(text string) {
	w.tooltip_text = text
}

func (w *WidgetWithTooltip[T]) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	adder.AddWidget(w.widget.Widget())

	if w.tooltip_text != "" {
		w.tooltip.SetText(w.tooltip_text)
		adder.AddWidget(&w.tooltip)
	}
	return nil
}

func (w *WidgetWithTooltip[T]) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	if w.tooltip_text != "" {
		layouter.LayoutWidget(&w.tooltip, widgetBounds.Bounds())
	}
	layouter.LayoutWidget(w.widget.Widget(), widgetBounds.Bounds())
}

func (w *WidgetWithTooltip[T]) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	return w.widget.Widget().Measure(ctx, constraints)
}

// TODO: Add tooltip with button
// TODO: Add tooltip with text