package CommonWidgets

import (
	"API-Client/basic"
	"image"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/guigui-gui/guigui/basicwidget/basicwidgetdraw"
	"github.com/hajimehoshi/ebiten/v2"
)

type TabItem[T any] struct {
	gui.DefaultWidget

	Text        string
	Size        gui.Size
	Value       T
	text_widget widget.Text
}

func (item *TabItem[T]) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	item.text_widget.SetValue(item.Text)
	item.text_widget.SetTabular(true)
	item.text_widget.SetVerticalAlign(widget.VerticalAlignMiddle)

	adder.AddChild(&item.text_widget)
	return nil
}

func (item *TabItem[T]) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	u := widget.UnitSize(ctx)
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Padding:   basic.NewPadding(u/4, u/2),
		Items: []gui.LinearLayoutItem{
			{
				Widget: &item.text_widget,
				Size:   item.Size,
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (item *TabItem[T]) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	point := item.text_widget.Measure(ctx, constraints)

	u := widget.UnitSize(ctx)
	point.X += u
	point.Y = u / 2
	return point
}

/*
func (tab_item *TabItem[T]) HandlePointingInput(ctx *gui.Context, widgetBounds *gui.WidgetBounds) gui.HandleInputResult {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		tab_item.tab.selected_index = tab_item.index
	}
	return gui.HandleInputResult{}
}
*/

func (tab_item *TabItem[T]) Draw(ctx *gui.Context, widgetBounds *gui.WidgetBounds, dst *ebiten.Image) {
	color_mod := ctx.ColorMode()
	background_color := basicwidgetdraw.ControlColor(color_mod, ctx.IsEnabled(tab_item))
	border_color := basicwidgetdraw.ControlSecondaryColor(color_mod, ctx.IsEnabled(tab_item))
	basicwidgetdraw.DrawRoundedRectBorder(ctx, dst, widgetBounds.Bounds(), background_color, border_color, 1, 1, basicwidgetdraw.RoundedRectBorderTypeRegular)
}

type tab[T any] struct {
	gui.DefaultWidget
	tab_items      []TabItem[T]
	selected_index int
}

func (tab *tab[T]) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	for i := range tab.tab_items {
		tab_item := &tab.tab_items[i]
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
	point.Y = widget.UnitSize(ctx)
	
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

func (tab *Tab[T]) OnSelect(fn func(text string, value T)) {

}

func (tab *Tab[T]) GetSelectedIndex() int {
	return 0
}

func (tab *Tab[T]) GetTabByIndex(index int) (text string, value T) {
	tab_item := &tab.tab.tab_items[index]
	return tab_item.Text, tab_item.Value
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
