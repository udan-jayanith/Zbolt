package def

import (
	"weak"

	gui "github.com/guigui-gui/guigui"
)

type RequestType uint8

const (
	HTTP RequestType = iota + 0
	Websocket
	GraphQL
	Grpc
)

func (t RequestType) IconName() string {
	switch t {
	case HTTP:
		return "large-icons/http"
	case Websocket:
		return "large-icons/websocket"
	case GraphQL:
		return "large-icons/graphql"
	case Grpc:
		return "large-icons/grpc"
	default:
		panic("Unknown request type")
	}
}

type Request struct {
	Type RequestType
	Name string
	data weak.Pointer[any]
}

func (r *Request) Data() any {
	return nil
}

type RequestWidget interface {
	gui.Widget
	RequestType() RequestType
}
