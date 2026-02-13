package CommonWidgets

import (
	"image"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type Attribute struct {
	gui.DefaultWidget

	Key, Value               string
	key_widget, value_widget widget.Text
	delete_widget            widget.Image
	is_bold                  bool
}

func (attr *Attribute) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	attr.key_widget.SetTabular(true)
	attr.key_widget.SetEditable(true)
	attr.key_widget.SetVerticalAlign(widget.VerticalAlignMiddle)
	attr.key_widget.SetHorizontalAlign(widget.HorizontalAlignLeft)
	attr.key_widget.SetValue(attr.Key)
	attr.key_widget.SetSelectable(true)
	attr.key_widget.SetBold(attr.is_bold)
	adder.AddChild(&attr.key_widget)

	attr.value_widget.SetTabular(true)
	attr.value_widget.SetEditable(true)
	attr.value_widget.SetVerticalAlign(widget.VerticalAlignMiddle)
	attr.value_widget.SetHorizontalAlign(widget.HorizontalAlignLeft)
	attr.value_widget.SetSelectable(true)
	attr.value_widget.SetValue(attr.Value)
	attr.value_widget.SetBold(attr.is_bold)
	adder.AddChild(&attr.value_widget)
	return nil
}

func (attr *Attribute) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Items: []gui.LinearLayoutItem{
			{
				Widget: &attr.key_widget,
				Size:   gui.FlexibleSize(1),
			},
			{
				Widget: &attr.value_widget,
				Size:   gui.FlexibleSize(1),
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

type attribute_table struct {
	gui.DefaultWidget
	header Attribute
	rows   []*Attribute
}

func (at *attribute_table) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	adder.AddChild(&at.header)
	
	i := len(at.rows) - 1
	if i == -1 || at.rows[i].Key != "" && at.rows[i].Value != "" {
		at.rows = append(at.rows, &Attribute{})
	}

	for _, attr_widget := range at.rows {
		adder.AddChild(attr_widget)
	}
	return nil
}

func (at *attribute_table) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Items:     make([]gui.LinearLayoutItem, 0, len(at.rows)+1),
	}

	u := widget.UnitSize(ctx)
	row_height := u + u/4
	
	layout.Items = append(layout.Items, gui.LinearLayoutItem{
		Widget: &at.header,
		Size:   gui.FixedSize(row_height),
	})

	for _, attr_widget := range at.rows {
		layout.Items = append(layout.Items, gui.LinearLayoutItem{
			Widget: attr_widget,
			Size:   gui.FixedSize(row_height),
		})
	}

	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (at *attribute_table) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	var point image.Point

	measurements := at.header.Measure(ctx, constraints)
	point.Y += measurements.Y
	point.X = measurements.X
	
	for _, attr_widget := range at.rows {
		measurements := attr_widget.Measure(ctx, constraints)
		point.Y += measurements.Y
	}

	return point
}

type AttributeTable struct {
	gui.DefaultWidget

	attribute_table attribute_table
	panel           widget.Panel
}

func (at *AttributeTable) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	at.panel.SetContentConstraints(widget.PanelContentConstraintsFixedWidth)
	at.panel.SetContent(&at.attribute_table)
	at.panel.SetStyle(widget.PanelStyleSide)
	adder.AddChild(&at.panel)
	return nil
}

func (at *AttributeTable) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layouter.LayoutWidget(&at.panel, widgetBounds.Bounds())
}

func (at *AttributeTable) SetHeader(column1, column2 string){
	at.attribute_table.header.Key = column1
	at.attribute_table.header.Value = column1
	at.attribute_table.header.is_bold = true
}

func (at *AttributeTable) AppendRows(rows []*Attribute){
	at.attribute_table.rows = append(at.attribute_table.rows, rows...) 
}