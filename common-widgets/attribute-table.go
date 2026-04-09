package CommonWidgets

import (
	"API-Client/basic"
	draw_color "API-Client/common-widgets/internal/draw"
	"API-Client/icons"
	attr "API-Client/widgets/request/attributes"
	"image"
	"slices"
	"strings"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/guigui-gui/guigui/basicwidget/basicwidgetdraw"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type table_row_widget struct {
	gui.DefaultWidget

	table *attribute_table

	index                int
	checkbox             widget.Checkbox
	key_cell, value_cell EditableText
	vr                   VerticalLine
	row_delete_btn       *icons.Icon
}

func (w *table_row_widget) padding(ctx *gui.Context) gui.Padding {
	return basic.NewPadding(widget.UnitSize(ctx) / 8)
}

func (w *table_row_widget) gap(ctx *gui.Context) int {
	return widget.UnitSize(ctx) / 4
}

func (w *table_row_widget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	if !w.table.checkbox_disabled {
		adder.AddWidget(&w.checkbox)
	}

	w.key_cell.SetEditable(!w.table.key_not_editable)
	adder.AddWidget(&w.key_cell)
	adder.AddWidget(&w.vr)
	adder.AddWidget(&w.value_cell)

	if !w.table.delete_disabled {
		if w.row_delete_btn == nil {
			l := widget.LineHeight(ctx)
			w.row_delete_btn = icons.NewIcon("delete", l-(l/6))
		}
		adder.AddWidget(w.row_delete_btn)
	}
	return nil
}

func (w *table_row_widget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Gap:       w.gap(ctx),
		Padding:   w.padding(ctx),
		Items:     []gui.LinearLayoutItem{},
	}

	left_column_layout := layout
	left_column_layout.Padding = gui.Padding{}
	left_column_layout.Items = []gui.LinearLayoutItem{
		{
			Widget: &w.checkbox,
		},
		{
			Widget: &w.key_cell,
			Size:   gui.FlexibleSize(1),
		},
	}
	if w.table.checkbox_disabled {
		left_column_layout.Items = left_column_layout.Items[1:]
	}

	right_column_layout := layout
	right_column_layout.Padding = gui.Padding{}
	right_column_layout.Items = []gui.LinearLayoutItem{
		{
			Widget: &w.value_cell,
			Size:   gui.FlexibleSize(1),
		},
	}
	if !w.table.delete_disabled {
		right_column_layout.Items = append(right_column_layout.Items, gui.LinearLayoutItem{
			Widget: w.row_delete_btn,
		})
	}

	layout.Items = []gui.LinearLayoutItem{
		{
			Layout: left_column_layout,
			Size:   gui.FlexibleSize(1),
		},
		{
			Widget: &w.vr,
		},
		{
			Layout: right_column_layout,
			Size:   gui.FlexibleSize(1),
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
		point.X = widget.UnitSize(ctx) * 6
	}

	return point
}

func (w *table_row_widget) HandlePointingInput(ctx *gui.Context, widgetBounds *gui.WidgetBounds) gui.HandleInputResult {
	return gui.HandleInputResult{}
}

func (row_widget *table_row_widget) on_delete(fn func(index int)) {
	if row_widget.row_delete_btn == nil {
		return
	}
	row_widget.row_delete_btn.OnClick(func() {
		fn(row_widget.index)
	})
}

type attribute_table struct {
	gui.DefaultWidget
	vr                       VerticalLine
	key_header, value_header widget.Text
	rows                     []*table_row_widget

	disable_auto_add                   bool
	checkbox_disabled, delete_disabled bool
	key_not_editable                   bool
	rwo_delete_fn                      func(index int)
}

func (at *attribute_table) push_row(row attr.AttrCheck) {
	row_widget := table_row_widget{}
	row_widget.index = len(at.rows)
	row_widget.table = at
	row_widget.checkbox.SetValue(row.Checked)
	row_widget.key_cell.SetValue(row.Key)
	row_widget.value_cell.SetValue(row.Value)
	at.rows = append(at.rows, &row_widget)
}

func (at *attribute_table) delete_row(index int) {
	at.rows = slices.Delete(at.rows, index, index+1)
	gui.RequestRebuild(at)
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
	adder.AddWidget(&at.value_header)

	l := len(at.rows)
	if !at.disable_auto_add && (l == 0 || strings.TrimSpace(at.rows[l-1].key_cell.Value()) != "") {
		at.push_row(attr.AttrCheck{
			Checked: true,
		})
	}

	for i, _ := range at.rows {
		row_widget := at.rows[i]
		if !at.delete_disabled {
			row_widget.on_delete(at.delete_row)
		}
		row_widget.index = i
		adder.AddWidget(at.rows[i])
	}

	return nil
}

func (at *attribute_table) header_height(ctx *gui.Context) int {
	return widget.LineHeight(ctx) + (widget.UnitSize(ctx)/8)*2
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
		Padding:   basic.NewPadding(u / 8),
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
			Widget: at.rows[i],
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

	line_color := draw_color.ScaleAlpha(draw_color.Color(	ctx.ColorMode(), draw_color.ColorTypeBase, 0), 6/32.0)
	width := 1 * float32(ctx.Scale())

	for i, _ := range at.rows {
		vector.StrokeLine(dst, float32(b.Min.X), float32(b.Min.Y), float32(b.Max.X), float32(b.Min.Y), width, line_color, true)
		b.Min.Y += at.rows[i].Measure(ctx, gui.Constraints{}).Y
	}
}

type AttributeTable struct {
	gui.DefaultWidget
	table gui.WidgetWithPadding[*attribute_table]

	panel widget.Panel
}

func (table *AttributeTable) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	padding := basic.NewPadding(widget.UnitSize(ctx) / 3)
	padding.Top = 0
	table.table.SetPadding(padding)

	table.panel.SetContent(&table.table)
	table.panel.SetContentConstraints(widget.PanelContentConstraintsFixedWidth)
	adder.AddWidget(&table.panel)
	return nil
}

func (table *AttributeTable) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layouter.LayoutWidget(&table.panel, widgetBounds.Bounds())
}

func (t *AttributeTable) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	return image.Pt(12*widget.UnitSize(ctx), 6*widget.UnitSize(ctx))
}

func (t *AttributeTable) Draw(ctx *gui.Context, widgetBounds *gui.WidgetBounds, dst *ebiten.Image) {
	border_radius := widget.UnitSize(ctx) / 4

	background_clr := basicwidgetdraw.ControlColor(ctx.ColorMode(), ctx.IsEnabled(t))
	basicwidgetdraw.DrawRoundedRect(ctx, dst, widgetBounds.Bounds(), background_clr, border_radius)

	border_clr1, border_clr2 := basicwidgetdraw.BorderColors(ctx.ColorMode(), basicwidgetdraw.RoundedRectBorderTypeInset)
	basicwidgetdraw.DrawRoundedRectBorder(ctx, dst, widgetBounds.Bounds(), border_clr1, border_clr2, border_radius, 1, basicwidgetdraw.RoundedRectBorderTypeInset)
}

func (t *AttributeTable) SetRows(rows []attr.Attribute) {
	table_rows := t.table.Widget().rows
	if len(table_rows) > len(rows) {
		table_rows = table_rows[:len(rows)]
	} else if len(table_rows) != len(rows) {
		table_rows = make([]*table_row_widget, len(rows))
	}
	table := t.table.Widget()

	for i, row := range rows {
		if table_rows[i] == nil {
			table_rows[i] = &table_row_widget{}
		}
		table_row := table_rows[i]
		table_row.table = table
		table_row.index = i

		table_row.key_cell.SetValue(row.Key)
		table_row.value_cell.SetValue(row.Value)
	}
	t.table.Widget().rows = table_rows
}

func (t *AttributeTable) SetRowsCheck(rows []attr.AttrCheck) {
	table_rows := t.table.Widget().rows
	//TODO: BUG: optimizations doesn't work figure out what is happening.

	//if len(table_rows) > len(rows) {
	//table_rows = table_rows[:len(rows)]
	//} else if len(table_rows) != len(rows) {
	table_rows = make([]*table_row_widget, len(rows))
	//}
	table := t.table.Widget()

	for i, row := range rows {
		if table_rows[i] == nil {
			table_rows[i] = &table_row_widget{}
		}

		table_row := table_rows[i]
		table_row.table = table
		table_row.index = i

		table_row.checkbox.SetValue(row.Checked)
		table_row.key_cell.SetValue(row.Key)
		table_row.value_cell.SetValue(row.Value)
	}
	t.table.Widget().rows = table_rows
}

func (t *AttributeTable) RowsCheck() []attr.AttrCheck {
	table_rows := t.table.Widget().rows
	rows := make([]attr.AttrCheck, 0, len(table_rows))

	for _, table_row := range table_rows {
		if strings.TrimSpace(table_row.key_cell.Value()) == "" {
			continue
		}
		rows = append(rows, attr.AttrCheck{
			Key:     table_row.key_cell.Value(),
			Value:   table_row.value_cell.Value(),
			Checked: table_row.checkbox.Value(),
		})
	}

	return rows
}

func (t *AttributeTable) Rows() []attr.Attribute {
	table_rows := t.table.Widget().rows
	rows := make([]attr.Attribute, 0, len(table_rows))

	for _, table_row := range table_rows {
		if strings.TrimSpace(table_row.key_cell.Value()) == "" {
			continue
		}
		rows = append(rows, attr.Attribute{
			Key:   table_row.key_cell.Value(),
			Value: table_row.value_cell.Value(),
		})
	}

	return rows
}

func (t *AttributeTable) DisableCheckbox(disable bool) {
	t.table.Widget().checkbox_disabled = disable
}

func (t *AttributeTable) DisableDelete(disable bool) {
	t.table.Widget().delete_disabled = disable
}

func (t *AttributeTable) KeyEditable(editable bool) {
	t.table.Widget().key_not_editable = !editable
}

func (t *AttributeTable) AutoAddRow(auto_add bool) {
	t.table.Widget().disable_auto_add = !auto_add
}

func (t *AttributeTable) Count() int {
	return len(t.table.Widget().rows)
}

func (t *AttributeTable) PushRow(row attr.AttrCheck) {
	t.table.Widget().push_row(row)
}
