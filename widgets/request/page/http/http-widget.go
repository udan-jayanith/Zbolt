package http_widget

import (
	CommonWidgets "API-Client/common-widgets"
	messages "API-Client/massages"
	"API-Client/widgets/request/def"
	url_utils "API-Client/widgets/request/url-utils"
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

	url_panel_widget url_panel_widget
	popup_widget     widget.Popup

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

	// Setup request widget
	brp.data = data
	brp.request_widget.SetHeaders(data.Headers)
	brp.request_widget.SetParameters(data.Parameters)
	brp.request_widget.SetAutowrap(data.RequestConfig.AutoWrap)
	brp.request_widget.SetFormat(data.RequestConfig.Formate)
	brp.request_widget.SetContentType(data.Body.ContentType)
	brp.request_widget.SetBody(&data.Body)
	if data.Method == "" {
		data.Method = "Get"
	}
	brp.request_widget.SetMethod(data.Method)
	brp.request_widget.SelectTab(data.SelectedRequestTab())

	u, err := url.Parse(data.URL.BaseURL)
	if err != nil {
		messages.Alerts.Push(err.Error())
	}
	u.Path = data.URL.EncodedPath()
	brp.request_widget.SetURL(u)
	brp.request_widget.DisableURLInput(data.URL.IsPattern())

	// Setup response widget
	data.ResponseData(func(res_data *def.HTTP_Response_Data) {
		brp.response_widget.SetHeaders(res_data.Headers)

		brp.response_widget.SetAutowrap(data.ResponseConfig.AutoWrap)
		brp.response_widget.SetFormat(data.ResponseConfig.Formate)
		brp.response_widget.SetResponseBody(&res_data.Body)
		brp.response_widget.SetSelectedTab(res_data.SelectedResponseTab)

		brp.response_widget.SetHTTPVersion(res_data.Version)
		brp.response_widget.SetResponseTime(res_data.ResponseTime)
		if res_data.Status_code != 0 {
			brp.response_widget.SetStatus(res_data.Status_code)
		}
	})
	gui.RequestRebuild(brp)
}

// TODO: SyncData should be run to save data before switching tabs, closing tabs or closing the app.
func (brp *HTTP_Widget) SyncData() {
	brp.data.Parameters = brp.request_widget.Parameters()
	brp.data.Headers = brp.request_widget.Headers()
	brp.data.Body.ContentType = brp.request_widget.ContentType()
	brp.data.Body.Content = brp.request_widget.Body()

	brp.data.SetSelectedRequestTab(brp.request_widget.SelectedTab())
	brp.data.ResponseData(func(value *def.HTTP_Response_Data) {
		value.SelectedResponseTab = brp.response_widget.SelectedTab()
	})
	// TODO: HTTP response data is synced in when request is finished
}

func (brp *HTTP_Widget) url_panel_popup_size(ctx *gui.Context, widgetBounds *gui.WidgetBounds) image.Rectangle {
	url_measurements := brp.url_panel_widget.Measure(ctx, gui.Constraints{})
	b := widgetBounds.Bounds()

	b.Min.X += (b.Dx() / 2) - (url_measurements.X / 2)
	b.Min.Y += (b.Dy() / 2) - (url_measurements.Y / 2)

	b.Max.X = b.Min.X + url_measurements.X
	b.Max.Y = b.Min.Y + url_measurements.Y

	return b
}

func (brp *HTTP_Widget) on_url_panel_open(ctx *gui.Context) {
	u, _ := url.Parse(brp.data.URL.BaseURL)
	brp.url_panel_widget.Set(u.Scheme, u.Host, brp.data.URL.RawPath(), brp.data.URL.Path.Pattern.Attributes)
	brp.popup_widget.SetOpen(true)
}

func (brp *HTTP_Widget) on_url_panel_close(ctx *gui.Context, reason widget.PopupCloseReason) {
	u, err := url.Parse(brp.url_panel_widget.URL())
	if err != nil {
		messages.Alerts.Push(err.Error())
	}
	brp.request_widget.SetURL(u)

	pattern, query_list := brp.url_panel_widget.Pattern()
	if len(query_list) > 0 {
		brp.data.URL.SetPattern(pattern, query_list)
		brp.request_widget.DisableURLInput(true)
	} else {
		brp.data.URL.SetPath(u.Path)
		brp.request_widget.DisableURLInput(false)
	}
	u.Path = ""
	url_utils.CleanURL(u)
	brp.data.URL.BaseURL = u.String()

	brp.url_panel_widget.Clear()
}

func (brp *HTTP_Widget) on_request_button_clicked(ctx *gui.Context, value string) {

}

func (brp *HTTP_Widget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	ctx.SetPreferredColorMode(ebiten.ColorModeDark)

	brp.request_widget.OnOpenIn(brp.on_url_panel_open)
	adder.AddWidget(&brp.popup_widget)
	if brp.popup_widget.IsOpen() {
		brp.popup_widget.SetAnimated(true)
		brp.popup_widget.SetBackgroundDark(true)
		brp.popup_widget.SetCloseByClickingOutside(true)
		brp.popup_widget.SetContent(&brp.url_panel_widget)
		brp.popup_widget.OnClose(brp.on_url_panel_close)
	}

	brp.request_widget.OnMethodChanged(func(method string) {
		brp.data.Method = method
	})

	brp.request_widget.OnURLInputChanged(func(context *gui.Context, text string, committed bool) {
		if !committed || brp.data.URL.IsPattern() {
			return
		}
		u, err := url.Parse(text)
		if err != nil {
			messages.Alerts.Push(err.Error())
		}
		brp.request_widget.SetURL(u)

		url_utils.CleanURL(u)
		brp.data.URL.SetPath(u.Path)
		u.Path = ""
		brp.data.URL.BaseURL = u.String()
	})

	brp.request_widget.OnRequestButtonClicked(brp.on_request_button_clicked)

	brp.request_widget.OnAutowrap(func(ctx *gui.Context, value bool) {
		brp.data.RequestConfig.AutoWrap = value
	})
	brp.request_widget.OnFormat(func(ctx *gui.Context, value bool) {
		brp.data.RequestConfig.Formate = value
	})

	brp.response_widget.OnAutowrapToggle(func(ctx *gui.Context, value bool) {
		brp.data.ResponseConfig.AutoWrap = value
	})
	brp.response_widget.OnFormatToggle(func(ctx *gui.Context, value bool) {
		brp.data.ResponseConfig.Formate = value
	})

	adder.AddWidget(&brp.request_widget)
	adder.AddWidget(&brp.vr)
	adder.AddWidget(&brp.response_widget)
	return nil
}

func (brp *HTTP_Widget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	if brp.popup_widget.IsOpen() {
		brp.popup_widget.SetBackgroundBounds(widgetBounds.Bounds())
		layouter.LayoutWidget(&brp.popup_widget, brp.url_panel_popup_size(ctx, widgetBounds))
	}

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
