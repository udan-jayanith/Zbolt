package CommonWidgets

import (
	"image"

	"API-Client/basic"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/guigui-gui/guigui/basicwidget/basicwidgetdraw"
	"github.com/hajimehoshi/ebiten/v2"
)

type Swap struct {
	From int
	To   int
}

type tab struct {
	gui.DefaultWidget
	tab_items []*tab_item

	holding_tab_item *tab_item
	closest          *tab_item
	swap             *Swap

	selected_item_index int
	on_select_fn        func(from, to TabItem, i, j int)
	on_swap             func(ctx *gui.Context, swap Swap)
	on_close            func(tab_item TabItem)
}

func (tab *tab) on_select(index int, tab_item TabItem) {
	// TODO: handle this
	/*
		if item.tabs_container.on_select_fn != nil {
			item.tabs_container.on_select_fn(item.tabs_container.tab_items[item.tabs_container.selected_item_index].tab_item, item.tab_item, item.tabs_container.selected_item_index, item.index)
		}
		item.tabs_container.selected_item_index = item.index
	*/
}

func (tab *tab) on_holding(index int, relative_cursor_x int) {
	// TODO:
	// item.tabs_container.holding_tab_item = item
}

func (tab *tab) on_holding_cancel(index int) {
	// TODO:
	/*
		  tab_widget := item.tabs_container

			if tab_widget.closest != nil {
				tab_widget.swap = &Swap{
					From: tab_widget.holding_tab_item.index,
					To:   tab_widget.closest.index,
				}
			}

			tab_widget.holding_tab_item = nil
			tab_widget.closest = nil
	*/
}

func (tab *tab) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	if tab.swap != nil {
		if tab.on_swap != nil {
			tab.on_swap(ctx, *tab.swap)
		}
		tab.swap = nil
	}

	for _, tab_item := range tab.tab_items {
		if tab.holding_tab_item != nil && tab.holding_tab_item.index == tab_item.index {
			continue
		}
		adder.AddWidget(tab_item)
	}

	if tab.holding_tab_item != nil {
		adder.AddWidget(tab.holding_tab_item)
	}
	return nil
}

func (tab *tab) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Items:     make([]gui.LinearLayoutItem, 0, len(tab.tab_items)),
	}

	b := widgetBounds.Bounds()

	if tab.holding_tab_item != nil {
		cursor_axis, _ := ebiten.CursorPosition()
		tab_item_bounds := image.Rectangle{
			Min: image.Point{
				X: cursor_axis - tab.holding_tab_item.relative_cursor_axis,
				Y: b.Min.Y,
			},
		}

		tab_item_bounds.Max.X = tab_item_bounds.Min.X + tab.holding_tab_item.Measure(ctx, gui.Constraints{}).X
		tab_item_bounds.Max.Y = b.Max.Y

		layouter.LayoutWidget(tab.holding_tab_item, tab_item_bounds)
	}

	for _, tab_item := range tab.tab_items {
		if tab.holding_tab_item != nil && tab.holding_tab_item.index == tab_item.index {
			w := tab.holding_tab_item.Measure(ctx, gui.Constraints{}).X
			layout.Items = append(layout.Items, gui.LinearLayoutItem{
				Size: gui.FixedSize(w),
			})
			continue
		}

		layout.Items = append(layout.Items, gui.LinearLayoutItem{
			Widget: tab_item,
		})
	}

	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (tab *tab) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	var point image.Point
	for _, tab_item := range tab.tab_items {
		measurement := tab_item.Measure(ctx, constraints)
		point.X += measurement.X
	}

	u := widget.UnitSize(ctx)
	point.Y = u

	return point
}

type Tab[T any] struct {
	gui.DefaultWidget
	panel widget.Panel
	tab   tab[T]
}

func (tab *tab) SetTabItems(tab_items []TabItem[T]) {
	if len(tab_items) != len(tab.tab.tab_items) {
		tab.tab.tab_items = make([]*tab_item[T], len(tab_items))
	}

	for i, item := range tab_items {
		tab.tab.tab_items[i] = &tab_item[T]{
			tab_item:   &item,
			tab_widget: &tab.tab,
			index:      i,
		}
	}
}

func (tab *tab) OnSwitch(fn func(from, to *TabItem[T])) {
	tab.tab.on_select_fn = fn
}

func (tab *tab) GetSelectedIndex() int {
	return tab.tab.selected_item_index
}

func (tab *tab) GetSelectedTab() (index int, value T) {
	if len(tab.tab.tab_items) == 0 {
		var t T
		return 0, t
	}
	return tab.tab.selected_item_index, tab.tab.tab_items[tab.tab.selected_item_index].tab_item.Value
}

func (tab *tab) GetTabByIndex(index int) (text string, value T) {
	tab_item := tab.tab.tab_items[index]
	return tab_item.tab_item.Text, tab_item.tab_item.Value
}

func (tab *tab) SelectTabItemByIndex(index int) {
	if tab.tab.on_select_fn != nil {
		tab.tab.on_select_fn(tab.tab.tab_items[tab.tab.selected_item_index].tab_item, tab.tab.tab_items[index].tab_item)
	}
	tab.tab.selected_item_index = index
}

func (tab *tab) OnSwap(fn func(ctx *gui.Context, swap Swap)) {
	tab.tab.on_swap = fn
}

func (tab *tab) OnClose(fn func(tab_item TabItem[T])) {
	tab.tab.on_close = fn
}

func (tab *tab) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	tab.panel.SetContent(&tab.tab)
	tab.panel.SetStyle(widget.PanelStyleSide)
	tab.panel.SetAutoBorder(true)
	tab.panel.SetContentConstraints(widget.PanelContentConstraintsFixedWidth)
	adder.AddWidget(&tab.panel)
	return nil
}

func (tab *tab) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Padding:   basic.NewPadding(2),
		Items: []gui.LinearLayoutItem{
			{
				Widget: &tab.panel,
				Size:   gui.FixedSize(tab.tab.Measure(ctx, gui.FixedHeightConstraints(1)).Y),
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (tab *tab) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	var point image.Point

	if w, ok := constraints.FixedWidth(); ok {
		point.X = w
	} else {
		point.X = widget.UnitSize(ctx)*4 + 4
	}

	if h, ok := constraints.FixedHeight(); ok {
		point.Y = h
	} else {
		point.Y = tab.tab.Measure(ctx, constraints).Y + 4
	}

	return point
}

func (tab *tab) Draw(ctx *gui.Context, widgetBounds *gui.WidgetBounds, dst *ebiten.Image) {
	cm := ctx.ColorMode()
	r := basic.BorderRadius(ctx)
	border_type := basicwidgetdraw.RoundedRectBorderTypeRegular

	background_color := basicwidgetdraw.BackgroundSecondaryColor(cm)
	basicwidgetdraw.DrawRoundedRect(ctx, dst, widgetBounds.Bounds(), background_color, r)

	clr1, clr2 := basicwidgetdraw.BorderColors(cm, border_type)
	basicwidgetdraw.DrawRoundedRectBorder(ctx, dst, widgetBounds.Bounds(), clr1, clr2, r, 2, border_type)
}
