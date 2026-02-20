package CommonWidgets

import (
	"image"

	gui "github.com/guigui-gui/guigui"
	"github.com/guigui-gui/guigui/basicwidget/basicwidgetdraw"
	"github.com/hajimehoshi/ebiten/v2"
)

type WidgetWithBorder[T gui.Widget] struct {
	gui.DefaultWidget

	widget T
}

func (item *WidgetWithBorder[T]) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	adder.AddChild(item.widget)
	return nil
}

func (item *WidgetWithBorder[T]) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layouter.LayoutWidget(item.widget, widgetBounds.Bounds())
}

func (item *WidgetWithBorder[T]) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	return item.widget.Measure(ctx, constraints)
}

func (wwb *WidgetWithBorder[T]) Draw(ctx *gui.Context, widgetBounds *gui.WidgetBounds, dst *ebiten.Image) {
	color_mod := ctx.ColorMode()
	background_color := basicwidgetdraw.ControlColor(color_mod, ctx.IsEnabled(wwb))
	border_color := basicwidgetdraw.ControlSecondaryColor(color_mod, ctx.IsEnabled(wwb))
	basicwidgetdraw.DrawRoundedRectBorder(ctx, dst, widgetBounds.Bounds(), background_color, border_color, 1, 1, basicwidgetdraw.RoundedRectBorderTypeRegular)
}

func (wwb *WidgetWithBorder[T]) SetWidget(widget T) {
	wwb.widget = widget
}

func (wwb *WidgetWithBorder[T]) GetWidget() T {
	return wwb.widget
}
