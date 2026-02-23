package request_page

import (
	"API-Client/basic"
	CommonWidgets "API-Client/common-widgets"
	"API-Client/icons"
	"API-Client/widgets/request/def"
	"API-Client/widgets/request/page/http"

	"image"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type RequestPage struct {
	gui.DefaultWidget

	background widget.Background

	sidebar       gui.WidgetWithPadding[*Sidebar[*def.Request]]
	sidebar_items []SidebarItem[*def.Request]

	tab_widget CommonWidgets.Tab[*def.Request]
	tab_items  []CommonWidgets.TabItem[*def.Request]

	request_widget CommonWidgets.WidgetWithPadding[def.RequestWidget]
	nothing_widget NothingWidget
	http_widget    http.HTTP_Widget

	popup_widget  widget.Popup
	popup_content sidebar_item_types_panel
	is_popup_open bool
}

func (rp *RequestPage) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	ctx.SetColorMode(gui.ColorModeDark)

	adder.AddChild(&rp.background)
	padding := basic.NewPadding(widget.UnitSize(ctx)/4, 0)

	sidebar := rp.sidebar.Widget()
	sidebar.SetItems(rp.sidebar_items)
	sidebar.OneItemsMoved(func(context *gui.Context, from, count, to int) {
		if to == len(rp.sidebar_items) {
			to--
		}

		f := rp.sidebar_items[from]
		rp.sidebar_items[from] = rp.sidebar_items[to]
		rp.sidebar_items[to] = f
	})

	rp.sidebar.SetPadding(padding)
	rp.sidebar.Widget().OnItemClicked(func(value *def.Request) {
		line_height := widget.LineHeight(ctx)
		size := line_height - line_height/4

		rp.tab_items = append(rp.tab_items, CommonWidgets.TabItem[*def.Request]{
			Value:    value,
			Text:     value.Name,
			Closable: true,
			Icon:     icons.NewIcon(value.Type.IconName(), size),
		})
	})
	adder.AddChild(&rp.sidebar)

	if len(rp.tab_items) > 0 {
		rp.tab_widget.SetTabItems(rp.tab_items)
		adder.AddChild(&rp.tab_widget)

		//_, req := rp.tab_widget.GetSelectedTab()
		rp.request_widget.SetWidget(&rp.http_widget)
		rp.request_widget.SetPadding(padding)
		adder.AddChild(&rp.request_widget)
	} else {
		adder.AddChild(&rp.nothing_widget)
	}

	rp.popup_widget.SetContent(&rp.popup_content)
	rp.popup_widget.SetBackgroundDark(true)
	rp.popup_widget.SetCloseByClickingOutside(true)
	rp.popup_widget.SetBackgroundBlurred(true)
	adder.AddChild(&rp.popup_widget)

	rp.sidebar.Widget().OnAddButtonClicked(func(ctx *gui.Context) {
		rp.popup_widget.SetOpen(true)
	})

	rp.popup_content.OnCreateButtonClicked(func(request *def.Request) {
		rp.sidebar_items = append(rp.sidebar_items, SidebarItem[*def.Request]{
			Value:    request,
			Text:     request.Name,
			IconName: request.Type.IconName(),
		})
		rp.popup_widget.SetOpen(false)
		rp.popup_content.Clear()
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
