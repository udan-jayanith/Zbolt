package CommonWidgets

import (
	"API-Client/basic"
	"image"
	"time"

	"sync"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/guigui-gui/guigui/basicwidget/basicwidgetdraw"
	"github.com/hajimehoshi/ebiten/v2"
)

type Notify struct {
	gui.DefaultWidget

	text_widget     widget.Text
	layered_padding gui.LayerWidget[*WidgetWithPadding[*widget.Text]]

	mutex                 sync.Mutex
	count                 int
	is_open               bool
	is_background_running bool
}

func (notify *Notify) background() {
	if notify.is_background_running {
		return
	}
	notify.is_background_running = true

	go func() {
		for {
			time.Sleep(time.Second)

			notify.mutex.Lock()
			notify.count--
			if notify.count <= 0 {
				notify.is_background_running = false
				notify.is_open = false
				notify.mutex.Unlock()
				break
			}
			notify.mutex.Unlock()
		}
	}()
}

func (notify *Notify) Open() {
	notify.mutex.Lock()
	defer notify.mutex.Unlock()

	notify.is_open = true
	notify.count++
	notify.background()
}

func (notify *Notify) IsOpen() bool {
	notify.mutex.Lock()
	defer notify.mutex.Unlock()
	return notify.is_open
}

func (notify *Notify) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	if !notify.IsOpen() {
		return nil
	}

	notify.text_widget.SetVerticalAlign(widget.VerticalAlignMiddle)
	notify.text_widget.SetHorizontalAlign(widget.HorizontalAlignCenter)
	notify.text_widget.SetOpacity(0.9)

	line_height := widget.LineHeight(ctx)
	notify.layered_padding.Widget().SetPadding(basic.NewPadding(line_height/3, line_height/2+line_height/3))
	notify.layered_padding.Widget().SetWidget(&notify.text_widget)
	adder.AddChild(&notify.layered_padding)
	return nil
}

func (notify *Notify) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	if notify.IsOpen() {
		layouter.LayoutWidget(&notify.layered_padding, widgetBounds.Bounds())
	}
}

func (notify *Notify) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	if !notify.IsOpen() {
		return image.Pt(0, 0)
	}
	return notify.layered_padding.Widget().Measure(ctx, constraints)
}

func (notify *Notify) Draw(ctx *gui.Context, widgetBounds *gui.WidgetBounds, dst *ebiten.Image) {
	if !notify.IsOpen() {
		return
	}
	background_color := basicwidgetdraw.BackgroundSecondaryColor(ctx.ColorMode())
	basicwidgetdraw.DrawRoundedRect(ctx, dst, widgetBounds.Bounds(), background_color, widget.LineHeight(ctx))
}

func (notify *Notify) LayoutWidget(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	if !notify.IsOpen() {
		return
	}

	measurements := notify.Measure(ctx, gui.Constraints{})
	b := widgetBounds.Bounds()

	notify_bounds := image.Rectangle{
		Min: image.Point{
			X: (b.Min.X + b.Dx()/2) - measurements.X/2,
			Y: b.Max.Y - (measurements.Y * 2),
		},
	}

	notify_bounds.Max = image.Point{
		X: notify_bounds.Min.X + measurements.X,
		Y: notify_bounds.Min.Y + measurements.Y,
	}

	layouter.LayoutWidget(notify, notify_bounds)
}

func (notify *Notify) SetText(text string) {
	notify.mutex.Lock()
	notify.text_widget.SetValue(text)
	notify.mutex.Unlock()
}
