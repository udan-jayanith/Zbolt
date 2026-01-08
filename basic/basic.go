package basic

import(
	gui "github.com/guigui-gui/guigui"
)

const (
	HeadingSize float64 = 1.8
)

func NewPadding(padding ...int) gui.Padding {
	l := len(padding)
	if l > 4 {
		panic("Extra arguments to NewPadding function")
	} else if l == 3 || l == 0 {
		panic("Invalid number of arguments to NewPadding function")
	} else if l == 4 {
		return gui.Padding{
			Top:    padding[0],
			Bottom: padding[2],
			Start:  padding[3],
			End:    padding[1],
		}
	} else if l == 2 {
		return gui.Padding{
			Top:    padding[0],
			Bottom: padding[0],
			Start:  padding[1],
			End:    padding[1],
		}
	}
	return gui.Padding{
		Top:    padding[0],
		Bottom: padding[0],
		Start:  padding[0],
		End:    padding[0],
	}
}

type Alignment int8

const (
	Start Alignment = iota + 1
	Center
	End
)

func align(items []gui.LinearLayoutItem, alignment Alignment) []gui.LinearLayoutItem {
	switch alignment {
	case Start:
		items = append(items, gui.LinearLayoutItem{
			Size: gui.FlexibleSize(1),
		})
	case Center:
		items = append([]gui.LinearLayoutItem{
			{
				Size: gui.FlexibleSize(1),
			},
		}, items...)
		items = append(items, gui.LinearLayoutItem{
			Size: gui.FlexibleSize(1),
		})
	case End:
		items = append(items, gui.LinearLayoutItem{
			Size: gui.FlexibleSize(1),
		})
	}
	return items
}

func Align(items []gui.LinearLayoutItem, horizontal, vertical Alignment) gui.LinearLayout {
	items = align(items, vertical)
	vertical_layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionVertical,
		Items:     items,
	}

	horizontal_layout := gui.LinearLayout{
		Direction: gui.LayoutDirectionHorizontal,
		Items: []gui.LinearLayoutItem{
			{
				Layout: vertical_layout,
			},
		},
	}
	horizontal_layout.Items = align(horizontal_layout.Items, horizontal)

	return horizontal_layout
}