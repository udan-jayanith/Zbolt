package CommonWidgets

import (
	"API-Client/basic"
	"image"
	"image/color"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	draw "github.com/guigui-gui/guigui/basicwidget/basicwidgetdraw"
	"github.com/hajimehoshi/ebiten/v2"
)

type path_segment_widget struct {
	gui.DefaultWidget

	path_name string
	is_end    bool

	text_widget widget.Text
}

func (psw *path_segment_widget) Build(context *gui.Context, adder *gui.ChildAdder) error {
	text := psw.path_name
	if psw.is_end {
		text += "/"
	}

	text_widget := &psw.text_widget
	text_widget.SetTabular(true)
	text_widget.SetVerticalAlign(widget.VerticalAlignMiddle)
	text_widget.SetHorizontalAlign(widget.HorizontalAlignCenter)
	adder.AddWidget(text_widget)
	return nil
}

func (psw *path_segment_widget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Padding:   basic.NewPadding(widget.UnitSize(ctx) / 4),
		Items: []gui.LinearLayoutItem{
			{
				Widget: &psw.text_widget,
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (psw *path_segment_widget) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	point := psw.Measure(ctx, constraints)
	padding := widget.UnitSize(ctx) / 4
	point.X += padding * 2
	point.Y += padding * 2
	return point
}

func (psw *path_segment_widget) Draw(ctx *gui.Context, widgetBounds *gui.WidgetBounds, dst *ebiten.Image) {
	var background_color color.Color
	cm := ctx.ColorMode()

	if widgetBounds.IsHitAtCursor() {
		background_color = draw.BackgroundSecondaryColor(cm)
	} else {
		background_color = draw.BackgroundColor(cm)
	}

	draw.DrawRoundedRect(ctx, dst, widgetBounds.Bounds(), background_color, widget.UnitSize(ctx)/4)
}

type path_widget struct {
	gui.DefaultWidget

	segments []path_segment_widget
}

func (pw *path_widget) Build(context *gui.Context, adder *gui.ChildAdder) error {
	for i, _ := range pw.segments {
		segment := &pw.segments[i]
		adder.AddWidget(segment)
	}
	return nil
}

func (pw *path_widget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Gap:       widget.UnitSize(ctx) / 4,
		Items:     make([]gui.LinearLayoutItem, 0, len(pw.segments)),
	}

	for i, _ := range pw.segments {
		segment := &pw.segments[i]
		layout.Items = append(layout.Items, gui.LinearLayoutItem{
			Widget: segment,
		})
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (pw *path_widget) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	var point image.Point

	padding := widget.UnitSize(ctx) / 4
	point.Y = widget.LineHeight(ctx) + padding*2

	for i, _ := range pw.segments {
		measurements := pw.segments[i].Measure(ctx, constraints)
		point.X += measurements.X
	}

	gaps := (len(pw.segments) - 1) * (widget.UnitSize(ctx) / 4)
	point.X += gaps

	return point
}

type Path struct {
	gui.DefaultWidget

	panel       widget.Panel
	path_widget path_widget
}

func (path *Path) Build(context *gui.Context, adder *gui.ChildAdder) error {
	path.panel.SetContent(&path.path_widget)
	path.panel.SetAutoBorder(true)
	path.panel.SetContentConstraints(widget.PanelContentConstraintsFixedWidth)
	path.panel.SetStyle(widget.PanelStyleSide)
	adder.AddWidget(&path.panel)
	return nil
}

func (path *Path) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layouter.LayoutWidget(&path.panel, widgetBounds.Bounds())
}

func (path *Path) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	var point image.Point

	measurements := path.path_widget.Measure(ctx, gui.Constraints{})
	if w, ok := constraints.FixedWidth(); ok {
		point.X = w
	} else {
		point.X = measurements.X
	}

	if h, ok := constraints.FixedHeight(); ok {
		point.Y = h
	} else {
		point.Y = measurements.Y
	}
	
	return point
}
