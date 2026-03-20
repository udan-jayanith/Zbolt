package http_widget

import (
	"API-Client/basic"
	CommonWidgets "API-Client/common-widgets"
	"image"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/hajimehoshi/ebiten/v2"
)

var url_panel *url_panel_widget_scrollable = func() *url_panel_widget_scrollable {
	w := &url_panel_widget{}
	w.host.SetValue("api.github.com")
	w.path.SetValue("/repos/{user-name}/{repo-name}")
	return &url_panel_widget_scrollable{
		content: w,
	}
}()

type long_text_input_widget struct {
	widget.TextInput
	t widget.TextInput
}

func (w *long_text_input_widget) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	point := w.t.Measure(ctx, gui.Constraints{})
	point.X *= 3
	return point
}

type url_panel_widget struct {
	gui.DefaultWidget

	form                              widget.Form
	scheme_text, host_text, path_text widget.Text
	scheme, host, path                long_text_input_widget

	query_header, pattern_header           widget.Text
	query_description, pattern_description CommonWidgets.Description
	query, pattern                         CommonWidgets.AttributeTable

	hr1, hr2         CommonWidgets.HorizontalLine
	pseudo_url CommonWidgets.Description
	
	url_preview_header widget.Text
	url_preview CommonWidgets.URLPreview
}

func (w *url_panel_widget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	w.scheme.SetValue("http")
	ctx.SetEnabled(&w.scheme, false)

	w.scheme_text.SetValue("Scheme")
	w.host_text.SetValue("Host")
	w.path_text.SetValue("Path")

	w.form.SetItems([]widget.FormItem{
		{
			PrimaryWidget:   &w.scheme_text,
			SecondaryWidget: &w.scheme,
		},
		{
			PrimaryWidget:   &w.host_text,
			SecondaryWidget: &w.host,
		},
		{
			PrimaryWidget:   &w.path_text,
			SecondaryWidget: &w.path,
		},
	})
	adder.AddWidget(&w.form)

	w.query_header.SetValue("Query")
	adder.AddWidget(&w.query_header)

	w.query_description.SetDescription("Attributes after the url path followed by ? mark.")
	adder.AddWidget(&w.query_description)

	w.pattern_header.SetValue("Patterns")
	adder.AddWidget(&w.pattern_header)

	w.pattern_description.SetDescription("Attributes in the url path enclosed by {}.")
	adder.AddWidget(&w.pattern_description)

	adder.AddWidget(&w.query)
	adder.AddWidget(&w.pattern)

	w.pseudo_url.SetDescription("The general form represented is:\n``[scheme:][//[userinfo@]host][/]path[?query][#fragment]``")
	adder.AddWidget(&w.pseudo_url)
	adder.AddWidget(&w.hr1)
	
	adder.AddWidget(&w.hr2)
	
	w.url_preview_header.SetValue("URL preview")
	adder.AddWidget(&w.url_preview_header)
	
	w.url_preview.SetURL("https://api.github.com/repos/udan-jayanith/Zbolt")
	adder.AddWidget(&w.url_preview)
	return nil
}

func (w *url_panel_widget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	u := widget.UnitSize(ctx)
	gap := gui.LinearLayoutItem{
		Size: gui.FixedSize(u / 3),
	}

	layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Padding:   basic.NewPadding(u / 2),
		Items: []gui.LinearLayoutItem{
			{
				Widget: &w.form,
			},
			gap,
			{
				Widget: &w.query_header,
			},
			{
				Widget: &w.query_description,
			},
			{
				Widget: &w.query,
				Size:   gui.FlexibleSize(1),
			},
			gap,
			{
				Widget: &w.pattern_header,
			},
			{
				Widget: &w.pattern_description,
			},
			{
				Widget: &w.pattern,
				Size:   gui.FlexibleSize(1),
			},
			gap,
			{
				Widget: &w.hr1,
			},
			gap,
			{
				Widget: &w.pseudo_url,
			},
			gap,
			{
				Widget: &w.hr2,
			},
			gap,
			{
				Widget: &w.url_preview_header,
			},
			gap,
			{
				Widget: &w.url_preview,
			},
		},
	}
	layout.LayoutWidgets(ctx, widgetBounds.Bounds(), layouter)
}

func (w *url_panel_widget) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	var point image.Point
	u := widget.UnitSize(ctx)
	point.X = u * 36
	point.Y = u * 30

	x, y := ebiten.WindowSize()
	if x < point.X && y < point.Y {
		point.X = u * 28
	}
	return point
}

type url_panel_widget_scrollable struct {
	gui.DefaultWidget

	content *url_panel_widget
	panel   widget.Panel
}

func (w *url_panel_widget_scrollable) Build(context *gui.Context, adder *gui.ChildAdder) error {
	w.panel.SetContentConstraints(widget.PanelContentConstraintsFixedWidth)
	w.panel.SetContent(w.content)
	adder.AddWidget(&w.panel)
	return nil
}

func (w *url_panel_widget_scrollable) Layout(context *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layouter.LayoutWidget(&w.panel, widgetBounds.Bounds())
}

func (w *url_panel_widget_scrollable) Measure(ctx *gui.Context, constraints gui.Constraints) image.Point {
	var point image.Point
	u := widget.UnitSize(ctx)
	point.X = u * 36
	point.Y = u * 28

	x, y := ebiten.WindowSize()
	if x < point.X && y < point.Y {
		point.X = u * 28
		point.Y = u*18
	}
	return point
}
