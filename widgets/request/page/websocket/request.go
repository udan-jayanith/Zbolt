package websocket_widget

import (
	gui "github.com/guigui-gui/guigui"
)

type request_widget struct {
	gui.DefaultWidget
}

func (rw *request_widget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	return nil
}

func (rw *request_widget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
}
