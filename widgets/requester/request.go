package Requester

import (
	"API-Client/basic"
	CommonWidgets "API-Client/common-widgets"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type RequestWidget struct {
	gui.DefaultWidget
	input_bar_widget RequestInputBar
	url_preview      widget.TextInput

	tab         CommonWidgets.Tab[string]
	tab_content struct {
		params, header           widget.Table[string]
		params_rows, header_rows []widget.TableRow[string]
		body                     widget.TextInput
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
				},
				{
					HeaderText:                "Value",
					HeaderTextHorizontalAlign: widget.HorizontalAlignLeft,
					MinWidth:                  u * 4,
				},
				{
					HeaderText:                "",
					HeaderTextHorizontalAlign: widget.HorizontalAlignLeft,
				},
			})

			if len(rw.tab_content.header_rows) == 0 {
				rw.tab_content.header_rows = append(rw.tab_content.header_rows, widget.TableRow[string]{
					Cells: []widget.TableCell{
						{},
						{},
						{
							Text: " +",
						},
					},
				})
			}
			rw.tab_content.header.SetItems(rw.tab_content.header_rows)
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
				{
					HeaderText:                "",
					HeaderTextHorizontalAlign: widget.HorizontalAlignLeft,
					Width:                     gui.FixedSize(u),
				},
			})

			if len(rw.tab_content.params_rows) == 0 {
				rw.tab_content.params_rows = append(rw.tab_content.params_rows, widget.TableRow[string]{
					Cells: []widget.TableCell{
						{},
						{},
						{
							Text: " +",
						},
					},
				})
			}
			rw.tab_content.params.SetItems(rw.tab_content.params_rows)
		}

		{
			rw.tab_content.body.SetAutoWrap(true)
			rw.tab_content.body.SetMultiline(true)
			rw.tab_content.body.SetEditable(true)
		}

		rw.tab.SetTabItems([]CommonWidgets.TabItem[string]{
			{
				Text: "Parameters",
			},
			{
				Text: "Headers",
			},
			{
				Text: "Body",
			},
			{
				Text: "Parameters Hello world",
			},
			{
				Text: "Headers",
			},
			{
				Text: "Body",
			},
		})
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
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}
