package http_widget

import (
	CommonWidgets "API-Client/common-widgets"
	attr "API-Client/widgets/request/attributes"
	"API-Client/widgets/request/def"
	"fmt"
	"net/http"
	"time"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type response_widget struct {
	gui.DefaultWidget
	header_widget response_header_widget
	tab           CommonWidgets.Tab[struct{}]
	tab_content   struct {
		response_header  widget.Table[struct{}]
		response_body    CommonWidgets.BodyWidget
		selected_content gui.Widget
	}
}

func (rw *response_widget) OnAutowrapToggle(fn func(ctx *gui.Context, value bool)) {
	rw.tab_content.response_body.OnAutowrapToggle(fn)
}

func (rw *response_widget) OnFormatToggle(fn func(ctx *gui.Context, value bool)) {
	rw.tab_content.response_body.OnFormatToggle(fn)
}

func (rw *response_widget) SetAutowrap(autowrap bool) {
	rw.tab_content.response_body.SetAutowrap(autowrap)
}

func (rw *response_widget) SetFormat(format bool) {
	rw.tab_content.response_body.SetFormat(format)
}

// SetStatus sets the http status code
func (rw *response_widget) SetStatus(status_code int) {
	status := fmt.Sprintf("%v %s", status_code, http.StatusText(status_code))
	rw.header_widget.status.SetValue(status)
}

// SetResponseTime sets the http response time
func (rw *response_widget) SetResponseTime(response_time time.Duration) {
	rw.header_widget.response_time.SetValue(response_time.String())
}

func (rw *response_widget) SetHTTPVersion(version def.Version) {
	rw.header_widget.proto.SetValue(fmt.Sprintf("HTTP v%v.%v", version.Major, version.Minor))
}

func (rw *response_widget) SetHeaders(headers []attr.AttrCheck) {
	header_items := make([]widget.TableRow[struct{}], 0, len(headers))

	for _, v := range headers {
		header_items = append(header_items, widget.TableRow[struct{}]{
			Cells: []widget.TableCell{
				{
					Text: v.Key,
				},
				{
					Text: v.Key,
				},
			},
		})
	}

	rw.tab_content.response_header.SetItems(header_items)
}

func (rw *response_widget) SetResponseBody(body *def.HTTP_Response_Body) {
	if body.File == nil {
		// TODO: handle this so that images and large files can render.
		rw.tab_content.response_body.SetBody(body.Content, body.ContentType)
	}

	// If file is not nil and the content type is jpg, png or a text format show it in the response body widget.
	// Other wise show not unable to open and close the file. User should be able to click the open with button to view it.
}

func (rw *response_widget) SetSelectedTab(index int) {
	rw.tab.SetTabItems([]CommonWidgets.TabItem[struct{}]{
		{
			Text: "Body",
		},
		{
			Text: "Header",
		},
	})
	rw.tab.SelectTabItemByIndex(index)
}

func (rw *response_widget) SelectedTab() int {
	return rw.tab.GetSelectedIndex()
}

func (rw *response_widget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	adder.AddWidget(&rw.header_widget)

	u := widget.UnitSize(ctx)
	rw.tab.SetTabItems([]CommonWidgets.TabItem[struct{}]{
		{
			Text: "Body",
		},
		{
			Text: "Header",
		},
	})

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

		switch rw.tab.GetSelectedIndex() {
		case 0:
			rw.tab_content.response_body.SetType(CommonWidgets.HTTP_Response)
			rw.tab_content.selected_content = &rw.tab_content.response_body
		case 1:
			rw.tab_content.selected_content = &rw.tab_content.response_header
		default:
			panic("Unknown tab selected")
		}

		adder.AddWidget(&rw.tab)
		adder.AddWidget(rw.tab_content.selected_content)
	}

	return nil
}

func (rw *response_widget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	u := widget.UnitSize(ctx)
	main_layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Gap:       u / 4,
		Items: []gui.LinearLayoutItem{
			{
				Widget: &rw.header_widget,
			},
			{
				Widget: &rw.tab,
			},
			{
				Widget: rw.tab_content.selected_content,
				Size:   gui.FlexibleSize(1),
			},
		},
	}

	main_layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}
