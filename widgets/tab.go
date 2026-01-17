package CWidget

import (
	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type TabItem struct {
	gui.DefaultWidget
	text widget.Text
	body gui.Widget
}

type TabBar struct {
	gui.DefaultWidget
	tab_items []*TabItem
}

type tab_body struct {
	gui.DefaultWidget
}

type TabContainer struct {
	gui.DefaultWidget
	tab_bar *TabBar
	body tab_body
}

