package CommonWidgets

import (
	"image"

	gui "github.com/guigui-gui/guigui"
)

type WidgetWithPadding[T gui.Widget] struct {
	gui.DefaultWidget

	is_set bool
	widget  T
	padding gui.Padding
}

func (wwp *WidgetWithPadding[T]) SetPadding(padding gui.Padding) {
	wwp.padding = padding
}

func (wwp *WidgetWithPadding[T]) SetWidget(widget T){
	wwp.widget = widget
	wwp.is_set = true
}

func (wwp *WidgetWithPadding[T]) IsSet() bool {
	return wwp.is_set
}

func (wwp *WidgetWithPadding[T]) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	adder.AddWidget(wwp.widget)
	return nil
}

func (wwp *WidgetWithPadding[T]) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layout := gui.LinearLayout{
		Padding: wwp.padding,
		Items: []gui.LinearLayoutItem{
			{
				Widget: wwp.widget,
				Size:   gui.FlexibleSize(1),
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (wwp *WidgetWithPadding[T]) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	point := wwp.widget.Measure(ctx, constraints)
	point.X += wwp.padding.Start + wwp.padding.End
	point.Y += wwp.padding.Top + wwp.padding.Bottom
	return point
}
