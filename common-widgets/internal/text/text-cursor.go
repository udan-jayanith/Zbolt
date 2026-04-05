package text

import (
	"API-Client/common-widgets/internal/text/internal/textutil"

	draw "API-Client/common-widgets/internal/draw"

	"github.com/guigui-gui/guigui"
	"github.com/guigui-gui/guigui/basicwidget/basicwidgetdraw"
	"github.com/hajimehoshi/ebiten/v2"
)

type textCursor struct {
	guigui.DefaultWidget

	text *Text

	counter   int
	prevAlpha float64
	prevPos   textutil.TextPosition
	prevOK    bool
}

func (t *textCursor) resetCounter() {
	t.counter = 0
}

func (t *textCursor) Tick(context *guigui.Context, widgetBounds *guigui.WidgetBounds) error {
	pos, ok := t.text.cursorPosition(context, widgetBounds)
	if t.prevPos != pos {
		t.resetCounter()
	}
	t.prevPos = pos
	t.prevOK = ok

	t.counter++
	if a := t.alpha(context, widgetBounds, t.text); t.prevAlpha != a {
		t.prevAlpha = a
		guigui.RequestRedraw(t)
	}
	return nil
}

func (t *textCursor) alpha(context *guigui.Context, widgetBounds *guigui.WidgetBounds, text *Text) float64 {
	if _, ok := text.cursorPosition(context, widgetBounds); !ok {
		return 0
	}
	s, e, ok := text.selectionToDraw(context)
	if !ok {
		return 0
	}
	if s != e {
		return 0
	}
	if text.cursorStatic {
		return 1
	}
	offset := ebiten.TPS() / 2
	if t.counter <= offset {
		return 1
	}
	interval := ebiten.TPS()
	c := (t.counter - offset) % interval
	if c < interval/5 {
		return 1 - float64(c)/float64(interval/5)
	}
	if c < interval*2/5 {
		return 0
	}
	if c < interval*3/5 {
		return float64(c-interval*2/5) / float64(interval/5)
	}
	return 1
}

func (t *textCursor) Draw(context *guigui.Context, widgetBounds *guigui.WidgetBounds, dst *ebiten.Image) {
	alpha := t.alpha(context, widgetBounds, t.text)
	if alpha == 0 {
		return
	}
	b := widgetBounds.Bounds()
	clr := draw.ScaleAlpha(draw.Color2(context.ColorMode(), draw.SemanticColorAccent, 0.5, 0.6), alpha)
	basicwidgetdraw.DrawRoundedRect(context, dst, b, clr, b.Dx()/2)
}
