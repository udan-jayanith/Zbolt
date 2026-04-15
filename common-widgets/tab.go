package CommonWidgets

import (
	"image"
	"image/color"

	"API-Client/basic"
	"API-Client/icons"

	draw_color "API-Client/common-widgets/internal/draw"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/guigui-gui/guigui/basicwidget/basicwidgetdraw"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type TabItem[T any] struct {
	Text     string
	Value    T
	Closable bool
	Icon     *icons.Icon
}

type tab_item[T any] struct {
	gui.DefaultWidget
	text_widget widget.Text
	close_icon  icons.Icon

	tab_item   *TabItem[T]
	tab_widget *tab[T]

	relative_cursor_axis int
	bounds               image.Rectangle

	is_hovering bool
	index       int
}

func (item *tab_item[T]) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	if item.tab_item.Icon != nil {
		adder.AddWidget(item.tab_item.Icon)
	}

	text_widget := &item.text_widget
	text_widget.SetValue(item.tab_item.Text)
	text_widget.SetTabular(true)
	text_widget.SetVerticalAlign(widget.VerticalAlignMiddle)

	adder.AddWidget(&item.text_widget)

	if item.tab_item.Closable {
		if item.close_icon.IconName() == "" {
			line_height := widget.LineHeight(ctx)
			size := line_height - line_height/4
			item.close_icon.SetSize(size)
		}

		if item.tab_widget.selected_item_index == item.index {
			item.close_icon.SetIcon("close")
		} else {
			item.close_icon.SetIcon("close-grey")
		}

		if item.tab_widget.on_close != nil {
			item.close_icon.OnClick(func() {
				item.tab_widget.on_close(*item.tab_item)
			})
		}
		adder.AddWidget(&item.close_icon)
	}

	return nil
}

func (item *tab_item[T]) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layout := gui.LinearLayout{
		Padding:   basic.NewPadding(0, widget.LineHeight(ctx)/2),
		Gap:       widget.UnitSize(ctx) / 4,
		Direction: gui.LayoutDirectionHorizontal,
		Items:     make([]gui.LinearLayoutItem, 0, 3),
	}

	if item.tab_item.Icon != nil {
		layout.Items = append(layout.Items, gui.LinearLayoutItem{
			Widget: item.tab_item.Icon,
		})
	}

	layout.Items = append(layout.Items, gui.LinearLayoutItem{
		Widget: &item.text_widget,
		Size:   gui.FlexibleSize(1),
	})

	if item.tab_item.Closable {
		layout.Items = append(layout.Items, gui.LinearLayoutItem{
			Widget: &item.close_icon,
		})
	}

	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (item *tab_item[T]) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	point := item.text_widget.Measure(ctx, constraints)
	padding := basic.NewPadding(0, widget.LineHeight(ctx)/2)

	gap := widget.UnitSize(ctx) / 4
	if item.tab_item.Icon != nil {
		icon_measurement := item.tab_item.Icon.Measure(ctx, constraints)
		point.X += icon_measurement.X + gap
	}

	if item.tab_item.Closable {
		icon_measurement := item.close_icon.Measure(ctx, constraints)
		point.X += icon_measurement.X + gap
	}
	point.X += padding.End + padding.Start
	point.Y = widget.UnitSize(ctx)
	return point
}

func (item *tab_item[T]) Draw(ctx *gui.Context, widgetBounds *gui.WidgetBounds, dst *ebiten.Image) {
	var background_color color.Color
	var border_type basicwidgetdraw.RoundedRectBorderType

	cm := ctx.ColorMode()
	if item.tab_widget.selected_item_index == item.index {
		background_color = draw_color.Color2(cm, draw_color.ColorTypeBase, 0.2, 0.2)
		border_type = basicwidgetdraw.RoundedRectBorderTypeInset
	} else if item.tab_widget.closest != nil && item.tab_widget.closest.index == item.index {
		background_color = draw_color.Color2(cm, draw_color.ColorTypeBase, 0.4, 0.4)
		border_type = basicwidgetdraw.RoundedRectBorderTypeRegular
	} else if widgetBounds.IsHitAtCursor() {
		background_color = draw_color.Color2(cm, draw_color.ColorTypeBase, 0.2, 0.2)
		border_type = basicwidgetdraw.RoundedRectBorderTypeRegular
	} else {
		background_color = basicwidgetdraw.BackgroundSecondaryColor(ctx.ColorMode())
		border_type = basicwidgetdraw.RoundedRectBorderTypeRegular
	}

	r := basic.BorderRadius(ctx)
	basicwidgetdraw.DrawRoundedRect(ctx, dst, widgetBounds.Bounds(), background_color, r)

	clr1, clr2 := basicwidgetdraw.BorderColors(cm, border_type)
	basicwidgetdraw.DrawRoundedRectBorder(ctx, dst, widgetBounds.Bounds(), clr1, clr2, r, 1, border_type)
}

func (item *tab_item[T]) HandlePointingInput(ctx *gui.Context, widgetBounds *gui.WidgetBounds) gui.HandleInputResult {
	item.bounds = widgetBounds.Bounds()
	item.is_hovering = widgetBounds.IsHitAtCursor()

	if item.is_hovering && inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		if item.tab_widget.on_select_fn != nil {
			item.tab_widget.on_select_fn(item.tab_widget.tab_items[item.tab_widget.selected_item_index].tab_item, item.tab_item)
		}
		item.tab_widget.selected_item_index = item.index
	} else if item.is_hovering && ebiten.IsMouseButtonPressed(ebiten.MouseButton0) && item.tab_widget.holding_tab_item == nil {
		item.tab_widget.holding_tab_item = item
		cursor_axis, _ := ebiten.CursorPosition()
		item.relative_cursor_axis = cursor_axis - widgetBounds.Bounds().Min.X
	} else if inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) && item.tab_widget.holding_tab_item != nil {
		tab_widget := item.tab_widget

		if tab_widget.closest != nil {
			tab_widget.swap = &Swap{
				From: tab_widget.holding_tab_item.index,
				To:   tab_widget.closest.index,
			}
		}

		tab_widget.holding_tab_item = nil
		tab_widget.closest = nil
	}

	is_overlapping := item.is_overlapping()
	if is_overlapping {
		if item.tab_widget.closest == nil {
			item.tab_widget.closest = item
		}

		dis1, dis2 := item.distance(), item.tab_widget.closest.distance()
		if dis1 < dis2 {
			item.tab_widget.closest = item
		}
	} else if !is_overlapping && item.tab_widget.closest != nil && item.tab_widget.closest.index == item.index {
		item.tab_widget.closest = nil
	}

	return gui.HandleInputResult{}
}

func (item *tab_item[T]) is_overlapping() bool {
	return item.tab_widget.holding_tab_item != nil &&
		item.tab_widget.holding_tab_item.index != item.index &&
		!item.bounds.Intersect(item.tab_widget.holding_tab_item.bounds).Empty()
}

func (item *tab_item[T]) distance() int {
	min_x, max_x := item.bounds.Min.X, item.bounds.Max.X
	x, _ := ebiten.CursorPosition()

	var min_dis, max_dis int
	min_dis = max(x, min_x) - min(x, min_x)
	max_dis = max(x, max_x) - min(x, max_x)
	return min(max_dis, min_dis)
}

type Swap struct {
	From int
	To   int
}

type tab[T any] struct {
	gui.DefaultWidget
	tab_items []*tab_item[T]

	holding_tab_item *tab_item[T]
	closest          *tab_item[T]
	swap             *Swap

	selected_item_index int
	on_select_fn        func(from, to *TabItem[T])
	on_swap             func(ctx *gui.Context, swap Swap)
	on_close            func(tab_item TabItem[T])
}

func (tab *tab[T]) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
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

func (tab *tab[T]) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
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

func (tab *tab[T]) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
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

func (tab *Tab[T]) SetTabItems(tab_items []TabItem[T]) {
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

func (tab *Tab[T]) OnSwitch(fn func(from, to *TabItem[T])) {
	tab.tab.on_select_fn = fn
}

func (tab *Tab[T]) GetSelectedIndex() int {
	return tab.tab.selected_item_index
}

func (tab *Tab[T]) GetSelectedTab() (index int, value T) {
	if len(tab.tab.tab_items) == 0 {
		var t T
		return 0, t
	}
	return tab.tab.selected_item_index, tab.tab.tab_items[tab.tab.selected_item_index].tab_item.Value
}

func (tab *Tab[T]) GetTabByIndex(index int) (text string, value T) {
	tab_item := tab.tab.tab_items[index]
	return tab_item.tab_item.Text, tab_item.tab_item.Value
}

func (tab *Tab[T]) SelectTabItemByIndex(index int) {
	if tab.tab.on_select_fn != nil {
		tab.tab.on_select_fn(tab.tab.tab_items[tab.tab.selected_item_index].tab_item, tab.tab.tab_items[index].tab_item)
	}
	tab.tab.selected_item_index = index
}

func (tab *Tab[T]) OnSwap(fn func(ctx *gui.Context, swap Swap)) {
	tab.tab.on_swap = fn
}

func (tab *Tab[T]) OnClose(fn func(tab_item TabItem[T])) {
	tab.tab.on_close = fn
}

func (tab *Tab[T]) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	tab.panel.SetContent(&tab.tab)
	tab.panel.SetStyle(widget.PanelStyleSide)
	tab.panel.SetAutoBorder(true)
	tab.panel.SetContentConstraints(widget.PanelContentConstraintsFixedWidth)
	adder.AddWidget(&tab.panel)
	return nil
}

func (tab *Tab[T]) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
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

func (tab *Tab[T]) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
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

func (tab *Tab[T]) Draw(ctx *gui.Context, widgetBounds *gui.WidgetBounds, dst *ebiten.Image) {
	cm := ctx.ColorMode()
	r := basic.BorderRadius(ctx)
	border_type := basicwidgetdraw.RoundedRectBorderTypeRegular

	background_color := basicwidgetdraw.BackgroundSecondaryColor(cm)
	basicwidgetdraw.DrawRoundedRect(ctx, dst, widgetBounds.Bounds(), background_color, r)

	clr1, clr2 := basicwidgetdraw.BorderColors(cm, border_type)
	basicwidgetdraw.DrawRoundedRectBorder(ctx, dst, widgetBounds.Bounds(), clr1, clr2, r, 2, border_type)
}
