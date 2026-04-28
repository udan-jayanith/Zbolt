package CommonWidgets

import (
	"image"
	"time"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/guigui-gui/guigui/basicwidget/basicwidgetdraw"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/iro"
)

type LazyLoading struct {
	gui.DefaultWidget
	t time.Time

	is_decreasing bool
	whiteness     float64
}

func (lazy_loading *LazyLoading) Tick(ctx *gui.Context, widgetBounds *gui.WidgetBounds) error {
	if time.Since(lazy_loading.t).Milliseconds() < 30 {
		return nil
	}

	max, min := 30.0, 18.0
	if lazy_loading.whiteness == 0 {
		lazy_loading.whiteness = min
	}

	if lazy_loading.is_decreasing {
		lazy_loading.whiteness -= 0.2
	} else {
		lazy_loading.whiteness += 0.1
	}

	if lazy_loading.whiteness >= max {
		lazy_loading.is_decreasing = true
	} else if lazy_loading.whiteness <= min {
		lazy_loading.is_decreasing = false
	}
	lazy_loading.t = time.Now()

	gui.RequestRedraw(lazy_loading)
	return nil
}

func (lazy_loading *LazyLoading) Draw(ctx *gui.Context, widgetBounds *gui.WidgetBounds, dst *ebiten.Image) {
	w := lazy_loading.whiteness / 100
	clr := iro.ColorFromSRGB(w, w, w, 1).SRGBColor()
	basicwidgetdraw.DrawRoundedRect(ctx, dst, widgetBounds.Bounds(), clr, widget.RoundedCornerRadius(ctx))
}

type WidgetWithLazyLoading[T gui.Widget] struct {
	gui.DefaultWidget

	widget       lazy_widget[T]
	lazy_loading LazyLoading
	lazy_load    bool
}

func (wi *WidgetWithLazyLoading[T]) Widget() T {
	return wi.widget.Widget()
}

func (wi *WidgetWithLazyLoading[T]) SetWidget(widget T) {
	wi.widget.SetWidget(widget)
}

func (wi *WidgetWithLazyLoading[T]) LazyLoad() bool {
	return wi.lazy_load
}

func (wi *WidgetWithLazyLoading[T]) SetLazyLoad(lazy_load bool) {
	wi.lazy_load = lazy_load
}

func (wi *WidgetWithLazyLoading[T]) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	if wi.lazy_load {
		adder.AddWidget(&wi.lazy_loading)
		return nil
	}
	adder.AddWidget(wi.widget.Widget())
	return nil
}

func (wi *WidgetWithLazyLoading[T]) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	if wi.lazy_load {
		layouter.LayoutWidget(&wi.lazy_loading, widgetBounds.Bounds())
		return
	}
	layouter.LayoutWidget(wi.widget.Widget(), widgetBounds.Bounds())
}

func (wi *WidgetWithLazyLoading[T]) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	return wi.widget.Widget().Measure(ctx, constraints)
}
