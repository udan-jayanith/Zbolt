package Requester

import (
	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type RequestPage struct{
	gui.DefaultWidget

	background widget.Background
	sidebar Sidebar[struct{}]
	requester_widget Requester
}

func (rp *RequestPage) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	adder.AddChild(&rp.background)
	
	rp.sidebar.SetItems([]SidebarItem[struct{}]{
		{
			IconName: "http",
			Text: "product-data",
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
	adder.AddChild(&rp.sidebar)
	
	adder.AddChild(&rp.requester_widget)
	return nil
}

func (rp *RequestPage) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layouter.LayoutWidget(&rp.background, widgetBounds.Bounds())
	
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Items: []gui.LinearLayoutItem{
			{
				Widget: &rp.sidebar,
				Size: gui.FlexibleSize(1),
			},
			{
				Widget: &rp.requester_widget,
				Size: gui.FlexibleSize(4),
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}
