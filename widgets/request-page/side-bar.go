package Requester

import (
	"API-Client/icons"
	"fmt"
	"image"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type SidebarItem[T comparable] struct {
	Text, IconName string
	Value          T
}

type sidebar_item_widget[T comparable] struct {
	gui.DefaultWidget

	icon_widget *icons.Icon

	text_widget    widget.Text
	sidebar_item   SidebarItem[T]
	sidebar_widget *Sidebar[T]
}

func (sd *sidebar_item_widget[T]) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	line_height := widget.LineHeight(ctx)

	if sd.icon_widget == nil && sd.sidebar_item.IconName == "" {
		sd.icon_widget = icons.NewIcon("request-page", line_height)
	} else if sd.icon_widget == nil {
		sd.icon_widget = icons.NewIcon(sd.sidebar_item.IconName, line_height)
	}
	adder.AddChild(sd.icon_widget)

	sd.text_widget.SetValue(sd.sidebar_item.Text)
	adder.AddChild(&sd.text_widget)

	return nil
}

func (sd *sidebar_item_widget[T]) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	u := widget.UnitSize(ctx)

	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Gap:       u / 6,
		Items: []gui.LinearLayoutItem{
			{},
			{
				Widget: sd.icon_widget,
			},
			{
				Widget: &sd.text_widget,
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (sd *sidebar_item_widget[T]) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	point := sd.text_widget.Measure(ctx, constraints)
	point.X += widget.UnitSize(ctx) / 4
	point.X += sd.icon_widget.Measure(ctx, constraints).X
	return point
}
func (sd *sidebar_item_widget[T]) HandlePointingInput(ctx *gui.Context, widgetBounds *gui.WidgetBounds) gui.HandleInputResult {
	if widgetBounds.IsHitAtCursor() && inpututil.IsMouseButtonJustPressed(ebiten.MouseButton2) {
		sd.sidebar_widget.right_clicked_item = sd
	}

	return gui.HandleInputResult{}
}

type Sidebar[T comparable] struct {
	gui.DefaultWidget

	options struct {
		search_widget widget.TextInput
		add_widget    widget.Button
	}
	list_widget       widget.List[T]
	list_widget_items []widget.ListItem[T]
	measurement       image.Point

	context_menu       widget.PopupMenu[struct{}]
	context_menu_pos   image.Point
	right_clicked_item *sidebar_item_widget[T]
	on_add_btn_clicked func(ctx *gui.Context)
}

func (sd *Sidebar[T]) SetItems(items []SidebarItem[T]) {
	sd.list_widget_items = make([]widget.ListItem[T], 0, len(items))
	for _, item := range items {
		content_widget := sidebar_item_widget[T]{
			sidebar_widget: sd,
			sidebar_item:   item,
		}
		sd.list_widget_items = append(sd.list_widget_items, widget.ListItem[T]{
			Content: &content_widget,
			Value:   item.Value,
		})
	}
}

func (sd *Sidebar[T]) OnAddButtonClicked(callback func(ctx *gui.Context)){
	sd.on_add_btn_clicked = callback
}

func (sd *Sidebar[T]) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	adder.AddChild(&sd.options.search_widget)

	sd.options.add_widget.SetIcon(icons.Store.Open("add"))
	if sd.on_add_btn_clicked != nil {
		sd.options.add_widget.SetOnDown(func(ctx *gui.Context) {
			sd.on_add_btn_clicked(ctx)
		})
	}
	adder.AddChild(&sd.options.add_widget)

	sd.list_widget.SetItems(sd.list_widget_items)
	adder.AddChild(&sd.list_widget)

	sd.context_menu.SetItemsByStrings([]string{"Rename", "Delete"})
	adder.AddChild(&sd.context_menu)
	return nil
}

func (sd *Sidebar[T]) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layouter.LayoutWidget(&sd.context_menu, image.Rectangle{
		Min: sd.context_menu_pos,
	})

	u := widget.UnitSize(ctx)
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Gap:       u / 4,
		Items: []gui.LinearLayoutItem{
			{
				Layout: gui.LinearLayout{
					Direction: gui.LayoutDirectionHorizontal,
					Gap:       u / 4,
					Items: []gui.LinearLayoutItem{
						{
							Widget: &sd.options.search_widget,
							Size:   gui.FlexibleSize(1),
						},
						{
							Widget: &sd.options.add_widget,
						},
					},
				},
			},
			{
				Widget: &sd.list_widget,
				Size:   gui.FlexibleSize(1),
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (sd *Sidebar[T]) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	return sd.list_widget.Measure(ctx, constraints)
}

func (sd *Sidebar[T]) HandlePointingInput(ctx *gui.Context, widgetBounds *gui.WidgetBounds) gui.HandleInputResult {
	if widgetBounds.IsHitAtCursor() && inpututil.IsMouseButtonJustReleased(ebiten.MouseButton2) && sd.right_clicked_item != nil {
		sd.context_menu_pos = image.Pt(ebiten.CursorPosition())
		sd.context_menu.SetOpen(true)

		fmt.Println("right clicked", sd.right_clicked_item.text_widget.Value())

		sd.right_clicked_item = nil
	}

	return gui.HandleInputResult{}
}
