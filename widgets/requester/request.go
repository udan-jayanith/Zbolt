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
				{
					HeaderText:                "",
					HeaderTextHorizontalAlign: widget.HorizontalAlignLeft,
					Width: gui.FixedSize(u),
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
				{
					HeaderText:                "",
					HeaderTextHorizontalAlign: widget.HorizontalAlignLeft,
					Width: gui.FixedSize(u),
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
				Name:   "Parameters",
				Widget: &rw.tab_content.params,
			},
			{
				Name:   "Headers",
				Widget: &rw.tab_content.header,
			}, {
				Name:   "Body",
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
				Size:   gui.FlexibleSize(1),
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}
