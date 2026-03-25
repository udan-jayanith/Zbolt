package def

import (
	"os"
	"time"
)

type Attribute struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

type HTTP_req_body struct {
	FilePath string `json:"filepath"`

	ContentType string `json:"content-type"`
	Content     string `json:"content"`
}

type HTTP_Data struct {
	Method string `json:"method"` // HTTP method

	URL struct {
		BaseURL string `json:"base-url"` // Everything before the path.

		Path struct {
			RawPath string `json:"raw-path"`
			Pattern struct {
				Pattern    string            `json:"pattern"`
				Attributes map[string]string `json:"attributes"`
			} `json:"pattern"`
		} `json:"path"` // Both path and pattern can't exists at once.

	} `json:"url"`

	Parameters []Attribute   `json:"parameters"`
	Headers    []Attribute   `json:"headers"`
	Body       HTTP_req_body `json:"body"` // Filepath of Content

	Response struct {
		AutoWrap bool `json:"auto-wrap"`
		Formate  bool `json:"formate"`
	} `json:"response-config"`

	temp TempHTTP_Data
}

func (data *HTTP_Data) Get_URL() string {
	return ""
}

func (data *HTTP_Data) TempData() *TempHTTP_Data {
	return &data.temp
}

type HTTP_res_body struct {
	File        *os.File
	ContentType string
	Content     string
}

type TempHTTP_Data struct {
	Status_code  int
	ResponseTime time.Duration
	ResponseSize int // In bytes
	Version      struct {
		Major, Minor int
	}
	Headers []Attribute
	ContentType string
	Body    HTTP_res_body
}
