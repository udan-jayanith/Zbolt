package http_widget

import (
	CommonWidgets "API-Client/common-widgets"
	"API-Client/widgets/request/def"
	"fmt"
	"net/http"
	"strconv"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type response_widget struct {
	gui.DefaultWidget
	header_widget response_header_widget
	tab           CommonWidgets.Tab[struct{}]
	tab_content   struct {
		response_header  widget.Table[string]
		response_body    response_body_widget
		selected_content gui.Widget
	}
}

func (rw *response_widget) SetAutowrap(autowrap bool) {

}

func (rw *response_widget) SetFormat(format bool) {

}

func (rw *response_widget) SetResponseData(res_data *def.TempHTTP_Data) {
	res_status := fmt.Sprintf("%v %s", res_data.Status_code, http.StatusText(res_data.Status_code))
	rw.header_widget.status.SetValue(res_status)

	rw.header_widget.response_time.SetValue(res_data.ResponseTime.String())
	rw.header_widget.size.SetValue(strconv.Itoa(res_data.ResponseSize))
	rw.header_widget.proto.SetValue(fmt.Sprintf("HTTP v%v.%v", res_data.Version.Major, res_data.Version.Minor))
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

		switch rw.tab.GetSelectedIndex() {
		case 0:
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
