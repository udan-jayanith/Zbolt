package request_page

import (
	"API-Client/basic"
	CommonWidgets "API-Client/common-widgets"
	"API-Client/widgets/request/def"
	http_widget "API-Client/widgets/request/page/http"
	websocket_widget "API-Client/widgets/request/page/websocket"

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

	tabs_handler   TabsHandler
	nothing_widget NothingWidget

	request_widget   CommonWidgets.WidgetWithPadding[def.RequestWidget]
	http_widget      http_widget.HTTP_Widget
	websocket_widget websocket_widget.WebsocketWidget

	request_create_widget sidebar_item_types_panel
	variable_panel_widget variable_panel_widget

	popup_content gui.Widget
	popup_size    image.Point
	popup_widget  widget.Popup

	notify_widget CommonWidgets.Notify
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

func (rp *RequestPage) open_popup(widget gui.Widget, ctx *gui.Context) {
	rp.popup_content = widget
	rp.popup_size = widget.Measure(ctx, gui.Constraints{})
	rp.popup_widget.SetContent(rp.popup_content)
	rp.popup_widget.SetOpen(true)
}

func (rp *RequestPage) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	ctx.SetPreferredColorMode(ebiten.ColorModeDark)

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
		rp.tabs_handler.Open(request, ctx)
	})
	adder.AddWidget(&rp.sidebar)

	if rp.tabs_handler.IsEmpty() {
		rp.nothing_widget.OnClick(func() {
			rp.request_create_widget.Clear()
			rp.open_popup(&rp.request_create_widget, ctx)
		})
		adder.AddWidget(&rp.nothing_widget)
	} else {
		rp.tabs_handler.Add(adder)

		data := rp.tabs_handler.GetData(rp.tabs_handler.SelectedTab())
		if data == nil {
			panic("Invalid")
		}
		switch data.Type {
		case def.HTTP:
			rp.http_widget.SetReq(data)
			rp.request_widget.SetWidget(&rp.http_widget)
		case def.Websocket:
			rp.websocket_widget.SetReq(data)
			rp.request_widget.SetWidget(&rp.websocket_widget)
		default:
			panic("request type not handled")
		}

		rp.request_widget.SetPadding(padding)
		adder.AddWidget(&rp.request_widget)
	}

	rp.popup_widget.SetBackgroundDark(true)
	rp.popup_widget.SetCloseByClickingOutside(true)
	rp.popup_widget.SetBackgroundBlurred(true)
	adder.AddWidget(&rp.popup_widget)

	sidebar.OnRequestCreate(func(ctx *gui.Context) {
		rp.request_create_widget.Clear()
		rp.open_popup(&rp.request_create_widget, ctx)
	})

	rp.request_create_widget.OnCreateButtonClicked(func(request *def.Request) {
		rp.create_sidebar_item(request)
		rp.popup_widget.SetOpen(false)
	})

	rp.sidebar.Widget().OnVariableClicked(func(ctx *gui.Context) {
		rp.open_popup(&rp.variable_panel_widget, ctx)
	})

	adder.AddWidget(&rp.notify_widget)
	return nil
}

func (rp *RequestPage) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layouter.LayoutWidget(&rp.background, widgetBounds.Bounds())

	b := widgetBounds.Bounds()
	if rp.popup_widget.IsOpen() {
		popup_content_bounds := rp.popup_size

		popup_size := image.Rectangle{
			Min: image.Pt(b.Min.X+b.Max.X/2-popup_content_bounds.X/2, b.Min.Y+b.Max.Y/2-popup_content_bounds.Y/2),
		}
		popup_size.Max = image.Pt(popup_size.Min.X+popup_content_bounds.X, popup_size.Min.Y+popup_content_bounds.Y)

		layouter.LayoutWidget(&rp.popup_widget, popup_size)
	}
	rp.notify_widget.LayoutWidget(ctx, widgetBounds, layouter)

	tab_container_layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
	}

	if !rp.tabs_handler.IsEmpty() {
		tab_container_layout.Items = []gui.LinearLayoutItem{
			{
				Widget: rp.tabs_handler.Widget(),
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
