package CommonWidgets

import (
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

	max, min := 24.0, 16.0
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
	return nil
}

func (lazy_loading *LazyLoading) Draw(ctx *gui.Context, widgetBounds *gui.WidgetBounds, dst *ebiten.Image) {
	w := lazy_loading.whiteness / 100
	clr := iro.ColorFromSRGB(w, w, w, 1).SRGBColor()
	basicwidgetdraw.DrawRoundedRect(ctx, dst, widgetBounds.Bounds(), clr, widget.RoundedCornerRadius(ctx))
}

// TODO: implement widget with lazy loading