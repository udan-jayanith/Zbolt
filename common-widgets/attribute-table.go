package CommonWidgets

import (
	"API-Client/basic"
	"image"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type Attribute struct {
	gui.DefaultWidget

	Key, Value                             string
	key_widget, value_widget               gui.WidgetWithPadding[*widget.Text]
	key_widget_border, value_widget_border WidgetWithBorder[*gui.WidgetWithPadding[*widget.Text]]
	delete_widget                          widget.Image
	is_bold                                bool
	Editable                               bool
}

func (attr *Attribute) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	padding := basic.NewPadding(0, widget.UnitSize(ctx)/3)

	attr.key_widget.SetPadding(padding)
	key_widget := attr.key_widget.Widget()
	key_widget.SetTabular(true)
	key_widget.SetEditable(attr.Editable)
	key_widget.SetVerticalAlign(widget.VerticalAlignMiddle)
	key_widget.SetHorizontalAlign(widget.HorizontalAlignLeft)
	key_widget.SetValue(attr.Key)
	key_widget.SetSelectable(true)
	key_widget.SetBold(attr.is_bold)
	key_widget.SetOnValueChanged(func(context *gui.Context, text string, committed bool) {
		attr.Key = text
	})

	attr.key_widget_border.SetWidget(&attr.key_widget)
	adder.AddChild(&attr.key_widget_border)

	attr.value_widget.SetPadding(padding)
	value_widget := attr.value_widget.Widget()
	value_widget.SetTabular(true)
	value_widget.SetEditable(attr.Editable)
	value_widget.SetVerticalAlign(widget.VerticalAlignMiddle)
	value_widget.SetHorizontalAlign(widget.HorizontalAlignLeft)
	value_widget.SetSelectable(true)
	value_widget.SetValue(attr.Value)
	value_widget.SetBold(attr.is_bold)
	value_widget.SetOnValueChanged(func(context *gui.Context, text string, committed bool) {
		attr.Value = text
	})

	attr.value_widget_border.SetWidget(&attr.value_widget)
	adder.AddChild(&attr.value_widget_border)

	return nil
}

func (attr *Attribute) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Items: []gui.LinearLayoutItem{
			{
				Widget: &attr.key_widget_border,
				Size:   gui.FlexibleSize(1),
			},
			{
				Widget: &attr.value_widget_border,
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
	if i == -1 || at.rows[i].Key != "" {
		at.rows = append(at.rows, &Attribute{
			Editable: true,
		})
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
			Size:   gui.FixedSize(widget.LineHeight(ctx)),
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
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Items: []gui.LinearLayoutItem{
			{
				Widget: &at.panel,
				Size:   gui.FlexibleSize(1),
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (at *AttributeTable) SetHeader(column1, column2 string) {
	attribute_table := &at.attribute_table
	attribute_table.header.Key = column1
	attribute_table.header.Value = column1
	attribute_table.header.is_bold = true
}

func (at *AttributeTable) AppendRows(rows []*Attribute) {
	attribute_table := &at.attribute_table
	attribute_table.rows = append(attribute_table.rows, rows...)
}
