package Requester

import (
	"API-Client/basic"
	"API-Client/icons"
	"image"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type sidebar_item_type_card struct {
	gui.DefaultWidget

	icon_widget icons.Icon
	text_widget widget.Text

	Text, Icon_name string
}

func (sitc *sidebar_item_type_card) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	sitc.icon_widget.IconName = "large-icons/" + sitc.Icon_name
	u := widget.UnitSize(ctx)
	icon_size := image.Pt(u*2, u*2)
	sitc.icon_widget.Point = &icon_size
	adder.AddChild(&sitc.icon_widget)

	sitc.text_widget.SetValue(sitc.Text)
	sitc.text_widget.SetVerticalAlign(widget.VerticalAlignMiddle)
	sitc.text_widget.SetHorizontalAlign(widget.HorizontalAlignCenter)
	adder.AddChild(&sitc.text_widget)

	return nil
}

func (sitc *sidebar_item_type_card) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	size := widget.UnitSize(ctx) / 4
	layout := gui.LinearLayout{
		Gap:       size,
		Padding:   basic.NewPadding(size),
		Direction: gui.LayoutDirectionVertical,

		Items: []gui.LinearLayoutItem{
			{
				Widget: &sitc.icon_widget,
			},
			{
				Widget: &sitc.text_widget,
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (sitc *sidebar_item_type_card) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	var point image.Point
	size := widget.UnitSize(ctx) / 4

	if h, ok := constraints.FixedHeight(); ok {
		point.Y = h
	} else {
		point.Y = sitc.icon_widget.Measure(ctx, constraints).Y
		point.Y += sitc.text_widget.Measure(ctx, constraints).Y
		point.Y += size * 3
	}

	if w, ok := constraints.FixedWidth(); ok {
		point.X = w
	} else {
		point.X = point.Y
	}

	return point
}

type sidebar_item_types_panel struct {
	gui.DefaultWidget

	http, websocket, graphql, grpc gui.WidgetWithSize[*sidebar_item_type_card]
	select_type_text_widget        widget.Text
}

func (sitp *sidebar_item_types_panel) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	sitp.select_type_text_widget.SetValue("Select request type")
	sitp.select_type_text_widget.SetVerticalAlign(widget.VerticalAlignMiddle)
	sitp.select_type_text_widget.SetHorizontalAlign(widget.HorizontalAlignCenter)
	sitp.select_type_text_widget.SetScale(1.2)
	adder.AddChild(&sitp.select_type_text_widget)

	u := widget.UnitSize(ctx)
	item_size := image.Pt(u*4, u*4)

	sitp.http.SetFixedSize(item_size)
	http := sitp.http.Widget()
	http.Text = "HTTP"
	http.Icon_name = "http"
	adder.AddChild(&sitp.http)

	sitp.websocket.SetFixedSize(item_size)
	websocket := sitp.websocket.Widget()
	websocket.Text = "Websocket"
	websocket.Icon_name = "websocket"
	adder.AddChild(&sitp.websocket)

	sitp.graphql.SetFixedSize(item_size)
	graphql := sitp.graphql.Widget()
	graphql.Text = "GraphQL"
	graphql.Icon_name = "graphql"
	adder.AddChild(&sitp.graphql)

	sitp.grpc.SetFixedSize(item_size)
	grpc := sitp.grpc.Widget()
	grpc.Text = "gRPC"
	grpc.Icon_name = "grpc"
	adder.AddChild(&sitp.grpc)

	return nil
}

func (sitp *sidebar_item_types_panel) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	size := widget.UnitSize(ctx) / 4
	layout := gui.LinearLayout{
		Gap:       size,
		Padding:   basic.NewPadding(size),
		Direction: gui.LayoutDirectionVertical,

		Items: []gui.LinearLayoutItem{
			{},
			{
				Widget: &sitp.select_type_text_widget,
				Size: gui.FlexibleSize(1),
			},
			{
				Layout: gui.LinearLayout{
					Gap:       size,
					Direction: gui.LayoutDirectionHorizontal,

					Items: []gui.LinearLayoutItem{
						{
							Widget: &sitp.http,
						},
						{
							Widget: &sitp.websocket,
						},
						{
							Widget: &sitp.graphql,
						},
						{
							Widget: &sitp.grpc,
						},
					},
				},
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (sitp *sidebar_item_types_panel) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	u := widget.UnitSize(ctx)
	
	gap, padding := u/4, u/4
	width := sitp.http.Measure(ctx, gui.Constraints{}).X*4+gap*3+padding*2
	height := (u*4)+(u/2)+gap*2+padding*2+sitp.select_type_text_widget.Measure(ctx, gui.Constraints{}).Y
	point := image.Pt(width, height)

	if h, ok := constraints.FixedHeight(); ok {
		point.Y = h
	}

	if w, ok := constraints.FixedWidth(); ok {
		point.X = w
	}

	return point
}
