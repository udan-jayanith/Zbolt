package inspect

import (
	"github.com/udan-jayanith/Zbolt/basic"
	CommonWidgets "github.com/udan-jayanith/Zbolt/common-widgets"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/guigui-gui/guigui/basicwidget/basicwidgetdraw"
	"github.com/hajimehoshi/ebiten/v2"
)

type InspectWidget struct {
	gui.DefaultWidget
	open         bool
	tabs         CommonWidgets.Tab[string]
	logs, stdout gui.WidgetWithPadding[*widget.Text]
	panel        widget.Panel
}

func (r *InspectWidget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	if !r.open {
		return nil
	}

	r.tabs.SetTabItems([]CommonWidgets.TabItem[string]{
		CommonWidgets.TabItem[string]{
			Text: "Logs",
		},
		CommonWidgets.TabItem[string]{
			Text: "Stdout",
		},
		CommonWidgets.TabItem[string]{
			Text: "Database",
		},
	})
	adder.AddWidget(&r.tabs)

	padding_end := gui.Padding{End: widget.UnitSize(ctx) / 2}
	r.logs.SetPadding(padding_end)
	logs_widget := r.logs.Widget()
	logs_widget.SetAutoWrap(true)
	logs_widget.SetMultiline(true)
	logs_widget.SetSelectable(true)

	r.stdout.SetPadding(padding_end)
	stdout_widget := r.stdout.Widget()
	stdout_widget.SetAutoWrap(true)
	stdout_widget.SetMultiline(true)
	stdout_widget.SetSelectable(true)

	selected_index := r.tabs.GetSelectedIndex()
	switch selected_index {
	case 0:
		r.panel.SetContent(&r.logs)
	case 1:
		r.panel.SetContent(&r.stdout)
	case 2:
		panic("Not implemented")
	default:
		panic("Not handled")
	}
	logs_widget.SetValue(`What is Lorem Ipsum?
Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.

Why do we use it?
It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout. The point of using Lorem Ipsum is that it has a more-or-less normal distribution of letters, as opposed to using 'Content here, content here', making it look like readable English. Many desktop publishing packages and web page editors now use Lorem Ipsum as their default model text, and a search for 'lorem ipsum' will uncover many web sites still in their infancy. Various versions have evolved over the years, sometimes by accident, sometimes on purpose (injected humour and the like).

`)

	r.panel.SetBorders(widget.PanelBorders{
		Top: true,
	})
	r.panel.SetContentConstraints(widget.PanelContentConstraintsFixedWidth)
	adder.AddWidget(&r.panel)
	return nil
}

func (r *InspectWidget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	if !r.open {
		return
	}

	u := widget.UnitSize(ctx)
	padding := u / 4

	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Padding:   basic.NewPadding(padding),
		Gap:       u / 4,
		Items: []gui.LinearLayoutItem{
			{
				Widget: &r.tabs,
			},
			{
				Widget: &r.panel,
				Size:   gui.FlexibleSize(1),
			},
		},
	}

	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (r *InspectWidget) Draw(ctx *gui.Context, widgetBounds *gui.WidgetBounds, dst *ebiten.Image) {
	cm := ctx.ColorMode()
	b := widgetBounds.Bounds()

	clr1, clr2 := basicwidgetdraw.BorderColors(cm, basicwidgetdraw.RoundedRectBorderTypeRegular)
	basicwidgetdraw.DrawRoundedRectBorder(ctx, dst, b, clr1, clr2, 1, 1, basicwidgetdraw.RoundedRectBorderTypeRegular)

	background_color := basicwidgetdraw.BackgroundSecondaryColor(cm)
	basicwidgetdraw.DrawRoundedRect(ctx, dst, b, background_color, 1)
}

func (r *InspectWidget) SetOpen(open bool) {
	r.open = open
}

func (r *InspectWidget) IsOpen() bool {
	return r.open
}
