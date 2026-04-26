package message_model

import (
	"API-Client/basic"
	"image"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/guigui-gui/guigui/basicwidget/basicwidgetdraw"
	"github.com/hajimehoshi/ebiten/v2"
)

type notify_widget struct {
	gui.DefaultWidget

	text_widget gui.WidgetWithPadding[*widget.Text]
	on_result   result_fn_type
}

func (w *notify_widget) SetMessage(message string) {
	w.text_widget.Widget().SetValue(message)
}

func (w *notify_widget) OnResult(fn result_fn_type) {
	w.on_result = fn
}

// TODO: close the notify_widget after X amount of time
func (notify_widget *notify_widget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	text_widget := notify_widget.text_widget.Widget()
	text_widget.SetVerticalAlign(widget.VerticalAlignMiddle)
	text_widget.SetHorizontalAlign(widget.HorizontalAlignCenter)
	text_widget.SetOpacity(0.9)

	line_height := widget.LineHeight(ctx)
	notify_widget.text_widget.SetPadding(basic.NewPadding(line_height/3, line_height/2+line_height/3))
	adder.AddWidget(&notify_widget.text_widget)
	return nil
}

func (notify_widget *notify_widget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layouter.LayoutWidget(&notify_widget.text_widget, widgetBounds.Bounds())
}

func (notify_widget *notify_widget) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	return notify_widget.text_widget.Widget().Measure(ctx, constraints)
}

func (notify_widget *notify_widget) Draw(ctx *gui.Context, widgetBounds *gui.WidgetBounds, dst *ebiten.Image) {
	background_color := basicwidgetdraw.BackgroundSecondaryColor(ctx.ColorMode())
	basicwidgetdraw.DrawRoundedRect(ctx, dst, widgetBounds.Bounds(), background_color, widget.LineHeight(ctx))
}

func (notify_widget *notify_widget) Bounds(ctx *gui.Context, widgetBounds *gui.WidgetBounds) image.Rectangle {
	measurements := notify_widget.Measure(ctx, gui.Constraints{})
	b := widgetBounds.Bounds()

	notify_widget_bounds := image.Rectangle{
		Min: image.Point{
			X: (b.Min.X + b.Dx()/2) - measurements.X/2,
			Y: b.Max.Y - (measurements.Y * 2),
		},
	}

	notify_widget_bounds.Max = image.Point{
		X: notify_widget_bounds.Min.X + measurements.X,
		Y: notify_widget_bounds.Min.Y + measurements.Y,
	}

	return notify_widget_bounds
}
