package CommonWidgets

import (
	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type path_segment_widget struct{
	gui.DefaultWidget
	
	text widget.Text
}

type path_widget struct{
	gui.DefaultWidget
	
	segments []path_segment_widget
}

type Path struct {
	gui.DefaultWidget
	
	panel widget.Panel
	path_widget path_widget
}