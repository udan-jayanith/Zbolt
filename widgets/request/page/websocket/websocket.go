package websocket_widget

import (
	"API-Client/basic"
	"API-Client/widgets/request/def"
	"image"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type WebsocketWidget struct {
	gui.DefaultWidget

	request_widget  request_widget
	response_widget response_widget
}

func (ww *WebsocketWidget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	adder.AddWidget(&ww.request_widget)
	adder.AddWidget(&ww.response_widget)
	return nil
}

func (ww *WebsocketWidget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Gap:       basic.Gap(ctx),
		Items: []gui.LinearLayoutItem{
			{
				Widget: &ww.request_widget,
				Size:   gui.FlexibleSize(1),
			},
			{
				Widget: &ww.response_widget,
				Size:   gui.FlexibleSize(1),
			},
		},
	}

	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (ww *WebsocketWidget) RequestType() def.RequestType {
	return def.Websocket
}

func (ww *WebsocketWidget) SetPopupWidget(w *widget.Popup, popup_size *image.Point) {
}

func (ww *WebsocketWidget) SetReq(req *def.Request)