package http_widget

import (
	messages "API-Client/massages"
	"API-Client/widgets/request/def"
	"image"
	"net/url"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/hajimehoshi/ebiten/v2"
	// "image"
)

type HTTP_Widget struct {
	gui.DefaultWidget

	request_widget  request_widget
	response_widget response_widget

	popup_widget *widget.Popup
	popup_size   *image.Point

	req *def.Request
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
}

func (brp *HTTP_Widget) Update() {
	
}

func (brp *HTTP_Widget) handle_popup() {
	if brp.popup_widget == nil {
		return
	}

	brp.request_widget.input_bar_widget.OnOpenIn(func(ctx *gui.Context) {
		url_panel := get_url_panel(ctx)
		url_panel.content.query.Empty()
		*brp.popup_size = url_panel.Measure(ctx, gui.Constraints{})
		u, err := url.Parse(brp.request_widget.input_bar_widget.input_widget.Value())
		if err != nil {
			messages.Alerts.Push(err.Error())
		}
		url_panel.SetURL(u, ctx)
		brp.popup_widget.SetContent(url_panel)
		brp.popup_widget.SetOpen(true)
	})

	brp.popup_widget.OnClose(func(ctx *gui.Context, reason widget.PopupCloseReason) {
		u1 := url_panel.content.generate_url()
		brp.request_widget.input_bar_widget.input_widget.SetValue(u1.String())

		u2, err := url.Parse(brp.request_widget.url_preview.URL())
		if err != nil {
			messages.Alerts.Push(err.Error())
		}
		u1.RawQuery = u2.RawQuery

		brp.request_widget.url_preview.SetURL(u1.String())
		brp.popup_widget.OnClose(func(_ *gui.Context, _ widget.PopupCloseReason) {})
	})
}

func (brp *HTTP_Widget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	ctx.SetColorMode(ebiten.ColorModeDark)

	brp.handle_popup()

	adder.AddWidget(&brp.request_widget)
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
				Widget: &brp.response_widget,
				Size:   gui.FlexibleSize(1),
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}
