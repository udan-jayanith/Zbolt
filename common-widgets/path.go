package CommonWidgets

import (
	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type path_separator_widget struct{
	gui.DefaultWidget
	
	sep widget.Text
}

type path_segment_widget struct{
	gui.DefaultWidget
	
	text widget.Text
}

type path_widget struct{
	gui.DefaultWidget
	
	segments []gui.Widget
}

type Path struct {
	gui.DefaultWidget
	
	panel widget.Panel
	path_widget path_widget
}