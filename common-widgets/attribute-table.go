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

type EditableText struct {
	gui.DefaultWidget
	widget widget.Text
}

func (et *EditableText) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	et.widget.SetTabular(true)
	et.widget.SetVerticalAlign(widget.VerticalAlignMiddle)
	et.widget.SetHorizontalAlign(widget.HorizontalAlignLeft)
	et.widget.SetSelectable(true)
	et.widget.SetEditable(true)
	adder.AddWidget(&et.widget)
	return nil
}

func (et *EditableText) SetValue(text string) {
	et.widget.SetValue(text)
}

func (et *EditableText) Value() string {
	return et.widget.Value()
}

func (et *EditableText) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layouter.LayoutWidget(&et.widget, widgetBounds.Bounds())
}

func (et *EditableText) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	return et.widget.Measure(ctx, constraints)
}

type AttributeTable struct {
	gui.DefaultWidget

	table_rows []widget.TableRow[struct{}]
	table      widget.Table[struct{}]
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

func (at *AttributeTable) PushRow(cell1_text, cell2_text string, ctx *gui.Context) {
	cell1, cell2 := &gui.WidgetWithPadding[*EditableText]{}, &gui.WidgetWithPadding[*EditableText]{}

	cell1.Widget().widget.SetValue(cell1_text)
	cell2.Widget().widget.SetValue(cell2_text)

	u := widget.UnitSize(ctx)
	padding := basic.NewPadding(0, u/3)
	cell1.SetPadding(padding)
	cell2.SetPadding(padding)

	line_height := widget.LineHeight(ctx)

	checkbox := widget.Checkbox{}
	checkbox.SetValue(true)
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
		at.PushRow("", "", ctx)
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

func (at *AttributeTable) Values() []url_pattern.Attribute{
	return []url_pattern.Attribute{}
}