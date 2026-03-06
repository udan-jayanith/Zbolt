package websocket_widget

import (
	gui "github.com/guigui-gui/guigui"
)

type request_body struct {
	gui.DefaultWidget
}

func (ww *request_body) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	return nil
}

func (ww *request_body) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
}
