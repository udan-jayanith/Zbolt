package def

import (
	"path/filepath"
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
	path string
	data weak.Pointer[any]
}

func (r *Request) Data() any {
	return nil
}

func (r *Request) Path() string {
	return r.path
}

func (r *Request) Name() string {
	return filepath.Base(r.path)
}

func NewRequest(t RequestType, path string) Request{
	return Request{
		Type: t,
		path: path,
	}
}

type RequestWidget interface {
	gui.Widget
	RequestType() RequestType
}

type Folder struct {
	path string
}

func (r *Folder) Path() string {
	return r.path
}

func (r *Folder) Name() string {
	return filepath.Base(r.path)
}

func NewFolder(path, name string) Folder{
	return Folder{
		path: filepath.Join(path, name),
	}
}

type FolderOrFile interface {
	Path() string
	Name() string
}
