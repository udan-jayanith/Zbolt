package Welcome

import (
	"API-Client/basic"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/hajimehoshi/ebiten/v2"
)

type Welcome struct {
	gui.DefaultWidget

	background      widget.Background
	create_a_text   widget.Text
	open            widget.Button
	new_request     widget.Button
	create_project  widget.Button
	recent_projects widget.Text
	recent_list     RecentList
}

func (wp *Welcome) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	ctx.SetColorMode(ebiten.ColorModeDark)
	adder.AddWidget(&wp.background)

	wp.create_a_text.SetValue("Create a")
	wp.create_a_text.SetBold(true)
	wp.create_a_text.SetScale(basic.HeadingSize)
	adder.AddWidget(&wp.create_a_text)

	wp.open.SetText("Open project")
	wp.open.SetTextBold(true)
	adder.AddWidget(&wp.open)

	wp.new_request.SetText("New Request")
	wp.new_request.SetTextBold(true)
	adder.AddWidget(&wp.new_request)

	wp.create_project.SetText("Create Project")
	wp.create_project.SetTextBold(true)
	adder.AddWidget(&wp.create_project)

	wp.recent_projects.SetValue("Recent projects")
	wp.recent_projects.SetBold(true)
	wp.recent_projects.SetScale(basic.HeadingSize)
	adder.AddWidget(&wp.recent_projects)

	wp.recent_list.Add([]*RecentItem{
		{
			Text: basic.NewTextWidget(`C:\Users\Udan\Documents\Dev\oss-contributions\guigui\example\todo`),
		},
		{
			Text: basic.NewTextWidget(`C:\Users\Udan\Documents\DIIT98\Internet and Email\00307060`),
		},
		{
			Text: basic.NewTextWidget(`C:\Users\Udan\Pictures\Camera Roll`),
		},
	})
	adder.AddWidget(&wp.recent_list)
	return nil
}

func (wp *Welcome) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	bounds := widgetBounds.Bounds()
	layouter.LayoutWidget(&wp.background, bounds)

	u := widget.UnitSize(ctx)
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Gap:       u / 4,
		Padding:   basic.NewPadding(u*2, u, u, u),
		Items: []gui.LinearLayoutItem{
			{
				Widget: &wp.create_a_text,
			},
			{
				Layout: gui.LinearLayout{
					Direction: gui.LayoutDirectionHorizontal,
					Gap:       u,
					Items: []gui.LinearLayoutItem{
						{
							Widget: &wp.open,
						},
						{
							Widget: &wp.new_request,
						},
						{
							Widget: &wp.create_project,
						},
					},
				},
			},
			{
				Size: gui.FixedSize(u * 2),
			},
			{
				Widget: &wp.recent_projects,
			},
			{
				Widget: &wp.recent_list,
			},
		},
	}

	layout = basic.Align(layout, basic.Center, basic.Start)
	layout.LayoutWidgets(ctx, bounds, layouter)
}
