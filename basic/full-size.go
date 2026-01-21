package basic

import(
	gui "github.com/guigui-gui/guigui"
)

type FullSizeWidget struct {
	gui.DefaultWidget
	widget gui.Widget
}

func (fsw *FullSizeWidget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	adder.AddChild(fsw.widget)
	return nil
}

func (fsw *FullSizeWidget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layout := gui.LinearLayout {
		Items: []gui.LinearLayoutItem{
			{
				Widget: fsw.widget,
				Size: gui.FlexibleSize(1),
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func NewFullSizeWidget(widget gui.Widget) *FullSizeWidget {
	return &FullSizeWidget{
		widget: widget,
	}
}