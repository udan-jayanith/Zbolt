package CommonWidgets

import (
	"API-Client/basic"
	"API-Client/icons"
	url_pattern "API-Client/widgets/request/url-pattern"
	"image"
	"slices"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
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
		Items: []gui.LinearLayoutItem{
			{
				Widget: &w.checkbox,
			},
			{
				Widget: &w.key_column,
				Size:   gui.FlexibleSize(1),
			},
			{
				Widget: &w.vr,
			},
			{
				Widget: &w.value_column,
				Size:   gui.FlexibleSize(1),
			},
			{
				Widget: w.row_delete_btn,
			},
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

/*
func (w *table_row_widget) Draw(ctx *gui.Context, widgetBounds *gui.WidgetBounds, dst *ebiten.Image) {
	//clr := draw_color.Color2(ctx.ResolvedColorMode(), draw_color.ColorTypeBase, 0.9, 0.4)
	//if !ctx.IsEnabled(w) {
		//clr = draw_color.Color2(ctx.ResolvedColorMode(), draw_color.ColorTypeBase, 0.8, 0.3)
		//	}

	//b := widgetBounds.Bounds()
}
*/

func (w *table_row_widget) HandlePointingInput(ctx *gui.Context, widgetBounds *gui.WidgetBounds) gui.HandleInputResult {
	return gui.HandleInputResult{}
}

type AttributeTable struct {
	gui.DefaultWidget

	table_rows        []widget.TableRow[struct{}]
	table             widget.Table[struct{}]
	checkbox_disabled bool

	panel widget.Panel
}

func (at *AttributeTable) find_row(key string) int {
	for i, row := range at.table_rows {
		row_key := row.Cells[1].Content.(*gui.WidgetWithPadding[*EditableText])
		if row_key.Widget().Value() == key {
			return i
		}
	}
	return -1
}

func (at *AttributeTable) DisableCheckBox(checkbox bool) {
	at.checkbox_disabled = checkbox
}

func (at *AttributeTable) delete_row(key string) {
	i := at.find_row(key)
	if i < 0 {
		return
	}
	at.table_rows = slices.Delete(at.table_rows, i, i+1)
}

func (at *AttributeTable) PushRow(attr url_pattern.Attribute, ctx *gui.Context) {
	cell1, cell2 := &gui.WidgetWithPadding[*EditableText]{}, &gui.WidgetWithPadding[*EditableText]{}

	cell1.Widget().widget.SetValue(attr.Key)
	cell2.Widget().widget.SetValue(attr.Value)

	u := widget.UnitSize(ctx)
	padding := basic.NewPadding(0, u/3)
	cell1.SetPadding(padding)
	cell2.SetPadding(padding)

	line_height := widget.LineHeight(ctx)

	checkbox := widget.Checkbox{}
	checkbox.SetValue(attr.Checked)
	ctx.SetEnabled(&checkbox, !at.checkbox_disabled)

	delete_icon := icons.NewIcon("delete", line_height/2)

	at.table_rows = append(at.table_rows, widget.TableRow[struct{}]{
		Movable: true,
		Cells: []widget.TableCell{
			{
				Content: &checkbox,
			},
			{
				Content: cell1,
			},
			{
				Content: cell2,
			},
			{
				Content: delete_icon,
			},
		},
	})

	delete_icon.OnClick(func() {
		at.delete_row(cell1.Widget().Value())
	})
}

func (at *AttributeTable) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	at.table.SetColumns([]widget.TableColumn{
		{
			Width: gui.FixedSize(widget.UnitSize(ctx)),
		},
		{
			HeaderText: "Key",
			Width:      gui.FlexibleSize(1),
		},
		{
			HeaderText: "Value",
			Width:      gui.FlexibleSize(1),
		},
		{
			Width: gui.FixedSize(widget.UnitSize(ctx)),
		},
	})

	no_of_rows := len(at.table_rows)
	var ok bool
	if no_of_rows > 0 {
		last_key_column := at.table_rows[no_of_rows-1].Cells[1].Content.(*gui.WidgetWithPadding[*EditableText])
		if last_key_column.Widget().widget.Value() != "" {
			ok = true
		}
	}

	if ok || no_of_rows == 0 {
		at.PushRow(url_pattern.Attribute{
			Checked: true,
		}, ctx)
	}

	at.table.SetItems(at.table_rows)
	at.table.OnItemsMoved(func(context *gui.Context, from, count, to int) {
		println("from", from)
		println("to", to)
		println()
	})
	adder.AddWidget(&at.table)
	return nil
}

func (at *AttributeTable) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Items: []gui.LinearLayoutItem{
			{
				Widget: &at.table,
				Size:   gui.FlexibleSize(1),
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (at *AttributeTable) Values() []url_pattern.Attribute {
	list := make([]url_pattern.Attribute, 0, len(at.table_rows))
	for _, row := range at.table_rows {
		is_checked := row.Cells[0].Content.(*widget.Checkbox)
		key := row.Cells[1].Content.(*gui.WidgetWithPadding[*EditableText]).Widget()
		value := row.Cells[2].Content.(*gui.WidgetWithPadding[*EditableText]).Widget()

		attr := url_pattern.Attribute{
			Checked: is_checked.Value(),

			Key:   key.Value(),
			Value: value.Value(),
		}
		list = append(list, attr)
	}

	return list
}
