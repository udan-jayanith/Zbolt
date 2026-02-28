package CommonWidgets

import (
	"API-Client/basic"
	"API-Client/icons"
	"image"

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

	panel widget.Panel
}

func (at *AttributeTable) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	at.table.SetColumns([]widget.TableColumn{
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
		last_key_column := at.table_rows[no_of_rows-1].Cells[0].Content.(*gui.WidgetWithPadding[*EditableText])
		if last_key_column.Widget().widget.Value() != "" {
			ok = true
		}
	}

	if ok || no_of_rows == 0 {
		cell1, cell2 := &gui.WidgetWithPadding[*EditableText]{}, &gui.WidgetWithPadding[*EditableText]{}
		u := widget.UnitSize(ctx)
		padding := basic.NewPadding(0, u/3)
		cell1.SetPadding(padding)
		cell2.SetPadding(padding)

		at.table_rows = append(at.table_rows, widget.TableRow[struct{}]{
			Movable: true,
			Cells: []widget.TableCell{
				{
					Content: cell1,
				},
				{
					Content: cell2,
				},
				{
					Content: icons.NewIcon("delete", widget.LineHeight(ctx)),
				},
			},
		})
	}

	at.table.SetItems(at.table_rows)
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
