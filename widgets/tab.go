// Custom widgets
package CWidget

import (
	gui "github.com/guigui-gui/guigui"
)

type TabItem struct {
	index int
	Text string
	Body gui.Widget
}

type Tab struct {
	gui.DefaultWidget
	tab_items []TabItem
}

func (t *Tab) SetTabItems(tab_items []TabItem){
	t.tab_items = tab_items
}

func (t *Tab) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	for i := 0; i <len(t.tab_items); i++ {
		//t.tab_items[i].Text.SetValue("Hello world")
		//adder.AddChild(&t.tab_items[i].Text)
	}
	return nil
}

func (t *Tab) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Items: make([]gui.LinearLayoutItem, len(t.tab_items)),
	}
	for i := 0; i <len(t.tab_items); i++ {
		layout.Items[i] = gui.LinearLayoutItem{
		//	Widget: &t.tab_items[i].Text,
		}
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}