package CommonWidgets

import (
	"log"

	gui "github.com/guigui-gui/guigui"
	widget "github.com/guigui-gui/guigui/basicwidget"
)

type TextInputWithContextMenu struct {
	contextmenu widget.ContextMenuArea[string]
	widget.TextInput
}

func (w *TextInputWithContextMenu) on_contextmenu_item_selected(_ *gui.Context, _ int) {
	item, _ := w.contextmenu.PopupMenu().SelectedItem()

	switch item.Value {
	case "cut":
		w.TextInput.Cut()
	case "copy":
		w.TextInput.Copy()
	case "copy-all":
		log.Fatalln("Not implemented yet")
	case "paste":
		w.TextInput.Paste()
	case "select-all":
		w.TextInput.SelectAll()
	case "undo":
		w.TextInput.Undo()
	case "redo":
		w.TextInput.Redo()
	}
}

func (w *TextInputWithContextMenu) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	ctx.SetButtonInputReceptive(&w.contextmenu, true)
	editable := w.TextInput.IsEditable()
	if editable {
		w.contextmenu.PopupMenu().SetItems([]widget.PopupMenuItem[string]{
			{
				Text:    "Cut",
				KeyText: "Ctrl+x",
				Value:   "cut",
			},
			{
				Text:    "Copy",
				KeyText: "Ctrl+c",
				Value:   "copy",
			},
			{
				Text:    "Paste",
				KeyText: "Ctrl+v",
				Value:   "paste",
			},
			{
				Text:    "Select all",
				KeyText: "Ctrl+a",
				Value:   "select-all",
			},
			{
				Text:  "Copy all",
				Value: "copy-all",
			},
			{
				Text:    "Undo",
				KeyText: "Ctrl+z",
				Value:   "undo",
			},
			{
				Text:    "Redo",
				KeyText: "Ctrl+y",
				Value:   "redo",
			},
		})
	} else {
		w.contextmenu.PopupMenu().SetItems([]widget.PopupMenuItem[string]{
			{
				Text:    "Copy",
				KeyText: "Ctrl+c",
				Value:   "copy",
			},
			{
				Text:    "Select all",
				KeyText: "Ctrl+a",
				Value:   "select-all",
			},
			{
				Text:  "Copy all",
				Value: "copy-all",
			},
		})
	}

	w.contextmenu.PopupMenu().OnItemSelected(w.on_contextmenu_item_selected)
	adder.AddWidget(&w.contextmenu)

	w.TextInput.Build(ctx, adder)
	return nil
}

func (w *TextInputWithContextMenu) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	layouter.LayoutWidget(&w.contextmenu, widgetBounds.Bounds())
	w.TextInput.Layout(ctx, widgetBounds, layouter)
}
