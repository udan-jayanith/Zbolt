package RequestPage

import (
	"image"

	"API-Client/basic"
	CommonWidgets "API-Client/common-widgets"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type RequestInputBar struct {
	gui.DefaultWidget
	method_select_widget widget.Select[string]
	input_widget         widget.TextInput
	request_btn_widget   widget.Button
}

func (rib *RequestInputBar) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	rib.method_select_widget.SetItemsByStrings([]string{
		"Get",
		"Post",
		"Put",
		"Patch",
		"Delete",
		"Options",
		"Head",
	})

	selected_item_index := max(rib.method_select_widget.SelectedItemIndex(), 0)
	rib.method_select_widget.SelectItemByIndex(selected_item_index)
	adder.AddChild(&rib.method_select_widget)

	rib.input_widget.SetEditable(true)
	adder.AddChild(&rib.input_widget)

	rib.request_btn_widget.SetText("Request")
	rib.request_btn_widget.SetTextBold(true)
	adder.AddChild(&rib.request_btn_widget)
	return nil
}

func (rib *RequestInputBar) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	u := widget.UnitSize(ctx)
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Items: []gui.LinearLayoutItem{
			{
				Size: gui.FixedSize(u),
				Layout: gui.LinearLayout{
					Direction: gui.LayoutDirectionHorizontal,
					Gap:       u / 4,
					Items: []gui.LinearLayoutItem{
						{
							Widget: &rib.method_select_widget,
						},
						{
							Widget: &rib.input_widget,
							Size:   gui.FlexibleSize(1),
						},
						{
							Widget: &rib.request_btn_widget,
						},
					},
				},
			},
		},
	}

	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (rib *RequestInputBar) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	u := widget.UnitSize(ctx)
	if w, ok := constraints.FixedWidth(); ok {
		return image.Pt(w, u*2)
	} else if h, ok := constraints.FixedHeight(); ok {
		return image.Pt(u*10, h)
	}
	return image.Pt(u*10, u*2)
}

type RequestWidget struct {
	gui.DefaultWidget
	input_bar_widget RequestInputBar
	url_preview      widget.TextInput

	tab         CommonWidgets.Tab[string]
	tab_content struct {
		params, header widget.Table[string]
		body           widget.TextInput
	}
}

func (rw *RequestWidget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	adder.AddChild(&rw.input_bar_widget)

	{
		rw.url_preview.SetEditable(false)
		rw.url_preview.SetValue("https://github.com/guigui-gui/guigui/issues?q=is%3Aissue%20state%3Aopen%20milestone%3Av0.1.0&page=2")
		rw.url_preview.SetMultiline(true)
		rw.url_preview.SetAutoWrap(true)
		adder.AddChild(&rw.url_preview)
	}

	u := widget.UnitSize(ctx)
	{
		{
			rw.tab_content.header.SetColumns([]widget.TableColumn{
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
		}

		{
			rw.tab_content.params.SetColumns([]widget.TableColumn{
				{
					HeaderText:                "Attribute name",
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
		}

		{
			rw.tab_content.body.SetAutoWrap(true)
			rw.tab_content.body.SetMultiline(true)
			rw.tab_content.body.SetEditable(true)
		}

		rw.tab.Tab_Items = []CommonWidgets.TabItem[string]{
			{
				Name: "Parameters",
				Widget: &rw.tab_content.params,
			},
			{
				Name: "Headers",
				Widget: &rw.tab_content.header,
			},			{
				Name: "Body",
				Widget: &rw.tab_content.body,
			},
		}
		
		adder.AddChild(&rw.tab)
	}
	return nil
}

func (rw *RequestWidget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	u := widget.UnitSize(ctx)
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Gap:       u / 4,
		Padding:   basic.NewPadding(u / 4),
		Items: []gui.LinearLayoutItem{
			{
				Widget: &rw.input_bar_widget,
				Size:   gui.FixedSize(u),
			},
			{
				Widget: &rw.url_preview,
				Size:   gui.FixedSize(u * 2),
			},
			{
				Widget: &rw.tab,
				Size: gui.FlexibleSize(1),
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

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

type BasicPage struct {
	gui.DefaultWidget
	background widget.Background
	panel      struct {
		request struct {
			panel   widget.Panel
			content gui.WidgetWithSize[*RequestWidget]
		}
		response struct {
			panel   widget.Panel
			content gui.WidgetWithSize[*ResponseWidget]
		}
	}
}

func (brp *BasicPage) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	ctx.SetColorMode(gui.ColorModeDark)
	adder.AddChild(&brp.background)

	brp.panel.request.panel.SetContent(&brp.panel.request.content)
	brp.panel.request.panel.SetBorders(widget.PanelBorders{
		End: true,
	})
	adder.AddChild(&brp.panel.request.panel)

	brp.panel.response.panel.SetContent(&brp.panel.response.content)
	adder.AddChild(&brp.panel.response.panel)
	return nil
}

func (brp *BasicPage) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layouter.LayoutWidget(&brp.background, widgetBounds.Bounds())
	b := widgetBounds.Bounds()

	panel_size := b.Max
	panel_size.X = panel_size.X / 2
	brp.panel.request.content.SetFixedSize(panel_size)
	brp.panel.response.content.SetFixedSize(panel_size)

	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Items: []gui.LinearLayoutItem{
			{
				Widget: &brp.panel.request.panel,
				Size:   gui.FlexibleSize(1),
			},
			{
				Widget: &brp.panel.response.panel,
				Size:   gui.FlexibleSize(1),
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}
