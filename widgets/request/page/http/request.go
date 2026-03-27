package http_widget

import (
	CommonWidgets "API-Client/common-widgets"
	"API-Client/widgets/request/def"
	url_pattern "API-Client/widgets/request/url-pattern"
	"image"
	"net/url"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type request_widget struct {
	gui.DefaultWidget
	input_bar_widget request_input_bar_widget
	url_preview      CommonWidgets.URLPreview

	tab         CommonWidgets.Tab[string]
	tab_content struct {
		params, header  CommonWidgets.AttributeTable
		body            widget.TextInput
		selected_widget gui.Widget
	}
}

// sets the http method
func (rw *request_widget) SetMethod(method string) {
	rw.input_bar_widget.method_select_widget.SelectItemByValue(method)
}

func (rw *request_widget) Method() string {
	selected_item, _ := rw.input_bar_widget.method_select_widget.SelectedItem()
	return selected_item.Value
}

func (rw *request_widget) SetURL(u *url.URL) {
	raw_query := u.RawQuery
	u.RawQuery = ""
	rw.input_bar_widget.input_widget.SetValue(u.String())

	u.RawQuery = raw_query
	rw.url_preview.SetURL(u.String())
}

// URL returns the url without parameters
func (rw *request_widget) URL() *url.URL {
	u, _ := url.Parse(rw.input_bar_widget.input_widget.Value())
	return u
}

func (rw *request_widget) SetParameters(parameters []url_pattern.Attribute) {
}

func (rw *request_widget) Parameters() []url_pattern.Attribute {
	return []url_pattern.Attribute{}
}

func (rw *request_widget) SetHeaders(headers []url_pattern.Attribute) {
}

func (rw *request_widget) Headers() []url_pattern.Attribute {
	return []url_pattern.Attribute{}
}

func (rw *request_widget) SetBody(body *def.HTTP_Request_Body) {
}

func (rw *request_widget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	rw.input_bar_widget.OnRequest(func(ctx *gui.Context, url, method string) {
		//http.Request{
		//Method: strings.ToUpper(method),
		//}
	})
	adder.AddWidget(&rw.input_bar_widget)

	adder.AddWidget(&rw.url_preview)

	{
		rw.tab_content.body.SetAutoWrap(true)
		rw.tab_content.body.SetMultiline(true)
		rw.tab_content.body.SetEditable(true)

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
		})

		switch rw.tab.GetSelectedIndex() {
		case 0:
			rw.tab_content.selected_widget = &rw.tab_content.params
		case 1:
			rw.tab_content.selected_widget = &rw.tab_content.header
		case 2:
			rw.tab_content.selected_widget = &rw.tab_content.body
		default:
			panic("Unknown tab was selected")
		}

		adder.AddWidget(rw.tab_content.selected_widget)
		adder.AddWidget(&rw.tab)
	}
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
