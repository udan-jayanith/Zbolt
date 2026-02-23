package request_page

import (
	"API-Client/icons"

	"image"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type NothingWidget struct {
	gui.DefaultWidget

	image  icons.Icon
}

func (nw *NothingWidget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	nw.image.IconName = "large-icons/add-box"
	u := widget.UnitSize(ctx)
	size := u * 8
	point := image.Pt(size, size)
	nw.image.Point = &point
	nw.image.OnClick(func() {})

	adder.AddChild(&nw.image)
	return nil
}

func (nw *NothingWidget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	flex1_item := gui.LinearLayoutItem{
		Size: gui.FlexibleSize(1),
	}
	
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Items: []gui.LinearLayoutItem{
			flex1_item,
			{
				Layout: gui.LinearLayout{
					Direction: gui.LayoutDirectionVertical,
					Items: []gui.LinearLayoutItem{
						flex1_item,
						{
							Widget: &nw.image,
						},
						flex1_item,
					},
				},
			},
			flex1_item,
		},
	}
	
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}
