package http_widget

import (
	CommonWidgets "API-Client/common-widgets"
	attr "API-Client/widgets/request"
	"API-Client/widgets/request/def"
	"image"
	"net/url"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type request_widget struct {
	gui.DefaultWidget
	input_bar_widget          request_input_bar_widget
	on_input_bar_value_change func(ctx *gui.Context, text string, committed bool, by_user bool)

	url_preview CommonWidgets.URLPreview

	tab         CommonWidgets.Tab[string]
	tab_content struct {
		params, header  []attr.AttrCheck
		table           CommonWidgets.AttributeTable // TODO: use one attribute table for params and headers.
		body            CommonWidgets.BodyWidget
		selected_widget gui.Widget
	}
}

// sets the http method
func (rw *request_widget) SetMethod(method string) {
	rw.input_bar_widget.select_method(method)
}

func (rw *request_widget) OnMethodChanged(fn func(method string)){
	rw.input_bar_widget.on_method_changed(fn)
}

func (rw *request_widget) OnOpenIn(fn func(ctx *gui.Context)) {
	rw.input_bar_widget.on_open_in_clicked(fn)
}

// This should be only used to set the url base including the path
func (rw *request_widget) SetURL_str(url string) {
	rw.input_bar_widget.set_url_input_value(url)
}

func (rw *request_widget) SetURL(u *url.URL) {
	raw_query := u.RawQuery
	u.RawQuery = ""
	// TODO: process the query and merge it with the parameter table.
	rw.input_bar_widget.set_url_input_value(u.String())

	// TODO: encode the parameter table and join it with the url base.
	u.RawQuery = raw_query
	rw.url_preview.SetURL(u.String()) // TODO: remove this once update_url_preview is finished implemented and url preview get updated every second.
}

func (rw *request_widget) update_url_preview() {

}

func (rw *request_widget) OnRequestButtonClicked(fn func(ctx *gui.Context, value string)) {
	rw.input_bar_widget.on_request_button_clicked(fn)
}

func (rw *request_widget) DisableURLInput(disabled bool) {
	rw.input_bar_widget.disable_url_input(disabled)
}

func (rw *request_widget) OnURLInputChanged(fn func(context *gui.Context, text string, committed bool)) {
	rw.input_bar_widget.on_url_input_value_changed(fn)
}

func (rw *request_widget) SetParameters(parameters []attr.AttrCheck) {
	_, selected_tab_value := rw.tab.GetSelectedTab()
	if selected_tab_value == "parameters" {
		rw.tab_content.table.SetRowsCheck(parameters)
	}

	rw.tab_content.params = parameters
}

func (rw *request_widget) Parameters() []attr.AttrCheck {
	_, selected_tab_value := rw.tab.GetSelectedTab()
	if selected_tab_value == "parameters" {
		rw.tab_content.params = rw.tab_content.table.RowsCheck()
	}
	return rw.tab_content.params
}

func (rw *request_widget) SetHeaders(headers []attr.AttrCheck) {
	_, selected_tab_value := rw.tab.GetSelectedTab()
	if selected_tab_value == "headers" {
		rw.tab_content.table.SetRowsCheck(headers)
	}
	rw.tab_content.header = headers
}

func (rw *request_widget) Headers() []attr.AttrCheck {
	_, selected_tab_value := rw.tab.GetSelectedTab()
	if selected_tab_value == "headers" {
		rw.tab_content.header = rw.tab_content.table.RowsCheck()
	}
	return rw.tab_content.header
}

// TODO: implement a function to receive body widget toggles values.

// SetBody set the http request body
func (rw *request_widget) SetBody(body *def.HTTP_Request_Body) {
	rw.tab_content.body.SetBody(body.Content, def.ContentType(body.ContentType))
}

func (rw *request_widget) SelectedTab() int {
	i, _ := rw.tab.GetSelectedTab()
	return i
}

func (rw *request_widget) SetTab(index int) {
	rw.tab.SelectTabItemByIndex(index)
}

func (rw *request_widget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	adder.AddWidget(&rw.input_bar_widget)
	adder.AddWidget(&rw.url_preview)

	rw.tab.SetTabItems([]CommonWidgets.TabItem[string]{
		{
			Text:  "Parameters",
			Value: "parameters",
		},
		{
			Text:  "Headers",
			Value: "headers",
		},
		{
			Text:  "Body",
			Value: "body",
		},
	})

	rw.tab.OnSwitch(func(from, to *CommonWidgets.TabItem[string]) {
		if from.Value == "parameters" && to.Value == "headers" {
			rw.tab_content.params = rw.tab_content.table.RowsCheck()
		} else if from.Value == "headers" && to.Value == "parameters" {
			rw.tab_content.header = rw.tab_content.table.RowsCheck()
		}

		if to.Value == "parameters" {
			rw.tab_content.table.SetRowsCheck(rw.tab_content.params)
		} else if to.Value == "headers" {
			rw.tab_content.table.SetRowsCheck(rw.tab_content.header)
		}
	})

	_, selected_tab_value := rw.tab.GetSelectedTab()
	switch selected_tab_value {
	case "parameters", "headers":
		rw.tab_content.selected_widget = &rw.tab_content.table
	case "body":
		rw.tab_content.selected_widget = &rw.tab_content.body
		rw.tab_content.body.SetType(CommonWidgets.HTTP_Request)
	default:
		panic("Unknown tab was selected")
	}

	adder.AddWidget(rw.tab_content.selected_widget)
	adder.AddWidget(&rw.tab)
	return nil
}

func (rw *request_widget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	u := widget.UnitSize(ctx)
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Gap:       u / 4,
		Items: []gui.LinearLayoutItem{
			{
				Widget: &rw.input_bar_widget,
				Size:   gui.FixedSize(u),
			},
			{
				Widget: &rw.url_preview,
			},
			{
				Widget: &rw.tab,
			},
			{
				Widget: rw.tab_content.selected_widget,
				Size:   gui.FlexibleSize(1),
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (rw *request_widget) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	point := rw.input_bar_widget.Measure(ctx, constraints)

	if h, ok := constraints.FixedHeight(); ok {
		point.Y = h
	} else {
		point.Y += rw.url_preview.Measure(ctx, constraints).Y
		point.Y += rw.tab.Measure(ctx, constraints).Y
		point.Y += rw.tab_content.selected_widget.Measure(ctx, constraints).Y
	}

	if w, ok := constraints.FixedWidth(); ok {
		point.X = w
	}

	return point
}
