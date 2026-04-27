package CommonWidgets

import (
	"reflect"

	gui "github.com/guigui-gui/guigui"
)

// Adapted from GUIGUI source code
type lazy_widget[T gui.Widget] struct {
	widget T
	is_set bool
}

func (l *lazy_widget[T]) Widget() T {
	if !l.is_set {
		t := reflect.TypeFor[T]()
		if t.Kind() == reflect.Pointer {
			l.widget = reflect.New(t.Elem()).Interface().(T)
		}
		l.is_set = true
	}
	return l.widget
}

// This is not nil safe
func (l *lazy_widget[T]) SetWidget(widget T) {
	l.widget = widget
}
