package Requester

import (
	"API-Client/basic"
	CommonWidgets "API-Client/common-widgets"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type ResponseWidget struct {
	gui.DefaultWidget
	header struct {
		status, response_time, size, proto widget.Text
	}
	tab         CommonWidgets.Tab[uint8]
	tab_content struct {
		response_header widget.Table[string]
		response_body   widget.TextInput
	}
}

func (rw *ResponseWidget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	{
		rw.header.status.SetTabular(true)
		rw.header.status.SetValue("200 Ok")
		adder.AddChild(&rw.header.status)

		rw.header.response_time.SetTabular(true)
		rw.header.response_time.SetValue("200 ms")
		adder.AddChild(&rw.header.response_time)

		rw.header.size.SetTabular(true)
		rw.header.size.SetValue("131 B")
		adder.AddChild(&rw.header.size)

		rw.header.proto.SetTabular(true)
		rw.header.proto.SetValue("HTTP v1.1")
		adder.AddChild(&rw.header.proto)
	}

	{
		{
			rw.tab_content.response_body.SetAutoWrap(true)
			rw.tab_content.response_body.SetMultiline(true)
			rw.tab_content.response_body.SetEditable(false)
			rw.tab_content.response_body.SetValue(`
			git clone https://github.com/guigui-gui/guigui.git
			cd guigui
			go run ./example/gallery


			hi


			Hello world
			`)
		}

		u := widget.UnitSize(ctx)
		{
			rw.tab_content.response_header.SetColumns([]widget.TableColumn{
				{
					HeaderText:                "Name",
					HeaderTextHorizontalAlign: widget.HorizontalAlignLeft,
					MinWidth:                  u * 4,
					Width:                     gui.FlexibleSize(1),
				},
				{
					HeaderText:                "Value",
					HeaderTextHorizontalAlign: widget.HorizontalAlignLeft,
					MinWidth:                  u * 4,
					Width:                     gui.FlexibleSize(1),
				},
			})
			rw.tab_content.response_header.SetItems([]widget.TableRow[string]{
				{
					Cells: []widget.TableCell{
						{Text: "Content-Type"},
						{Text: "test/json"},
					},
				},
				{
					Cells: []widget.TableCell{
						{Text: "Content-Length"},
						{Text: "141"},
					},
				},
			})
		}

		rw.tab.Tab_Items = []CommonWidgets.TabItem[uint8]{
			{
				Widget: &rw.tab_content.response_body,
				Name:   "Body",
			},
			{
				Widget: &rw.tab_content.response_header,
				Name:   "Header",
			},
		}

		adder.AddChild(&rw.tab)
	}

	return nil
}

func (rw *ResponseWidget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	u := widget.UnitSize(ctx)
	header_layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Gap:       u / 4,
		Padding:   basic.NewPadding(u / 4),
		Items: []gui.LinearLayoutItem{
			{
				Widget: &rw.header.status,
			},
			{
				Widget: &rw.header.response_time,
			},
			{
				Widget: &rw.header.size,
			},
			{
				Size: gui.FlexibleSize(1),
			},
			{
				Widget: &rw.header.proto,
			},
		},
	}

	main_layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Gap:       u / 4,
		Padding:   basic.NewPadding(u / 4),
		Items: []gui.LinearLayoutItem{
			{
				Layout: header_layout,
			},
			{
				Widget: &rw.tab,
				Size:   gui.FlexibleSize(1),
			},
		},
	}

	main_layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}