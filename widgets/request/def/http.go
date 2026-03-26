package def

import (
	"net/url"
	"os"
	"time"
)

type Attribute struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

type HTTP_Request_Body struct {
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
		
		u *url.URL

	} `json:"url"`

	Parameters []Attribute   `json:"parameters"`
	Headers    []Attribute   `json:"headers"`
	Body       HTTP_Request_Body `json:"body"` // Filepath of Content

	ResponseConfig struct {
		AutoWrap bool `json:"auto-wrap"`
		Formate  bool `json:"formate"`
	} `json:"response-config"`

	response_data HTTP_Response_Data
}

func (data *HTTP_Data) Do() {
	
}

func (data *HTTP_Data) Get_URL(parameters bool) string {
	return ""
}

func (data *HTTP_Data) ResponseData() *HTTP_Response_Data {
	return &data.response_data
}

func (data *HTTP_Data) UpdateResponseData() {
}

type HTTP_Response_Body struct {
	File        *os.File
	ContentType string
	Content     string
}

type HTTP_Response_Data struct {
	Status_code  int
	ResponseTime time.Duration
	ResponseSize int // In bytes
	Version      struct {
		Major, Minor int
	}
	Headers []Attribute
	ContentType string
	Body    HTTP_Response_Body
}
