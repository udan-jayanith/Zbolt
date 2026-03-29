package CommonWidgets

import (
	"API-Client/basic"
	draw_color "API-Client/common-widgets/internal/draw"
	"API-Client/icons"
	"image"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type table_row_widget struct {
	gui.DefaultWidget

	checkbox                 widget.Checkbox
	key_column, value_column EditableText
	vr                       VerticalLine
	row_delete_btn           *icons.Icon
}

func (w *table_row_widget) padding(ctx *gui.Context) gui.Padding {
	return basic.NewPadding(widget.UnitSize(ctx) / 3)
}

func (w *table_row_widget) gap(ctx *gui.Context) int {
	return widget.UnitSize(ctx) / 4
}

func (w *table_row_widget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	adder.AddWidget(&w.checkbox)
	adder.AddWidget(&w.key_column)
	adder.AddWidget(&w.vr)
	adder.AddWidget(&w.value_column)

	if w.row_delete_btn == nil {
		l := widget.LineHeight(ctx)
		w.row_delete_btn = icons.NewIcon("delete", l-(l/3))
	}
	adder.AddWidget(w.row_delete_btn)
	return nil
}

func (w *table_row_widget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Gap:       w.gap(ctx),
		Padding:   w.padding(ctx),
		Items: []gui.LinearLayoutItem{},
	}

	left_column_layout := layout
	left_column_layout.Padding = gui.Padding{}
	left_column_layout.Items = []gui.LinearLayoutItem{
		{
			Widget: &w.checkbox,
		},
		{
			Widget: &w.key_column,
			Size:   gui.FlexibleSize(1),
		},
	}

	right_column_layout := layout
	right_column_layout.Padding = gui.Padding{}
	right_column_layout.Items = []gui.LinearLayoutItem{
		{
			Widget: &w.value_column,
			Size:   gui.FlexibleSize(1),
		},
		{
			Widget: w.row_delete_btn,
		},
	}

	layout.Items = []gui.LinearLayoutItem{
		{
			Layout: left_column_layout,
			Size: gui.FlexibleSize(1),
		},
		{
			Widget: &w.vr,
		},
		{
			Layout: right_column_layout,
			Size: gui.FlexibleSize(1),
		},
	}
	
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (row_widget *table_row_widget) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	var point image.Point
	padding := row_widget.padding(ctx)
	point.Y = widget.LineHeight(ctx) + padding.Top + padding.Bottom

	if w, ok := constraints.FixedWidth(); ok {
		point.X = w
	} else {
		point.X = row_widget.checkbox.Measure(ctx, gui.Constraints{}).X
		point.X += row_widget.key_column.Measure(ctx, gui.Constraints{}).X * 2
		l := widget.LineHeight(ctx)
		point.X += l - (l / 3)
		point.X += row_widget.gap(ctx) * 4
		point.X = padding.Start + padding.End
	}

	return point
}

func (w *table_row_widget) HandlePointingInput(ctx *gui.Context, widgetBounds *gui.WidgetBounds) gui.HandleInputResult {
	return gui.HandleInputResult{}
}

type attribute_table struct {
	gui.DefaultWidget
	vr                       VerticalLine
	key_header, value_header widget.Text
	rows                     []table_row_widget
}

func (at *attribute_table) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	at.key_header.SetValue("Key")
	at.key_header.SetBold(true)
	at.key_header.SetVerticalAlign(widget.VerticalAlignMiddle)
	adder.AddWidget(&at.key_header)

	adder.AddWidget(&at.vr)

	at.value_header.SetValue("Value")
	at.value_header.SetBold(true)
	at.value_header.SetVerticalAlign(widget.VerticalAlignMiddle)
	adder.AddWidget(&at.key_header)

	for i, _ := range at.rows {
		adder.AddWidget(&at.rows[i])
	}

	return nil
}

func (at *attribute_table) header_height(ctx *gui.Context) int {
	return widget.LineHeight(ctx) + (widget.UnitSize(ctx)/3)*2
}

func (at *attribute_table) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Items:     make([]gui.LinearLayoutItem, 0, len(at.rows)),
	}

	u := widget.UnitSize(ctx)
	header_layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Gap:       u / 4,
		Padding:   basic.NewPadding(u / 3),
		Items: []gui.LinearLayoutItem{
			{
				Widget: &at.key_header,
				Size:   gui.FlexibleSize(1),
			},
			{
				Widget: &at.vr,
			},
			{
				Widget: &at.value_header,
				Size:   gui.FlexibleSize(1),
			},
		},
	}
	layout.Items = append(layout.Items, gui.LinearLayoutItem{
		Layout: header_layout,
		Size:   gui.FixedSize(at.header_height(ctx)),
	})

	for i, _ := range at.rows {
		layout.Items = append(layout.Items, gui.LinearLayoutItem{
			Widget: &at.rows[i],
		})
	}

	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (at *attribute_table) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	var point image.Point

	if w, ok := constraints.FixedWidth(); ok {
		point.X = w
	} else {
		point.X = widget.UnitSize(ctx) * 6
	}

	point.Y = at.header_height(ctx)
	for i, _ := range at.rows {
		point.Y += at.rows[i].Measure(ctx, gui.Constraints{}).Y
	}

	return point
}

func (at *attribute_table) Draw(ctx *gui.Context, widgetBounds *gui.WidgetBounds, dst *ebiten.Image) {
	b := widgetBounds.Bounds()
	b.Min.Y += at.header_height(ctx)

	line_color := draw_color.ScaleAlpha(draw_color.Color(ctx.ResolvedColorMode(), draw_color.ColorTypeBase, 0), 6/32.0)
	width := 1 * float32(ctx.Scale())

	for i, _ := range at.rows {
		vector.StrokeLine(dst, float32(b.Min.X), float32(b.Min.Y), float32(b.Max.X), float32(b.Min.Y), width, line_color, true)
		b.Min.Y += at.rows[i].Measure(ctx, gui.Constraints{}).Y
	}
}

type AttributeTable struct {
	gui.DefaultWidget
	
	table attribute_table
	panel widget.Panel
}

func (table *AttributeTable) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	table.panel.SetContentConstraints(widget.PanelContentConstraintsFixedWidth)
	table.panel.SetContent(&table.table)
	return nil
}

func (table *AttributeTable) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layouter.LayoutWidget(&table.panel, widgetBounds.Bounds())
}

func (t *AttributeTable) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	return image.Pt(12*widget.UnitSize(ctx), 6*widget.UnitSize(ctx))
}