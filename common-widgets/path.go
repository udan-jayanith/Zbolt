package CommonWidgets

import (
	"image"
	"image/color"
	"path/filepath"
	"strings"

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
	if !psw.is_end {
		text += "/"
	}

	text_widget := &psw.text_widget
	text_widget.SetTabular(true)
	text_widget.SetVerticalAlign(widget.VerticalAlignMiddle)
	text_widget.SetHorizontalAlign(widget.HorizontalAlignCenter)
	text_widget.SetValue(text)
	adder.AddWidget(text_widget)
	return nil
}

func (psw *path_segment_widget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Items: []gui.LinearLayoutItem{
			{
				Widget: &psw.text_widget,
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (psw *path_segment_widget) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	return psw.text_widget.Measure(ctx, gui.Constraints{})
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
	point.Y = widget.LineHeight(ctx)

	for i, _ := range pw.segments {
		measurements := pw.segments[i].Measure(ctx, constraints)
		point.X += measurements.X
	}
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
	adder.AddWidget(&path.panel)
	return nil
}

func (path *Path) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Items: []gui.LinearLayoutItem{
			{
				Widget: &path.panel,
				Size: gui.FlexibleSize(1),
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
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

func (path_widget *Path) SetPath(directory_path string) {
	directory_path = filepath.Clean(directory_path)
	list := strings.Split(filepath.ToSlash(directory_path), "/")
	l := len(list)
	path_widget.path_widget.segments = make([]path_segment_widget, 0, l)
	for i, path_name := range list {
		path_widget.path_widget.segments = append(path_widget.path_widget.segments, path_segment_widget{
			path_name: path_name,
			is_end: l == i,
		})
	}
}
