package CommonWidgets

import (
	"image"

	"API-Client/basic"
	"API-Client/icons"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/hajimehoshi/ebiten/v2"
)

type TabItem[T any] struct {
	Text     string
	Value    T
	Closable bool
	Icon     *icons.Icon
}

type tab_item[T any] struct {
	gui.DefaultWidget
	text_widget  widget.Text
	close_widget struct {
		icon     *icons.Icon
		normal   *icons.Icon
		selected *icons.Icon
	}
	border_widget WidgetWithBorder[*tab_item[T]]

	tab_item   *TabItem[T]
	tab_widget *tab[T]
}

func (item *tab_item[T]) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	if item.tab_item.Icon != nil {
		adder.AddChild(item.tab_item.Icon)
	}

	text_widget := &item.text_widget
	text_widget.SetValue(item.tab_item.Text)
	text_widget.SetTabular(true)
	text_widget.SetVerticalAlign(widget.VerticalAlignMiddle)
	//text_widget.SetBold(item.tab_item)

	if true {
		text_widget.SetOpacity(1)
	} else {
		text_widget.SetOpacity(0.6)
	}
	adder.AddChild(&item.text_widget)

	if item.tab_item.Closable {
		if item.close_widget.icon == nil {
			line_height := widget.LineHeight(ctx)
			size := line_height - line_height/4
			item.close_widget.normal = icons.NewIcon("close", size)
			item.close_widget.selected = icons.NewIcon("close-grey", size)
		}

		if true {
			item.close_widget.icon = item.close_widget.normal
		} else {
			item.close_widget.icon = item.close_widget.selected
		}
		adder.AddChild(item.close_widget.icon)
	}
	return nil
}

func (item *tab_item[T]) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	u := widget.UnitSize(ctx)
	layout := gui.LinearLayout{
		Padding:   basic.NewPadding(u/4, u/2),
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
			Widget: item.close_widget.icon,
		})
	}

	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (item *tab_item[T]) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	point := item.text_widget.Measure(ctx, constraints)
	u := widget.UnitSize(ctx)
	padding := basic.NewPadding(u/4, u/2)

	gap := widget.UnitSize(ctx) / 4
	if item.tab_item.Icon != nil {
		icon_measurement := item.tab_item.Icon.Measure(ctx, constraints)
		point.X += icon_measurement.X + gap
	}

	if item.tab_item.Closable {
		icon_measurement := item.close_widget.icon.Measure(ctx, constraints)
		point.X += icon_measurement.X + gap
	}
	point.X += padding.End + padding.Start
	point.Y = u + u/4

	return point
}

func (item *tab_item[T]) Draw(ctx *gui.Context, widgetBounds *gui.WidgetBounds, dst *ebiten.Image) {
	item.border_widget.Draw(ctx, widgetBounds, dst)
}

type tab[T any] struct {
	gui.DefaultWidget
	tab_items []*tab_item[T]

	selected_item_index int
	on_select_fn        func(tab_item *TabItem[T], index int)
}

func (tab *tab[T]) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	for _, tab_item := range tab.tab_items {
		adder.AddChild(tab_item)
	}
	return nil
}

func (tab *tab[T]) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Items:     make([]gui.LinearLayoutItem, 0, len(tab.tab_items)),
	}

	for _, tab_item := range tab.tab_items {
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
	point.Y = u + u/4

	return point
}

type Tab[T any] struct {
	gui.DefaultWidget
	panel widget.Panel
	tab   tab[T]
}

func (tab *Tab[T]) SetTabItems(tab_items []TabItem[T]) {
	tab.tab.tab_items = make([]*tab_item[T], 0, len(tab_items))
	for _, item := range tab_items {
		tab.tab.tab_items = append(tab.tab.tab_items, &tab_item[T]{
			tab_item:   &item,
			tab_widget: &tab.tab,
		})
	}
}

func (tab *Tab[T]) OnSelect(fn func(tab_item *TabItem[T], index int)) {
	tab.tab.on_select_fn = fn
}

func (tab *Tab[T]) GetSelectedIndex() int {
	return tab.tab.selected_item_index
}

func (tab *Tab[T]) GetSelectedTab() (index int, value T) {
	return tab.tab.selected_item_index, tab.tab.tab_items[tab.tab.selected_item_index].tab_item.Value
}

func (tab *Tab[T]) GetTabByIndex(index int) (text string, value T) {
	tab_item := tab.tab.tab_items[index]
	return tab_item.tab_item.Text, tab_item.tab_item.Value
}

func (tab *Tab[T]) SelectTabItemByIndex(index int) {
	tab.tab.selected_item_index = index
	tab.tab.on_select_fn(tab.tab.tab_items[index].tab_item, index)
}

func (tab *Tab[T]) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	tab.panel.SetContent(&tab.tab)
	tab.panel.SetStyle(widget.PanelStyleSide)
	tab.panel.SetContentConstraints(widget.PanelContentConstraintsFixedWidth)
	adder.AddChild(&tab.panel)
	return nil
}

func (tab *Tab[T]) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
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
		point.X = widget.UnitSize(ctx) * 4
	}

	if h, ok := constraints.FixedHeight(); ok {
		point.Y = h
	} else {
		point.Y = tab.tab.Measure(ctx, constraints).Y
	}

	return point
}
