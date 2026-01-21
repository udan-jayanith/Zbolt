package CommonWidgets

import(
	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type TabItem[T comparable] struct {
	Name string
	Widget gui.Widget
	Value T
}

type Tab[T comparable] struct {
	gui.DefaultWidget
	
	tabs widget.SegmentedControl[string]
	Tab_Items []TabItem[T]
	panel widget.Panel
}

func (tab *Tab[T]) OnSelect(fn func(ctx *gui.Context, tab_item TabItem[T])) {
	tab.tabs.SetOnItemSelected(func(context *gui.Context, i int) {
		if i < 0 {
			return
		}
		
		tab_item := tab.Tab_Items[i]
		fn(context, tab_item)
	})
}

func (tab *Tab[T]) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	type segmented_control_item_value[T comparable] struct {
		widget gui.Widget
		value T
	}
	segmented_control_items := make([]widget.SegmentedControlItem[segmented_control_item_value[T]], len(tab.Tab_Items))
	for i := range tab.Tab_Items {
		tab_item := &tab.Tab_Items[i]
		segment_item := widget.SegmentedControlItem[segmented_control_item_value[T]]{
			Text: tab_item.Name,
			Value: segmented_control_item_value[T]{
				widget: tab_item.Widget,
				value: tab_item.Value,
			},
		}
		
		segmented_control_items[i] = segment_item
	}
	
	selected_item_index := tab.tabs.SelectedItemIndex() 
	if selected_item_index == -1 {
		selected_item_index = 1
		tab.tabs.SelectItemByIndex(selected_item_index)
	}
	adder.AddChild(&tab.tabs)
	
	selected_widget := tab.Tab_Items[selected_item_index].Widget
	tab.panel.SetContentConstraints(widget.PanelContentConstraintsNone)
	tab.panel.SetContent(selected_widget)
	adder.AddChild(&tab.panel)
	return nil
}

func (tab *Tab[T]) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Items: []gui.LinearLayoutItem{
			{
				Widget: &tab.tabs,
			},
			{
				Widget: &tab.panel,
				Size: gui.FlexibleSize(1),
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}