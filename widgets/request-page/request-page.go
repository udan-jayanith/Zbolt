package Requester

import (
	"API-Client/basic"
	CommonWidgets "API-Client/common-widgets"
	"weak"

	"image"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type RequestType uint8

const (
	HTTP RequestType = iota + 0
	Websocket
	GraphQL
	Grpc
)

type Request struct {
	Type RequestType
	Name string
	data weak.Pointer[any]
}

func (r *Request) Data() any {
	return nil
}

type RequestPage struct {
	gui.DefaultWidget

	background widget.Background

	sidebar       gui.WidgetWithPadding[*Sidebar[Request]]
	sidebar_items []SidebarItem[Request]

	tab_widget       CommonWidgets.Tab[Request]
	requester_widget gui.WidgetWithPadding[*HTTP_request]

	popup_widget  widget.Popup
	popup_content sidebar_item_types_panel
	is_popup_open bool
}

func (rp *RequestPage) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	adder.AddChild(&rp.background)
	padding := basic.NewPadding(widget.UnitSize(ctx)/4, 0)

	sidebar := rp.sidebar.Widget()
	sidebar.SetItems(rp.sidebar_items)

	rp.sidebar.SetPadding(padding)
	adder.AddChild(&rp.sidebar)

	rp.tab_widget.SetTabItems([]CommonWidgets.TabItem[Request]{
		{
			Text:     "product-data",
			Closable: true,
		},
		{
			Text:     "discover",
			Closable: true,
		},
	})
	adder.AddChild(&rp.tab_widget)

	rp.requester_widget.SetPadding(padding)
	adder.AddChild(&rp.requester_widget)

	rp.popup_widget.SetContent(&rp.popup_content)
	rp.popup_widget.SetBackgroundDark(true)
	rp.popup_widget.SetCloseByClickingOutside(true)
	rp.popup_widget.SetBackgroundBlurred(true)
	adder.AddChild(&rp.popup_widget)

	rp.sidebar.Widget().OnAddButtonClicked(func(ctx *gui.Context) {
		rp.popup_widget.SetOpen(true)
	})
	
	rp.popup_content.OnCreateButtonClicked(func(request Request) {
		rp.sidebar_items = append(rp.sidebar_items, SidebarItem[Request]{
			Value: request,
			Text: request.Name,
		})
		rp.popup_widget.SetOpen(false)
	})

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
				Layout: gui.LinearLayout{
					Direction: gui.LayoutDirectionVertical,
					Items: []gui.LinearLayoutItem{
						{
							Widget: &rp.tab_widget,
						},
						{
							Widget: &rp.requester_widget,
							Size:   gui.FlexibleSize(1),
						},
					},
				},
				Size: gui.FlexibleSize(4),
			},
			{},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}
