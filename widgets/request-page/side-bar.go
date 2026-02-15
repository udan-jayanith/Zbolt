package Requester

import (
	"API-Client/icons"
	"image"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type SidebarItem[T comparable] struct {
	Text, IconName string
	Value          T
}

type sidebar_item_widget[T comparable] struct {
	gui.DefaultWidget

	icon_widget  *icons.Icon
	text_widget  widget.Text
	sidebar_item SidebarItem[T]
}

func (sd *sidebar_item_widget[T]) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	line_height := widget.LineHeight(ctx)

	if sd.icon_widget == nil && sd.sidebar_item.IconName == "" {
		sd.icon_widget = icons.NewIcon("circle", line_height)
	} else if sd.icon_widget == nil {
		sd.icon_widget = icons.NewIcon(sd.sidebar_item.IconName, line_height)
	}
	adder.AddChild(sd.icon_widget)

	sd.text_widget.SetValue(sd.sidebar_item.Text)
	adder.AddChild(&sd.text_widget)
	return nil
}

func (sd *sidebar_item_widget[T]) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Gap:       widget.UnitSize(ctx) / 4,
		Items: []gui.LinearLayoutItem{
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

type Sidebar[T comparable] struct {
	gui.DefaultWidget

	options struct {
		search_widget widget.TextInput
		add_widget    widget.Button
	}
	list_widget       widget.List[T]
	list_widget_items []SidebarItem[T]
	measurement       image.Point
}

func (sd *Sidebar[T]) SetItems(items []SidebarItem[T]){
	sd.list_widget_items = items
}

func (sd *Sidebar[T]) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	adder.AddChild(&sd.options.search_widget)

	sd.options.add_widget.SetIcon(icons.Store.Open("add"))
	adder.AddChild(&sd.options.add_widget)

	items := make([]widget.ListItem[T], 0, len(sd.list_widget_items))
	for _, item := range sd.list_widget_items {
		content_widget := sidebar_item_widget[T]{
			sidebar_item: item,
		}
		items = append(items, widget.ListItem[T]{
			Content: &content_widget,
			Value:   item.Value,
		})
	}
	sd.list_widget.SetItems(items)
	adder.AddChild(&sd.list_widget)

	return nil
}

func (sd *Sidebar[T]) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	u := widget.UnitSize(ctx)
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Items: []gui.LinearLayoutItem{
			{
				Size: gui.FixedSize(u / 4),
			},
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
				Size: gui.FixedSize(u / 4),
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
