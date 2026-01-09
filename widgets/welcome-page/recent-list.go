package WelcomePage

import (
	"image"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type RecentList struct {
	gui.DefaultWidget

	recent_projects []*widget.Text
}

func (rl *RecentList) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	for _, recent_project := range rl.recent_projects {
		adder.AddChild(recent_project)
	}
	return nil
}

func (rl *RecentList) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	u := widget.UnitSize(ctx)
	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Gap:       u / 4,
		Items:     make([]gui.LinearLayoutItem, len(rl.recent_projects)),
	}
	for i, recent_project := range rl.recent_projects {
		layout.Items[i] = gui.LinearLayoutItem{
			Widget: recent_project,
			Size:   gui.FlexibleSize(1),
		}
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (rl *RecentList) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	u := widget.UnitSize(ctx)
	h := u * 4
	l := len(rl.recent_projects)
	if l > 0 {
		project := rl.recent_projects[0]
		points := project.Measure(ctx, constraints)
		h = points.Y * l
	}
	if w, ok := constraints.FixedWidth(); ok {
		return image.Pt(w, h+u)
	} else if h, ok := constraints.FixedHeight(); ok {
		return image.Pt(u*4, h)
	}
	return image.Pt(u*4, u*2)
}

func (rl *RecentList) Clear(){
	rl.recent_projects = make([]*widget.Text, 0, 4)
}

func (rl *RecentList) Add(text string) {
	text_widget := &widget.Text{}
	text_widget.SetValue(text)
	rl.recent_projects = append(rl.recent_projects, text_widget)
}
