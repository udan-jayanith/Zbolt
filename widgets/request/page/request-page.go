package request_page

import (
	"API-Client/basic"
	CommonWidgets "API-Client/common-widgets"
	"API-Client/icons"
	"API-Client/widgets/request/def"
	http_widget "API-Client/widgets/request/page/http"
	websocket_widget "API-Client/widgets/request/page/websocket"
	"fmt"

	"image"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/hajimehoshi/ebiten/v2"
)

type sidebar_item struct {
	IsFolder bool
	Data     def.FolderOrFile
}

func (si *sidebar_item) GetPath() string {
	return si.Data.Path()
}

type RequestPage struct {
	gui.DefaultWidget

	background widget.Background

	sidebar gui.WidgetWithPadding[*Sidebar[sidebar_item]]

	// made sidebar_items map[string][]SidebarItem[sidebar_item]
	current_directory string
	sidebar_items     []SidebarItem[sidebar_item]

	tab_widget CommonWidgets.Tab[*def.Request]
	tab_items  []CommonWidgets.TabItem[*def.Request]

	nothing_widget NothingWidget

	request_widget   CommonWidgets.WidgetWithPadding[def.RequestWidget]
	http_widget      http_widget.HTTP_Widget
	websocket_widget websocket_widget.WebsocketWidget

	popup_widget  widget.Popup
	popup_content sidebar_item_types_panel
	is_popup_open bool

	notify_widget CommonWidgets.Notify
}

func (rp *RequestPage) open_tab(request *def.Request, ctx *gui.Context) error {
	for i, req := range rp.tab_items {
		if req.Value.Path() == request.Path() {
			rp.tab_widget.SelectTabItemByIndex(i)
			return fmt.Errorf("%s is already opened", request.Path())
		}
	}

	line_height := widget.LineHeight(ctx)
	size := line_height - line_height/4
	rp.tab_items = append(rp.tab_items, CommonWidgets.TabItem[*def.Request]{
		Value:    request,
		Text:     request.Name(),
		Closable: true,
		Icon:     icons.NewIcon(request.Type.IconName(), size),
	})

	rp.tab_widget.SelectTabItemByIndex(len(rp.tab_items) - 1)
	return nil
}

func (rp *RequestPage) create_sidebar_item(request *def.Request) {
	request_container := sidebar_item{
		Data: request,
	}
	rp.sidebar_items = append(rp.sidebar_items, SidebarItem[sidebar_item]{
		IconName: request.Type.IconName(),
		Text:     request.Name(),
		Value:    request_container,
	})

	//TODO: Open the the sidebar item if there were no items before on creation.
}

func (rp *RequestPage) create_folder(path string, name string) {
	folder := def.NewFolder(path, name)
	request_container := sidebar_item{
		IsFolder: true,
		Data:     &folder,
	}
	rp.sidebar_items = append(rp.sidebar_items, SidebarItem[sidebar_item]{
		IconName: "folder",
		Text:     folder.Name(),
		Value:    request_container,
	})

	//TODO: Open the the sidebar item if there were no items before on creation.
}

func (rp *RequestPage) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	ctx.SetColorMode(ebiten.ColorModeDark)

	adder.AddWidget(&rp.background)
	padding := basic.NewPadding(widget.UnitSize(ctx)/4, 0)

	sidebar := rp.sidebar.Widget()
	sidebar.SetItems(rp.sidebar_items)
	rp.sidebar.SetPadding(padding)

	sidebar.OnFolderCreate(func(ctx *gui.Context, folder_name string) {
		rp.create_folder(rp.current_directory, folder_name)
	})

	sidebar.OnItemClicked(func(item sidebar_item) {
		if item.IsFolder {
			return
		}

		request := item.Data.(*def.Request)
		err := rp.open_tab(request, ctx)
		if err != nil {
			rp.notify_widget.SetText(err.Error())
			rp.notify_widget.Open()
		}
	})
	adder.AddWidget(&rp.sidebar)

	if len(rp.tab_items) > 0 {
		rp.tab_widget.SetTabItems(rp.tab_items)
		rp.tab_widget.OnSwap(func(ctx *gui.Context, swap CommonWidgets.Swap) {
			temp := rp.tab_items[swap.From]
			rp.tab_items[swap.From] = rp.tab_items[swap.To]
			rp.tab_items[swap.To] = temp
		})
		adder.AddWidget(&rp.tab_widget)

		_, req := rp.tab_widget.GetSelectedTab()
		switch req.Type {
		case def.HTTP:
			rp.request_widget.SetWidget(&rp.http_widget)
		case def.Websocket:
			rp.request_widget.SetWidget(&rp.websocket_widget)
		default:
			panic("request type not handled")
		}

		rp.request_widget.SetPadding(padding)
		adder.AddWidget(&rp.request_widget)
	} else {
		rp.nothing_widget.OnClick(func() {
			rp.popup_widget.SetOpen(true)
		})
		adder.AddWidget(&rp.nothing_widget)
	}

	rp.popup_widget.SetContent(&rp.popup_content)
	rp.popup_widget.SetBackgroundDark(true)
	rp.popup_widget.SetCloseByClickingOutside(true)
	rp.popup_widget.SetBackgroundBlurred(true)
	adder.AddWidget(&rp.popup_widget)

	sidebar.OnAddButtonClicked(func(ctx *gui.Context) {
		rp.popup_widget.SetOpen(true)
	})

	rp.popup_content.OnCreateButtonClicked(func(request *def.Request) {
		rp.create_sidebar_item(request)
		rp.popup_widget.SetOpen(false)
		rp.popup_content.Clear()
	})

	adder.AddWidget(&rp.notify_widget)
	return nil
}

func (rp *RequestPage) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layouter.LayoutWidget(&rp.background, widgetBounds.Bounds())

	b := widgetBounds.Bounds()
	popup_content_bounds := rp.popup_content.Measure(ctx, gui.Constraints{})

	popup_size := image.Rectangle{
		Min: image.Pt(b.Min.X+b.Max.X/2-popup_content_bounds.X/2, b.Min.Y+b.Max.Y/2-popup_content_bounds.Y/2),
	}
	popup_size.Max = image.Pt(popup_size.Min.X+popup_content_bounds.X, popup_size.Min.Y+popup_content_bounds.Y)

	layouter.LayoutWidget(&rp.popup_widget, popup_size)

	rp.notify_widget.LayoutWidget(ctx, widgetBounds, layouter)

	tab_container_layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
	}

	if len(rp.tab_items) > 0 {
		tab_container_layout.Items = []gui.LinearLayoutItem{
			{
				Widget: &rp.tab_widget,
			},
			{
				Widget: &rp.request_widget,
				Size:   gui.FlexibleSize(1),
			},
		}
	} else {
		tab_container_layout.Items = append(tab_container_layout.Items, gui.LinearLayoutItem{
			Size:   gui.FlexibleSize(1),
			Widget: &rp.nothing_widget,
		})
	}

	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Gap:       widget.UnitSize(ctx) / 4,
		Items: []gui.LinearLayoutItem{
			{},
			{
				Widget: &rp.sidebar,
				Size:   gui.FlexibleSize(1),
			},
			{
				Layout: tab_container_layout,
				Size:   gui.FlexibleSize(4),
			},
			{},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}
