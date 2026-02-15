package Requester

import (
	"github.com/sqweek/dialog"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type response_body_widget struct {
	gui.DefaultWidget
	not_first_build bool
	file_type       widget.Text

	options struct {
		auto_wrap struct {
			text   widget.Text
			toggle widget.Toggle
		}
		format struct {
			text   widget.Text
			toggle widget.Toggle
		}
		open_with widget.Button
	}
	view widget.TextInput
}

func (rbw *response_body_widget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	{
		rbw.file_type.SetValue("Json")
		rbw.file_type.SetVerticalAlign(widget.VerticalAlignMiddle)
		adder.AddChild(&rbw.file_type)
	}
	{
		rbw.options.auto_wrap.text.SetValue("Auto wrap")
		rbw.options.auto_wrap.text.SetVerticalAlign(widget.VerticalAlignMiddle)
		adder.AddChild(&rbw.options.auto_wrap.text)

		if !rbw.not_first_build {
			rbw.options.auto_wrap.toggle.SetValue(true)
		}
		adder.AddChild(&rbw.options.auto_wrap.toggle)
	}
	{
		rbw.options.format.text.SetValue("Format")
		rbw.options.format.text.SetVerticalAlign(widget.VerticalAlignMiddle)
		adder.AddChild(&rbw.options.format.text)

		if !rbw.not_first_build {
			rbw.options.format.toggle.SetValue(true)
		}
		adder.AddChild(&rbw.options.format.toggle)
	}
	{
		rbw.options.open_with.SetText("Open")
		rbw.options.open_with.SetOnUp(func(context *gui.Context) {
			path, err := dialog.File().Load()
			if err != nil {
				println(err.Error())
				return
			}
			println(path)
		})
		adder.AddChild(&rbw.options.open_with)
	}
	{
		rbw.view.SetAutoWrap(true)
		rbw.view.SetMultiline(true)
		rbw.view.SetEditable(false)
		rbw.view.SetValue(`
		git clone https://github.com/guigui-gui/guigui.git
		cd guigui
		go run ./example/gallery


		hi


		Hello world
		`)
		adder.AddChild(&rbw.view)
	}

	rbw.not_first_build = true
	return nil
}

func (rbw *response_body_widget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	u := widget.UnitSize(ctx)
	toggle_size := gui.FixedSize(u*2 - u/3)
	space := gui.LinearLayoutItem{
		Size: gui.FixedSize(widget.UnitSize(ctx) / 4),
	}

	header_layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Gap:       u / 4,
		Items: []gui.LinearLayoutItem{
			{
				Widget: &rbw.options.auto_wrap.text,
			},
			{
				Widget: &rbw.options.auto_wrap.toggle,
				Size:   toggle_size,
			},
			space,
			{
				Widget: &rbw.options.format.text,
			},
			{
				Widget: &rbw.options.format.toggle,
				Size:   toggle_size,
			},
			{
				Size: gui.FlexibleSize(1),
			},
			{
				Widget: &rbw.file_type,
			},
			{
				Widget: &rbw.options.open_with,
			},
		},
	}

	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Gap:       u / 4,
		Items: []gui.LinearLayoutItem{
			space,
			{
				Layout: header_layout,
			},
			{
				Widget: &rbw.view,
				Size:   gui.FlexibleSize(1),
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}