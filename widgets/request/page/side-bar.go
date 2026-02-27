package request_page

import (
	CommonWidgets "API-Client/common-widgets"
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
		sd.sidebar_widget.list.contextmenu.right_clicked_item = sd
	}

	return gui.HandleInputResult{}
}

type Sidebar[T comparable] struct {
	gui.DefaultWidget

	options struct {
		search_widget widget.TextInput
		add           struct {
			widget         widget.Button
			add_button_pos image.Point
			menu           widget.PopupMenu[struct{}]

			on_request_clicked func(ctx *gui.Context)
			folder_popup       CommonWidgets.SimpleFormPopup
		}
	}

	list struct {
		widget widget.List[T]
		items  []widget.ListItem[T]

		contextmenu struct {
			menu     widget.PopupMenu[struct{}]
			position image.Point

			rename_popup_widget CommonWidgets.SimpleFormPopup
			right_clicked_item  *sidebar_item_widget[T]
		}
	}

	on_item_clicked func(value T)
	on_items_moved  func(context *gui.Context, from int, count int, to int)
}

func (sd *Sidebar[T]) SetItems(items []SidebarItem[T]) {
	sd.list.items = make([]widget.ListItem[T], 0, len(items))
	for _, item := range items {
		content_widget := sidebar_item_widget[T]{
			sidebar_widget: sd,
			sidebar_item:   item,
		}
		sd.list.items = append(sd.list.items, widget.ListItem[T]{
			Content: &content_widget,
			Value:   item.Value,
			Movable: false,
		})
	}
}

func (sd *Sidebar[T]) OneItemsMoved(f func(context *gui.Context, from int, count int, to int)) {
	sd.on_items_moved = f
}

func (sd *Sidebar[T]) OnAddButtonClicked(callback func(ctx *gui.Context)) {
	sd.options.add.on_request_clicked = callback
}

func (sd *Sidebar[T]) OnItemClicked(callback func(value T)) {
	sd.on_item_clicked = callback
}

func (sd *Sidebar[T]) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	adder.AddChild(&sd.options.search_widget)

	if sd.on_items_moved != nil {
		sd.list.widget.SetOnItemsMoved(sd.on_items_moved)
	}

	sd.options.add.widget.SetIcon(icons.Store.Open("add"))
	sd.options.add.menu.SetItemsByStrings([]string{"Request", "Folder"})
	sd.options.add.widget.SetOnDown(func(ctx *gui.Context) {
		sd.options.add.add_button_pos = image.Pt(ebiten.CursorPosition())
		sd.options.add.menu.SetOpen(true)
	})

	sd.options.add.folder_popup.SetButtonText("Create")
	sd.options.add.folder_popup.SetFieldValue("Enter folder name")

	sd.options.add.menu.SetOnItemSelected(func(context *gui.Context, index int) {
		if sd.options.add.on_request_clicked != nil && index == 0 {
			sd.options.add.on_request_clicked(ctx)
		} else if index == 1 {
			sd.options.add.folder_popup.SetOpen(true)
		}
	})
	adder.AddChild(&sd.options.add.widget)

	sd.list.widget.SetItems(sd.list.items)
	sd.list.widget.SetOnItemSelected(func(context *gui.Context, index int) {
		if sd.on_item_clicked != nil {
			sd.on_item_clicked(sd.list.items[index].Value)
		}
	})
	adder.AddChild(&sd.list.widget)

	sd.list.contextmenu.menu.SetItemsByStrings([]string{"Rename", "Delete"})
	sd.list.contextmenu.menu.SetOnItemSelected(func(context *gui.Context, index int) {
		if index == 0 {
			sd.list.contextmenu.rename_popup_widget.SetOpen(true)
		}
	})

	sd.list.contextmenu.rename_popup_widget.SetButtonText("Rename")
	sd.list.contextmenu.rename_popup_widget.SetFieldValue("Set new name")

	adder.AddChild(&sd.list.contextmenu.rename_popup_widget)
	adder.AddChild(&sd.list.contextmenu.menu)
	adder.AddChild(&sd.options.add.menu)
	adder.AddChild(&sd.options.add.folder_popup)
	return nil
}

func (sd *Sidebar[T]) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layouter.LayoutWidget(&sd.list.contextmenu.menu, image.Rectangle{
		Min: sd.list.contextmenu.position,
	})

	layouter.LayoutWidget(&sd.options.add.menu, image.Rectangle{
		Min: sd.options.add.add_button_pos,
	})

	form_measurements := sd.options.add.folder_popup.Measure(ctx, gui.Constraints{})
	layouter.LayoutWidget(&sd.options.add.folder_popup, image.Rectangle{
		Min: sd.options.add.add_button_pos,
		Max: sd.options.add.add_button_pos.Add(form_measurements),
	})

	rename_measurements := sd.list.contextmenu.rename_popup_widget.Measure(ctx, gui.Constraints{})
	layouter.LayoutWidget(&sd.list.contextmenu.rename_popup_widget, image.Rectangle{
		Min: sd.list.contextmenu.position,
		Max: sd.list.contextmenu.position.Add(rename_measurements),
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
							Widget: &sd.options.add.widget,
						},
					},
				},
			},
			{
				Widget: &sd.list.widget,
				Size:   gui.FlexibleSize(1),
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (sd *Sidebar[T]) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	return sd.list.contextmenu.menu.Measure(ctx, constraints)
}

func (sd *Sidebar[T]) HandlePointingInput(ctx *gui.Context, widgetBounds *gui.WidgetBounds) gui.HandleInputResult {
	if widgetBounds.IsHitAtCursor() && inpututil.IsMouseButtonJustReleased(ebiten.MouseButton2) && sd.list.contextmenu.right_clicked_item != nil {
		sd.list.contextmenu.position = image.Pt(ebiten.CursorPosition())
		sd.list.contextmenu.menu.SetOpen(true)

		fmt.Println("right clicked", sd.list.contextmenu.right_clicked_item.text_widget.Value())

		sd.list.contextmenu.right_clicked_item = nil
	}

	return gui.HandleInputResult{}
}
