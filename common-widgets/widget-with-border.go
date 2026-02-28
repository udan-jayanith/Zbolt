package CommonWidgets

import (
	"image"

	gui "github.com/guigui-gui/guigui"
	"github.com/guigui-gui/guigui/basicwidget/basicwidgetdraw"
	"github.com/hajimehoshi/ebiten/v2"
)

type WidgetWithBorder[T gui.Widget] struct {
	gui.DefaultWidget

	widget      T
	radius      int
	borderWidth float32
	borderType  basicwidgetdraw.RoundedRectBorderType
}

func (item *WidgetWithBorder[T]) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	adder.AddWidget(item.widget)
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
	basicwidgetdraw.DrawRoundedRectBorder(ctx, dst, widgetBounds.Bounds(), background_color, border_color, max(wwb.radius, 1), max(wwb.borderWidth, 1), wwb.borderType)
}

func (wwb *WidgetWithBorder[T]) SetWidget(widget T) {
	wwb.widget = widget
}

func (wwb *WidgetWithBorder[T]) GetWidget() T {
	return wwb.widget
}

func (wwb *WidgetWithBorder[T]) SetRadius(r int) {
	wwb.radius = r
}

func (wwb *WidgetWithBorder[T]) SetBorderType(borderType basicwidgetdraw.RoundedRectBorderType) {
	wwb.borderType = borderType
}

func (wwb *WidgetWithBorder[T]) SetBorderWidth(borderWidth float32) {
	wwb.borderWidth = borderWidth
}
