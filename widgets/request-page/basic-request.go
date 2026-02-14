package Requester

import (
	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type Requester struct {
	gui.DefaultWidget
	background widget.Background
	panel      struct {
		request struct {
			panel   widget.Panel
			content gui.WidgetWithSize[*RequestWidget]
		}
		response struct {
			panel   widget.Panel
			content gui.WidgetWithSize[*ResponseWidget]
		}
	}
}

func (brp *Requester) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	ctx.SetColorMode(gui.ColorModeDark)
	adder.AddChild(&brp.background)

	brp.panel.request.panel.SetContent(&brp.panel.request.content)
	brp.panel.request.panel.SetBorders(widget.PanelBorders{
		End: true,
	})
	adder.AddChild(&brp.panel.request.panel)

	brp.panel.response.panel.SetContent(&brp.panel.response.content)
	adder.AddChild(&brp.panel.response.panel)
	return nil
}

func (brp *Requester) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layouter.LayoutWidget(&brp.background, widgetBounds.Bounds())
	b := widgetBounds.Bounds()

	panel_size := b.Max
	panel_size.X = panel_size.X / 2
	brp.panel.request.content.SetFixedSize(panel_size)
	brp.panel.response.content.SetFixedSize(panel_size)

	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Items: []gui.LinearLayoutItem{
			{
				Widget: &brp.panel.request.panel,
				Size:   gui.FlexibleSize(1),
			},
			{
				Widget: &brp.panel.response.panel,
				Size:   gui.FlexibleSize(1),
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}
