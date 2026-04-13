package CommonWidgets

import (
	"fmt"
	"image"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/guigui-gui/guigui/basicwidget/basicwidgetdraw"
	"github.com/hajimehoshi/ebiten/v2"
)

type LoadingBar struct {
	gui.DefaultWidget

	show_percentage   bool
	total, downloaded int
	parentage_widget  widget.Text
	tooltip           widget.TooltipArea
}

func (wi *LoadingBar) ShowPercentage(show bool) {
	wi.show_percentage = show
}

func (wi *LoadingBar) SetValue(total, downloaded int) {
	wi.total = total
	wi.downloaded = downloaded
}

func (wi *LoadingBar) percentage() float32 {
	if wi.downloaded == 0 {
		return 0
	}
	return float32(wi.downloaded / wi.total * 100)
}

func (wi *LoadingBar) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	if !wi.show_percentage {
		return nil
	}

	wi.parentage_widget.SetValue(fmt.Sprintf("%v%%", wi.percentage()))
	wi.parentage_widget.SetVerticalAlign(widget.VerticalAlignMiddle)
	wi.parentage_widget.SetScale(0.4)
	adder.AddWidget(&wi.parentage_widget)

	wi.tooltip.SetText(fmt.Sprintf("%v of %v (%v%%)", wi.downloaded, wi.total, wi.percentage()))
	adder.AddWidget(&wi.tooltip)
	return nil
}

func (wi *LoadingBar) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	if !wi.show_percentage {
		return
	}

	b := widgetBounds.Bounds()
	middle := b.Min.X + b.Dx()/2

	percentage_size := wi.parentage_widget.Measure(ctx, gui.Constraints{})
	percentage_bounds := image.Rectangle{
		Min: image.Point{
			X: middle - percentage_size.X/2,
			Y: b.Min.Y,
		},
	}
	percentage_bounds.Max = image.Point{
		X: percentage_bounds.Min.X + percentage_size.X,
		Y: b.Max.Y,
	}

	layouter.LayoutWidget(&wi.parentage_widget, percentage_bounds)
	layouter.LayoutWidget(&wi.tooltip, b)
}

func (wi *LoadingBar) Draw(ctx *gui.Context, widgetBounds *gui.WidgetBounds, dst *ebiten.Image) {
	background_color := basicwidgetdraw.BackgroundSecondaryColor(ctx.ColorMode())
	basicwidgetdraw.DrawRoundedRect(ctx, dst, widgetBounds.Bounds(), background_color, 4)

	if wi.show_percentage {
		return
	}
}

func (wi *LoadingBar) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	w, ok := constraints.FixedWidth()
	if !ok {
		w = widget.UnitSize(ctx) * 6
	}

	return image.Pt(w, widget.LineHeight(ctx)/2)
}
