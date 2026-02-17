package CommonWidgets

import (
	"image"

	"API-Client/basic"
	"API-Client/icons"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type tab_item[T any] struct {
	gui.DefaultWidget
	text_widget  widget.Text
	close_widget struct {
		icon *icons.Icon
		normal  *icons.Icon
		selected *icons.Icon
	}

	tab_item *TabItem[T]
}

func (item *tab_item[T]) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	if item.tab_item.Icon != nil {
		adder.AddChild(item.tab_item.Icon)
	}

	text_widget := &item.text_widget
	text_widget.SetValue(item.tab_item.Text)
	text_widget.SetTabular(true)
	text_widget.SetVerticalAlign(widget.VerticalAlignMiddle)
	text_widget.SetBold(item.tab_item.is_hovering)

	if item.tab_item.tab.selected_index == item.tab_item.index {
		text_widget.SetOpacity(1)
	} else {
		text_widget.SetOpacity(0.6)
	}
	adder.AddChild(&item.text_widget)

	if item.tab_item.Closable {
		if item.close_widget.icon == nil {
			line_height := widget.LineHeight(ctx)
			size :=  line_height-line_height/4
			item.close_widget.normal = icons.NewIcon("close", size)
			item.close_widget.selected = icons.NewIcon("close-grey", size)
		}
	
		if item.tab_item.tab.selected_index == item.tab_item.index {
			item.close_widget.icon = item.close_widget.normal 
		} else {
			item.close_widget.icon = item.close_widget.selected
		}
		adder.AddChild(item.close_widget.icon)
	}
	return nil
}

func (item *tab_item[T]) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layout := gui.LinearLayout{
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

	gap := widget.UnitSize(ctx) / 4
	if item.tab_item.Icon != nil {
		icon_measurement := item.tab_item.Icon.Measure(ctx, constraints)
		point.X += icon_measurement.X + gap
	}

	if item.tab_item.Closable {
		icon_measurement := item.close_widget.icon.Measure(ctx, constraints)
		point.X += icon_measurement.X + gap
	}

	return point
}

type TabItem[T any] struct {
	gui.DefaultWidget

	Text     string
	Size     gui.Size
	Value    T
	Closable bool
	Icon     *icons.Icon

	border_widget WidgetWithBorder[*gui.WidgetWithPadding[*tab_item[T]]]

	index       int
	tab         *tab[T]
	is_hovering bool
}

func (item *TabItem[T]) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	u := widget.UnitSize(ctx)

	padding := gui.WidgetWithPadding[*tab_item[T]]{}
	padding.SetPadding(basic.NewPadding(u/4, u/2))

	padding.Widget().tab_item = item

	item.border_widget.SetWidget(&padding)
	adder.AddChild(&item.border_widget)
	return nil
}

func (item *TabItem[T]) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Items: []gui.LinearLayoutItem{
			{
				Widget: &item.border_widget,
				Size:   item.Size,
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (item *TabItem[T]) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	return item.border_widget.Measure(ctx, constraints)
}

func (tab_item *TabItem[T]) HandlePointingInput(ctx *gui.Context, widgetBounds *gui.WidgetBounds) gui.HandleInputResult {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) && widgetBounds.IsHitAtCursor() {
		tab_item.tab.selected_index = tab_item.index
		if tab_item.tab.on_select_fn != nil {
			tab_item.tab.on_select_fn(tab_item, tab_item.index)
		}
	}
	tab_item.is_hovering = widgetBounds.IsHitAtCursor()

	return gui.HandleInputResult{}
}

type tab[T any] struct {
	gui.DefaultWidget
	tab_items      []TabItem[T]
	selected_index int
	on_select_fn   func(tab_item *TabItem[T], index int)
}

func (tab *tab[T]) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	for i := range tab.tab_items {
		tab_item := &tab.tab_items[i]
		if tab_item.tab == nil {
			tab_item.tab = tab
			tab_item.index = i
		}
		adder.AddChild(tab_item)
	}
	return nil
}

func (tab *tab[T]) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Items:     make([]gui.LinearLayoutItem, 0, len(tab.tab_items)),
	}

	for i := range tab.tab_items {
		tab_item := &tab.tab_items[i]
		layout.Items = append(layout.Items, gui.LinearLayoutItem{
			Size:   tab_item.Size,
			Widget: tab_item,
		})
	}

	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (tab *tab[T]) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	var point image.Point
	for i := range tab.tab_items {
		tab_item := &tab.tab_items[i]
		mesurement := tab_item.Measure(ctx, constraints)
		point.X += mesurement.X
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
	tab.tab.tab_items = tab_items
}

func (tab *Tab[T]) OnSelect(fn func(tab_item *TabItem[T], index int)) {
	tab.tab.on_select_fn = fn
}

func (tab *Tab[T]) GetSelectedIndex() int {
	return tab.tab.selected_index
}

func (tab *Tab[T]) GetSelectedTab() (index int, value T) {
	return tab.tab.selected_index, tab.tab.tab_items[tab.tab.selected_index].Value
}

func (tab *Tab[T]) GetTabByIndex(index int) (text string, value T) {
	tab_item := &tab.tab.tab_items[index]
	return tab_item.Text, tab_item.Value
}

func (tab *Tab[T]) SelectTabItemByIndex(index int) {
	tab.tab.selected_index = index
	tab.tab.on_select_fn(&tab.tab.tab_items[index], index)
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
	return tab.tab.Measure(ctx, constraints)
}
