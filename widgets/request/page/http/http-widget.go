package http_widget

import (
	CommonWidgets "API-Client/common-widgets"
	messages "API-Client/massages"
	"API-Client/widgets/request/def"
	"image"
	"net/url"
	"time"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/hajimehoshi/ebiten/v2"
)

type HTTP_Widget struct {
	gui.DefaultWidget

	request_widget  request_widget
	vr              CommonWidgets.VerticalLine
	response_widget response_widget

	// TODO: add a popup and a url_panel_widget widget here.

	req  *def.Request
	data *def.HTTP_Data
	t    time.Time
}

func (brp *HTTP_Widget) RequestType() def.RequestType {
	return def.HTTP
}

// SetReq runs when switching tabs and tab data are passed to this.
func (brp *HTTP_Widget) SetReq(req *def.Request) {
	if req.Type != def.HTTP {
		panic("Invalid request type")
	}
	brp.req = req

	data, ok := req.Data().(*def.HTTP_Data)
	if !ok {
		panic("Invalid data type")
	} else if data == brp.data {
		return
	}

	brp.data = data
	brp.request_widget.SetHeaders(data.Headers)
	brp.request_widget.SetParameters(data.Parameters)
	brp.request_widget.SetTab(data.SelectedRequestTab())
	if data.Method == "" {
		data.Method = "Get"
	}
	brp.request_widget.SetMethod(data.Method)

	u, err := url.Parse(brp.data.URL.BaseURL)
	u.Path = data.URL.GetPath()
	if err != nil {
		messages.Alerts.Push(err.Error())
	}
	brp.request_widget.SetURL_InputValue(u.String())

	brp.response_widget.SetAutowrap(data.ResponseConfig.AutoWrap)
	brp.response_widget.SetFormat(data.ResponseConfig.Formate)

	temp := data.ResponseData()
	brp.response_widget.SetHeaders(temp.Headers)
	brp.response_widget.SetResponseBody(&temp.Body)
}

// TODO: SyncData should be run to save data before switching tabs, closing tabs or closing the app.
func (brp *HTTP_Widget) SyncData() {
	brp.data.Headers = brp.request_widget.Headers()
	brp.data.Parameters = brp.request_widget.Parameters()
	brp.data.SetSelectedRequestTab(brp.request_widget.SelectedTab())
	brp.data.Method = brp.request_widget.Method()
	// TODO: sync the request.
}

func (brp *HTTP_Widget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	ctx.SetPreferredColorMode(ebiten.ColorModeDark)

	// TODO: handle the url panel popup here
	 
	adder.AddWidget(&brp.request_widget)
	adder.AddWidget(&brp.vr)
	adder.AddWidget(&brp.response_widget)
	return nil
}

// TODO: finish this
func (brp *HTTP_Widget) url_panel_popup_size() image.Rectangle {
	return image.Rectangle{}
}

func (brp *HTTP_Widget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Gap:       widget.UnitSize(ctx) / 4,
		Items: []gui.LinearLayoutItem{
			{
				Widget: &brp.request_widget,
				Size:   gui.FlexibleSize(1),
			},
			{
				Widget: &brp.vr,
			},
			{
				Widget: &brp.response_widget,
				Size:   gui.FlexibleSize(1),
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}
