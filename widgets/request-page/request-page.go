package Requester

import (
	"API-Client/basic"
	CommonWidgets "API-Client/common-widgets"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type RequestPage struct {
	gui.DefaultWidget

	background       widget.Background
	sidebar          gui.WidgetWithPadding[*Sidebar[struct{}]]
	tab_widget       CommonWidgets.Tab[struct{}]
	requester_widget gui.WidgetWithPadding[*HTTP_request]
}

func (rp *RequestPage) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	adder.AddChild(&rp.background)
	padding := basic.NewPadding(widget.UnitSize(ctx)/4, 0)

	sidebar := rp.sidebar.Widget()
	sidebar.SetItems([]SidebarItem[struct{}]{
		{
			IconName: "http",
			Text:     "product-data",
		},
		{
			Text: "update-product-data",
		},
		{
			Text: "search",
		},
		{
			Text: "discover",
		},
	})

	rp.sidebar.SetPadding(padding)
	adder.AddChild(&rp.sidebar)

	rp.tab_widget.SetTabItems([]CommonWidgets.TabItem[struct{}]{
		{
			Text: "product-data",
			Closable: true,
		},
		{
			Text: "discover",
			Closable: true,
		},
	})
	adder.AddChild(&rp.tab_widget)

	rp.requester_widget.SetPadding(padding)
	adder.AddChild(&rp.requester_widget)
	return nil
}

func (rp *RequestPage) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layouter.LayoutWidget(&rp.background, widgetBounds.Bounds())

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
