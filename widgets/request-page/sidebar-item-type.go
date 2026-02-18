package Requester

import (
	"API-Client/basic"
	"API-Client/icons"
	"image"
	"log"

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

	http, websocket, graphql, grpc gui.WidgetWithSize[*widget.Button]
	selected_request_type          RequestType

	select_type_text_widget widget.Text
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
	http.SetContent(&sidebar_item_type_card{
		Text:      "HTTP",
		Icon_name: "http",
	})
	http.SetOnDown(func(context *gui.Context) {
		sitp.selected_request_type = HTTP
	})
	adder.AddChild(&sitp.http)

	sitp.websocket.SetFixedSize(item_size)
	websocket := sitp.websocket.Widget()
	websocket.SetContent(&sidebar_item_type_card{
		Text:      "Websocket",
		Icon_name: "websocket",
	})
	websocket.SetOnDown(func(context *gui.Context) {
		sitp.selected_request_type = Websocket
	})

	adder.AddChild(&sitp.websocket)

	sitp.graphql.SetFixedSize(item_size)
	graphql := sitp.graphql.Widget()
	graphql.SetContent(&sidebar_item_type_card{
		Text:      "GraphQL",
		Icon_name: "graphql",
	})
	graphql.SetOnDown(func(context *gui.Context) {
		sitp.selected_request_type = GraphQL
	})
	adder.AddChild(&sitp.graphql)

	sitp.grpc.SetFixedSize(item_size)
	grpc := sitp.grpc.Widget()
	grpc.SetContent(&sidebar_item_type_card{
		Text:      "gRPC",
		Icon_name: "grpc",
	})
	grpc.SetOnDown(func(context *gui.Context) {
		sitp.selected_request_type = Grpc
	})
	adder.AddChild(&sitp.grpc)

	http.SetType(widget.ButtonTypeNormal)
	websocket.SetType(widget.ButtonTypeNormal)
	graphql.SetType(widget.ButtonTypeNormal)
	grpc.SetType(widget.ButtonTypeNormal)

	switch sitp.selected_request_type {
	case HTTP:
		http.SetType(widget.ButtonTypePrimary)
	case Websocket:
		websocket.SetType(widget.ButtonTypePrimary)
	case GraphQL:
		graphql.SetType(widget.ButtonTypePrimary)
	case Grpc:
		grpc.SetType(widget.ButtonTypePrimary)
	default:
		log.Fatal("Unknown request type selected")
	}
	return nil
}

func (sitp *sidebar_item_types_panel) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	size := widget.UnitSize(ctx) / 4
	layout := gui.LinearLayout{
		Gap:       size,
		Padding:   basic.NewPadding(size * 2),
		Direction: gui.LayoutDirectionVertical,

		Items: []gui.LinearLayoutItem{
			{
				Widget: &sitp.select_type_text_widget,
				Size:   gui.FlexibleSize(1),
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

	gap, padding := u/4, u/2
	width := sitp.http.Measure(ctx, gui.Constraints{}).X*4 + gap*3 + padding*2
	height := (u * 4) + (u / 2) + gap + padding*2 + sitp.select_type_text_widget.Measure(ctx, gui.Constraints{}).Y
	point := image.Pt(width, height)

	if h, ok := constraints.FixedHeight(); ok {
		point.Y = h
	}

	if w, ok := constraints.FixedWidth(); ok {
		point.X = w
	}

	return point
}
