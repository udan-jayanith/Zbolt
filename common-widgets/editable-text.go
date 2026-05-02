package CommonWidgets

import (
	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type EditableText struct {
	not_editable bool
	widget.Text
}

func (et *EditableText) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	et.Text.SetEditable(!et.not_editable)
	et.Text.Build(ctx, adder)
	return nil
}

func (et *EditableText) SetEditable(editable bool) {
	et.not_editable = !editable
}
