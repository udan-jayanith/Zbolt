package http_widget

import (
	CommonWidgets "API-Client/common-widgets"
	"image"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type query_table_widget struct {
	gui.DefaultWidget

	table widget.Table[struct{}]
	items []widget.TableRow[struct{}]
}

func (w *query_table_widget) push_row(name, value string) {
	v := &CommonWidgets.EditableText{}
	v.SetValue(value)
	w.items = append(w.items, widget.TableRow[struct{}]{
		Cells: []widget.TableCell{
			{
				Text: name,
			},
			{
				Content: v,
			},
		},
	})
	w.table.SetItems(w.items)
}

func (w *query_table_widget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	w.table.SetColumns([]widget.TableColumn{
		{
			Width: gui.FlexibleSize(1),
			HeaderText: "Name",
		},
		{
			Width: gui.FlexibleSize(1),
			HeaderText: "Value",
		},
	})

	adder.AddWidget(&w.table)
	return nil
}

func (w *query_table_widget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layouter.LayoutWidget(&w.table, widgetBounds.Bounds())
}

func (w *query_table_widget) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	return w.table.Measure(ctx, constraints)
}

func (w *query_table_widget) GetValues() []string {
	values := make([]string, 0, len(w.items))
	for _, cell := range w.items {
		w, ok := cell.Cells[1].Content.(*CommonWidgets.EditableText)
		if !ok {
			panic("Unexpected widget")
		}
		values = append(values, w.Value())
	}
	return values
}

func (w *query_table_widget) Empty(){
	w.items = make([]widget.TableRow[struct{}], 0, len(w.items))
}
