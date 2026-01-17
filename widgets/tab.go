package CWidget

import (
	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	basic "API-Client/basic"
)

type TabItem struct {
	gui.DefaultWidget
	text widget.Text
}

type TabBar struct {
	gui.DefaultWidget
	tab_items []*TabItem
}

func (tb *TabBar) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	for _, item := range tb.tab_items {
		adder.AddChild(item)
	}

	return nil
}

func (tb *TabBar) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Items: make([]gui.LinearLayoutItem, 0, len(tb.tab_items)),
	}
	
	for _, item := range tb.tab_items {
		layout.Items = append(layout.Items, gui.LinearLayoutItem{
			Widget: item,
		})
	}

	layout = basic.Align(layout, basic.Start, basic.Center)
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

type TabContainer struct {
	gui.DefaultWidget
	tab_bar *TabBar
}

