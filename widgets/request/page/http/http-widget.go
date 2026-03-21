package http_widget

import (
	"API-Client/widgets/request/def"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/hajimehoshi/ebiten/v2"
	// "image"
)

type HTTP_Widget struct {
	gui.DefaultWidget
	//https://api.github.com/repos/udan-jayanith/Zbolt
	request_widget  request_widget
	response_widget response_widget
	
	popup_content *gui.Widget 
	popup_widget *widget.Popup
}

func (brp *HTTP_Widget) RequestType() def.RequestType {
	return def.HTTP
}

func (brp *HTTP_Widget) Popup(popup_content *gui.Widget, popup_widget *widget.Popup) {
	brp.popup_content = popup_content
	brp.popup_widget = popup_widget
}

func (brp *HTTP_Widget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	ctx.SetColorMode(ebiten.ColorModeDark)

	brp.request_widget.input_bar_widget.OnOpenIn(func(ctx *gui.Context) {
		if brp.popup_widget == nil {
			return
		}
		*brp.popup_content = get_url_panel(ctx)
		brp.popup_widget.SetOpen(true)
	})
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
