package http_widget

import (
	"API-Client/basic"
	CommonWidgets "API-Client/common-widgets"
	"image"
	"net/url"
	"time"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
	"github.com/hajimehoshi/ebiten/v2"
)

var url_panel *url_panel_widget_scrollable

func get_url_panel(ctx *gui.Context) *url_panel_widget_scrollable {
	if url_panel == nil {
		w := &url_panel_widget{}
		u, _ := url.Parse("")
		w.set_url(u, ctx)
		url_panel = &url_panel_widget_scrollable{
			content: w,
		}
	}
	return url_panel
}

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

	query_header           widget.Text
	query_description CommonWidgets.Description
	query                      query_table_widget

	hr1, hr2   CommonWidgets.HorizontalLine
	pseudo_url CommonWidgets.Description

	url_preview_header widget.Text
	url_preview        CommonWidgets.URLPreview

	t time.Time
}

func (w *url_panel_widget) generate_url() *url.URL {
	q, _ := Parse_url_path_query(w.path.Value())
	values := w.query.GetValues()
	
	if len(q.List) == len(values) {
		for i, v := range values{
			q.List[i].V = v
		}
	}
	
	u := &url.URL{
		Scheme: "http",
		Host:   w.host.Value(),
		Path:   q.Path(),
	}
	return u
}

func (w *url_panel_widget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	w.scheme.SetValue("http")
	ctx.SetEnabled(&w.scheme, false)

	w.scheme_text.SetValue("Scheme")
	w.host_text.SetValue("Host")
	w.path_text.SetValue("Path")

	w.path.OnValueChanged(func(context *gui.Context, text string, committed bool) {
		w.query.Empty()
		q, _ := Parse_url_path_query(text)
		for _, v := range q.List {
			w.query.push_row(string(v.K), string(v.V))
		}
	})
	
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

	w.query_description.SetDescription("Attributes enclosed by {} in path.")
	adder.AddWidget(&w.query_description)

	adder.AddWidget(&w.query)

	w.pseudo_url.SetDescription("The general form of the URL is:\n``[scheme:][//[host][/]path[?query]``")
	adder.AddWidget(&w.pseudo_url)
	adder.AddWidget(&w.hr1)

	adder.AddWidget(&w.hr2)

	w.url_preview_header.SetValue("URL preview")
	adder.AddWidget(&w.url_preview_header)

	if time.Now().Sub(w.t).Seconds() > 1 {
		u := w.generate_url()
		w.url_preview.SetURL(u.String())

		w.t = time.Now()
	}
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
	point.Y = u * 20

	x, y := ebiten.WindowSize()
	if x < point.X && y < point.Y {
		point.X = u * 28
	}
	return point
}

func (w *url_panel_widget) set_url(u *url.URL, ctx *gui.Context) {
	w.host.SetValue(u.Host)
	w.path.SetValue(u.Path)

	q, _ := Parse_url_path_query(u.Path)
	for _, v := range q.List {
		w.query.push_row(string(v.K), string(v.V))
	}
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
	point := w.content.Measure(ctx, gui.Constraints{})
	u := widget.UnitSize(ctx)
	
	x, y := ebiten.WindowSize()
	if x < point.X && y < point.Y {
		point.X = u * 28
		point.Y = u * 18
	}
	return point
}

func (w *url_panel_widget_scrollable) SetURL(u *url.URL, ctx *gui.Context) {
	w.content.set_url(u, ctx)
}
