package message_model

import (
	"image"
	"log"

	queue "github.com/golang-ds/queue/slicequeue"
	gui "github.com/guigui-gui/guigui"
)

type MessageModelType uint8

const (
	Alert        MessageModelType = iota + 0 // Used for showing some kind of error or importent message
	Prompt                                   // Ask a user for a yes no answer
	Notification                             // Shows a notification in the bouttom right corner
	Notify                                   // Used for notifying user about somthing
)

type result_fn_type = func(ok bool, ctx *gui.Context)

type message_data struct {
	message   string
	t         MessageModelType
	on_result result_fn_type
}

var messsage_queue queue.SliceQueue[message_data]

func Show(message string, t MessageModelType, on_result func(ok bool, ctx *gui.Context)) {
	messsage_queue.Enqueue(message_data{
		message:   message,
		t:         t,
		on_result: on_result,
	})
}

type message_model interface {
	gui.Widget
	SetMessage(message string)
	OnResult(fn func(ok bool, ctx *gui.Context))
	Bounds(ctx *gui.Context, widgetBounds *gui.WidgetBounds) image.Rectangle
//	Open(open bool)
//	IsOpen() bool
}

type message_model_widget struct {
	gui.DefaultWidget
	is_showing_message   bool
	current_message_data message_data

	alert_widget alert_widget
	notify       notify_widget
}

func (wi *message_model_widget) Build(ctx *gui.Context, adder *gui.ChildAdder) error {
	if !wi.is_showing_message && !messsage_queue.IsEmpty() {
		wi.is_showing_message = true
		wi.current_message_data, _ = messsage_queue.Dequeue()
	} else if !wi.is_showing_message {
		return nil
	}

	var modeled_widget message_model
	switch wi.current_message_data.t {
	case Alert:
		modeled_widget = &wi.alert_widget
	case Notify:
		modeled_widget = &wi.notify
	default:
		//case Prompt:
		//case Notification:
		log.Fatalln("Not implemented")
	}
	modeled_widget.SetMessage(wi.current_message_data.message)
	modeled_widget.OnResult(func(ok bool, ctx *gui.Context) {
		wi.is_showing_message = false
		if wi.current_message_data.on_result != nil {
			wi.current_message_data.on_result(ok, ctx)
		}
	})
	adder.AddWidget(modeled_widget)

	return nil
}

func (wi *message_model_widget) Layout(ctx *gui.Context, widgetBounds *gui.WidgetBounds, layouter *gui.ChildLayouter) {
	if !wi.is_showing_message {
		return
	}

	var modeled_widget message_model
	switch wi.current_message_data.t {
	case Alert:
		modeled_widget = &wi.alert_widget
	case Notify:
		modeled_widget = &wi.notify
	default:
		//case Prompt:
		//case Notification:
		log.Fatalln("Not implemented")
	}
	layouter.LayoutWidget(modeled_widget, modeled_widget.Bounds(ctx, widgetBounds))
}

// This whidget only must be used by the rott widget
var MessageModel = message_model_widget{}
