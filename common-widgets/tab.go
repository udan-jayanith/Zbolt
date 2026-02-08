package CommonWidgets

import (
	"image"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/guigui-gui/guigui/basicwidget/basicwidgetdraw"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type TabItem[T any] struct {
	gui.DefaultWidget

	Text                     string
	Value                    T
	text_widget              widget.Text
	is_selected, is_hovering bool
}

func (item *TabItem[T]) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	item.text_widget.SetValue(item.Text)
	item.text_widget.SetTabular(true)
	adder.AddChild(&item.text_widget)
	return nil
}

func (item *TabItem[T]) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layouter.LayoutWidget(&item.text_widget, widgetBounds.Bounds())
}

func (tab_item *TabItem[T]) HandlePointingInput(ctx *gui.Context, widgetBounds *gui.WidgetBounds) gui.HandleInputResult {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		tab_item.is_selected = true
	}
	return gui.HandleInputResult{}
}

func (tab_item *TabItem[T]) Draw(ctx *gui.Context, widgetBounds *gui.WidgetBounds, dst *ebiten.Image) {
	color_mod := ctx.ColorMode()
	basicwidgetdraw.DrawRoundedRect(ctx, dst, widgetBounds.Bounds(), basicwidgetdraw.ControlColor(color_mod, false), 0)
}

type tab[T any] struct {
	gui.DefaultWidget
	tab_items []TabItem[T]
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
			Widget: tab_item,
		})
	}

	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (tab *tab[T]) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	u := widget.UnitSize(ctx)
	if w, ok := constraints.FixedWidth(); ok {
		return image.Pt(w, u)
	}
	return image.Pt(10, u)
}

type Tab[T any] struct {
	gui.DefaultWidget
	panel widget.Panel
	tab   tab[T]
}

func (tab *Tab[T]) OnSelect(fn func(text string, value T)) {

}

func (tab *Tab[T]) GetSelectedIndex() int {
	return 0
}

func (tab *Tab[T]) GetTab(index int) (text string, value T) {
	tab_item := &tab.tab.tab_items[index]
	return tab_item.Text, tab_item.Value
}

func (tab *Tab[T]) SetTabItems(tab_items []TabItem[T]) {
	tab.tab.tab_items = tab_items
}

func (tab *Tab[T]) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	tab.panel.SetContent(&tab.tab)
	tab.panel.SetContentConstraints(widget.PanelContentConstraintsFixedHeight)
	tab.panel.SetStyle(widget.PanelStyleSide)
	adder.AddChild(&tab.panel)
	return nil
}

func (tab *Tab[T]) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layouter.LayoutWidget(&tab.panel, widgetBounds.Bounds())
}
