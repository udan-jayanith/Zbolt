package basic

import (
	"image"

	widget "github.com/guigui-gui/guigui/basicwidget"
)

type TooltipHelper struct {
	Widget widget.TooltipArea
	IsOpen bool
	Bounds image.Rectangle
}

func (t *TooltipHelper) Open(is_open bool, text string, bounds image.Rectangle) {
	t.IsOpen = is_open
	t.Widget.SetText(text)
	t.Bounds = bounds
}
