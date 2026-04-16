package CommonWidgets

import (
	"image"

	"API-Client/basic"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/guigui-gui/guigui/basicwidget/basicwidgetdraw"
	"github.com/hajimehoshi/ebiten/v2"
)

type TabItemContainer struct {
	Index int
	Item  TabItem
}

type tabs_container struct {
	gui.DefaultWidget
	tab_items []*tab_item

	holding struct {
		is_holding               bool
		tab_item_index           int
		relative_cursor_position int
		closest_index            int // closest_index is the closest tab item that is closest to the holding tab item
	}
	selected_item_index int

	listeners struct {
		on_select func(from TabItemContainer, to TabItemContainer)
		on_swap   func(from TabItemContainer, to TabItemContainer)
		on_close  func(closed TabItemContainer)
	}
}

func (tab *tabs_container) on_select(index int, tab_item TabItem) {
	if tab.listeners.on_select != nil {
		from := TabItemContainer{
			Index: tab.selected_item_index,
			Item:  tab.tab_items[tab.selected_item_index].tab_item,
		}
		to := TabItemContainer{
			Index: index,
			Item:  tab_item,
		}
		tab.listeners.on_select(from, to)
	}
	tab.selected_item_index = index
}

func (tab *tabs_container) on_holding(index int, relative_cursor_x int) {
	tab.holding.is_holding = true
	tab.holding.tab_item_index = index
	tab.holding.relative_cursor_position = relative_cursor_x
}

func (tab *tabs_container) on_holding_cancel(index int) {
	tab.holding.is_holding = false
	if tab.holding.closest_index != -1 && len(tab.tab_items) > 0 && tab.listeners.on_swap != nil {
		from := TabItemContainer{
			Index: index,
			Item:  tab.tab_items[index].tab_item,
		}
		to := TabItemContainer{
			Index: tab.holding.closest_index,
			Item:  tab.tab_items[tab.holding.closest_index].tab_item,
		}
		tab.listeners.on_swap(from, to)
	}
}

func (tab *tabs_container) on_close(index int, item TabItem) {
	if tab.listeners.on_close == nil {
		return
	}
	tab.listeners.on_close(TabItemContainer{
		Index: index,
		Item:  item,
	})
}

func (tab *tabs_container) update_tab_items(tab_items []TabItem) {
	for i, widget := range tab.tab_items {
		widget.tab_item = tab_items[i]
	}
}

func (tab *tabs_container) set_tab_items(tab_items []TabItem) {
	if len(tab.tab_items) == len(tab_items) {
		tab.update_tab_items(tab_items)
		return
	}

	tab.tab_items = make([]*tab_item, 0, len(tab_items))
	for i, item := range tab_items {
		tab_item_widget := tab_item{}
		tab_item_widget.index = i
		tab_item_widget.tabs_container = tab
		tab_item_widget.tab_item = item
		tab.tab_items = append(tab.tab_items, &tab_item_widget)
	}

	if len(tab_items) > 0 {
		tab.on_select(0, tab_items[0])
	}
}

func (tab *tabs_container) select_tab(index int) {
	tab.on_select(index, tab.tab_items[index].tab_item)
}

func (tab *tabs_container) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	for _, tab_item := range tab.tab_items {
		if tab.holding.is_holding && tab.holding.tab_item_index == tab_item.index {
			continue
		}
		adder.AddWidget(tab_item)
	}

	if tab.holding.is_holding {
		adder.AddWidget(tab.tab_items[tab.holding.tab_item_index])
	}
	return nil
}

func (tab *tabs_container) holding_item_bounds(ctx *gui.Context, b image.Rectangle) image.Rectangle {
	cursor_axis, _ := ebiten.CursorPosition()
	tab_item_bounds := image.Rectangle{
		Min: image.Point{
			X: cursor_axis - tab.holding.relative_cursor_position,
			Y: b.Min.Y,
		},
	}

	holding_tab_item := tab.tab_items[tab.holding.tab_item_index]
	tab_item_bounds.Max.X = tab_item_bounds.Min.X + holding_tab_item.Measure(ctx, gui.Constraints{}).X
	tab_item_bounds.Max.Y = b.Max.Y
	return tab_item_bounds
}

func (tab *tabs_container) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Items:     make([]gui.LinearLayoutItem, 0, len(tab.tab_items)),
	}

	if tab.holding.is_holding {
		b := widgetBounds.Bounds()
		layouter.LayoutWidget(tab.tab_items[tab.holding.tab_item_index], tab.holding_item_bounds(ctx, b))
	}

	for _, tab_item := range tab.tab_items {
		if tab.holding.is_holding && tab.holding.tab_item_index == tab_item.index {
			holding_tab_item := tab.tab_items[tab.holding.tab_item_index]
			w := holding_tab_item.Measure(ctx, gui.Constraints{}).X
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

func (tab *tabs_container) HandlePointingInput(ctx *gui.Context, widgetBounds *gui.WidgetBounds) gui.HandleInputResult {
	if !widgetBounds.IsHitAtCursor() || !tab.holding.is_holding {
		return gui.HandleInputResult{}
	}

	b := widgetBounds.Bounds()
	holding_item_bounds := tab.holding_item_bounds(ctx, b)

	var closest_overlaying_item_bounds image.Rectangle
	closest_overlaying_item_index := -1

	for _, tab_item := range tab.tab_items {
		if tab_item.index == tab.holding.tab_item_index {
			continue
		}

		w := tab_item.Measure(ctx, gui.Constraints{}).X
		bounds := image.Rect(b.Min.X, b.Min.Y, b.Min.X+w, b.Max.Y)
		b.Min.X += w

		intersection_w := holding_item_bounds.Intersect(bounds).Dx()
		if intersection_w > 0 && closest_overlaying_item_index == -1 || intersection_w > closest_overlaying_item_bounds.Intersect(holding_item_bounds).Dx() {
			closest_overlaying_item_bounds = bounds
			closest_overlaying_item_index = tab_item.index
		}
	}

	tab.holding.closest_index = closest_overlaying_item_index
	return gui.HandleInputResult{}
}

func (tab *tabs_container) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	var point image.Point
	for _, tab_item := range tab.tab_items {
		measurement := tab_item.Measure(ctx, constraints)
		point.X += measurement.X
	}

	u := widget.UnitSize(ctx)
	point.Y = u

	return point
}

type Tab struct {
	gui.DefaultWidget
	panel         widget.Panel
	tab_container tabs_container
}

func (tab *Tab) SetTabItem(items []TabItem) {
	tab.tab_container.set_tab_items(items)
}

func (tab *Tab) SelectTab(index int) {
	tab.tab_container.select_tab(index)
}

func (tab *Tab) SelectedTab() int {
	return tab.tab_container.selected_item_index
}

func (tab *Tab) OnSelect(fn func(from TabItemContainer, to TabItemContainer)) {
	tab.tab_container.listeners.on_select = fn
}

func (tab *Tab) OnSwap(fn func(from TabItemContainer, to TabItemContainer)) {
	tab.tab_container.listeners.on_swap = fn
}

func (tab *Tab) OnClose(fn func(closed TabItemContainer)) {
	tab.tab_container.listeners.on_close = fn
}

func (tab *Tab) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	tab.panel.SetContent(&tab.tab_container)
	tab.panel.SetStyle(widget.PanelStyleSide)
	tab.panel.SetAutoBorder(true)
	tab.panel.SetContentConstraints(widget.PanelContentConstraintsFixedWidth)
	adder.AddWidget(&tab.panel)
	return nil
}

func (tab *Tab) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Padding:   basic.NewPadding(2),
		Items: []gui.LinearLayoutItem{
			{
				Widget: &tab.panel,
				Size:   gui.FixedSize(tab.tab_container.Measure(ctx, gui.Constraints{}).Y),
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (tab *Tab) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	var point image.Point

	if w, ok := constraints.FixedWidth(); ok {
		point.X = w
	} else {
		point.X = widget.UnitSize(ctx)*4 + 4
	}

	if h, ok := constraints.FixedHeight(); ok {
		point.Y = h
	} else {
		point.Y = tab.tab_container.Measure(ctx, constraints).Y + 4
	}

	return point
}

func (tab *Tab) Draw(ctx *gui.Context, widgetBounds *gui.WidgetBounds, dst *ebiten.Image) {
	cm := ctx.ColorMode()
	r := basic.BorderRadius(ctx)
	border_type := basicwidgetdraw.RoundedRectBorderTypeRegular

	background_color := basicwidgetdraw.BackgroundSecondaryColor(cm)
	basicwidgetdraw.DrawRoundedRect(ctx, dst, widgetBounds.Bounds(), background_color, r)

	clr1, clr2 := basicwidgetdraw.BorderColors(cm, border_type)
	basicwidgetdraw.DrawRoundedRectBorder(ctx, dst, widgetBounds.Bounds(), clr1, clr2, r, 2, border_type)
}
