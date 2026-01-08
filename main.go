package main

import (
	"fmt"
	"os"

	"API-Client/basic"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/hajimehoshi/ebiten/v2"
)

type RecentList struct {
}

type WelcomePage struct {
	gui.DefaultWidget

	background      widget.Background
	create_a_text   widget.Text
	open            widget.Button
	new_request     widget.Button
	create_project  widget.Button
	recent_projects widget.Text
}

func (wp *WelcomePage) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	ctx.SetColorMode(gui.ColorModeDark)
	adder.AddChild(&wp.background)

	wp.create_a_text.SetValue("Create a")
	wp.create_a_text.SetBold(true)
	wp.create_a_text.SetScale(basic.HeadingSize)
	adder.AddChild(&wp.create_a_text)

	wp.open.SetText("Open project")
	wp.open.SetTextBold(true)
	adder.AddChild(&wp.open)

	wp.new_request.SetText("New Request")
	wp.new_request.SetTextBold(true)
	adder.AddChild(&wp.new_request)

	wp.create_project.SetText("Create Project")
	wp.create_project.SetTextBold(true)
	adder.AddChild(&wp.create_project)

	wp.recent_projects.SetValue("Recent projects")
	wp.recent_projects.SetBold(true)
	adder.AddChild(&wp.recent_projects)

	return nil
}

func (wp *WelcomePage) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	bounds := widgetBounds.Bounds()
	layouter.LayoutWidget(&wp.background, bounds)

	u := widget.UnitSize(ctx)
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Items: []gui.LinearLayoutItem{
			{
				Widget: &wp.create_a_text,
			},
			{
				Size: gui.FixedSize(u / 2),
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
		},
	}

	layout = basic.Align(layout.Items, basic.Center, basic.Start)
	layout.Padding = basic.NewPadding(u*2, u, u, u)
	layout.LayoutWidgets(ctx, bounds, layouter)
}

func main() {
	op := &gui.RunOptions{
		Title: "API Client",
		RunGameOptions: &ebiten.RunGameOptions{
			ApplePressAndHoldEnabled: true,
		},
	}
	if err := gui.Run(&WelcomePage{}, op); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
