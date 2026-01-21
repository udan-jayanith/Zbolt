package CommonWidgets

import (
	"errors"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type TabItem[T comparable] struct {
	Name   string
	Widget gui.Widget
	Value  T
}

type Tab[T comparable] struct {
	gui.DefaultWidget

	tabs      widget.SegmentedControl[string]
	Tab_Items []TabItem[T]
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

func (tab *Tab[T]) GetSelectedWidget() gui.Widget{
	selected_item_index := tab.tabs.SelectedItemIndex()
	if selected_item_index == -1 {
		selected_item_index = 0
		tab.tabs.SelectItemByIndex(selected_item_index)
	}

	return tab.Tab_Items[selected_item_index].Widget
}

func (tab *Tab[T]) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	if len(tab.Tab_Items) == 0 {
		return errors.New("No selected items found")
	}
	
	segmented_control_items := make([]widget.SegmentedControlItem[string], len(tab.Tab_Items))
	for i := range tab.Tab_Items {
		tab_item := &tab.Tab_Items[i]
		segment_item := widget.SegmentedControlItem[string]{
			Text: tab_item.Name,
		}

		segmented_control_items[i] = segment_item
	}
	tab.tabs.SetItems(segmented_control_items)
	adder.AddChild(&tab.tabs)
	
	selected_widget :=	tab.GetSelectedWidget()
	adder.AddChild(selected_widget)
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
				Widget: tab.GetSelectedWidget(),
				Size:   gui.FlexibleSize(1),
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}