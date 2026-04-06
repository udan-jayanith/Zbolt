package http_widget

import (
	CommonWidgets "API-Client/common-widgets"
	messages "API-Client/massages"
	attr "API-Client/widgets/request"
	"API-Client/widgets/request/def"
	"image"
	"net/url"
	"time"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/hajimehoshi/ebiten/v2"
	// "image"
)

type HTTP_Widget struct {
	gui.DefaultWidget

	request_widget  request_widget
	vr              CommonWidgets.VerticalLine
	response_widget response_widget

	popup_widget *widget.Popup
	popup_size   *image.Point

	req  *def.Request
	data *def.HTTP_Data
	t    time.Time
}

func (brp *HTTP_Widget) RequestType() def.RequestType {
	return def.HTTP
}

func (brp *HTTP_Widget) SetPopupWidget(w *widget.Popup, popup_size *image.Point) {
	brp.popup_widget = w
	brp.popup_size = popup_size
}

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

	brp.response_widget.SetAutowrap(data.ResponseConfig.AutoWrap)
	brp.response_widget.SetFormat(data.ResponseConfig.Formate)

	temp := data.ResponseData()
	brp.response_widget.SetHeaders(temp.Headers)
	brp.response_widget.SetResponseBody(&temp.Body)
}

func (brp *HTTP_Widget) update() {
	d := brp.data
	brp.response_widget.OnAutowrapToggle(func(ctx *gui.Context, value bool) {
		d.ResponseConfig.AutoWrap = value
	})

	brp.response_widget.OnFormatToggle(func(ctx *gui.Context, value bool) {
		d.ResponseConfig.Formate = value
	})

	brp.response_widget.SetResponseData(brp.data.ResponseData())

	if time.Now().Sub(brp.t).Seconds() < 1 {
		return
	}
	brp.t = time.Now()

	//TODO: update url preview every second
	brp.SyncData()
}

// TODO: run SyncData before switching tabs
func (brp *HTTP_Widget) SyncData() {
	brp.data.Headers = brp.request_widget.Headers()
	brp.data.Parameters = brp.request_widget.Parameters()
	brp.data.SetSelectedRequestTab(brp.request_widget.SelectedTab())
	brp.data.Method = brp.request_widget.Method()
}

func (brp *HTTP_Widget) handle_popup() {
	if brp.popup_widget == nil {
		return
	}

	brp.request_widget.input_bar_widget.OnOpenIn(func(ctx *gui.Context) {
		url_panel := get_url_panel(ctx)
		*brp.popup_size = url_panel.Measure(ctx, gui.Constraints{})
		
		// TODO: only parse the url from the url input if there is no pattern
		// otherwise use the pattern 
		u, err := url.Parse(brp.request_widget.input_bar_widget.input_widget.Value()) // Gets the url from the url input bar 
		if err != nil {
			messages.Alerts.Push(err.Error())
		}
		
		url_panel.SetURL(u, ctx)
		brp.popup_widget.SetContent(url_panel)
		brp.popup_widget.SetOpen(true)
	})

	brp.popup_widget.OnClose(func(ctx *gui.Context, reason widget.PopupCloseReason) {
		url_panel_content := get_url_panel(ctx).content

		u, is_pattern := url_panel_content.generate_url()
		brp.request_widget.input_bar_widget.input_widget.SetValue(u.String())
		path := u.Path
		u.Path = ""
		brp.data.URL.BaseURL = u.String()
		if is_pattern {
			brp.data.URL.Path.RawPath = ""
			brp.data.URL.Path.Pattern.Pattern = url_panel_content.path.Value()
			brp.data.URL.Path.Pattern.Attributes = url_panel_content.query.Rows()
		} else {
			brp.data.URL.Path.Pattern.Pattern = ""
			brp.data.URL.Path.Pattern.Attributes = []attr.Attribute{}
			brp.data.URL.Path.RawPath = path
		}

		url_panel_content.query.SetRows([]attr.Attribute{})
		brp.popup_widget.OnClose(func(_ *gui.Context, _ widget.PopupCloseReason) {})
	})
}

func (brp *HTTP_Widget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	ctx.SetPreferredColorMode(ebiten.ColorModeDark)

	brp.update()
	brp.handle_popup()

	adder.AddWidget(&brp.request_widget)
	adder.AddWidget(&brp.vr)
	adder.AddWidget(&brp.response_widget)
	return nil
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
