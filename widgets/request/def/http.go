package def

import (
	"os"
	"time"
)

type Attribute struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
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
		
		Parameters []Attribute `json:"parameters"`
	} `json:"url"`

	Headers []Attribute `json:"headers"`
	Body    struct {
		FilePath string `json:"filepath"`

		ContentType string `json:"content-type"`
		Content     string `json:"content"`
	} `json:"body"` // Filepath of Content

	Response struct {
		AutoWrap bool `json:"auto-wrap"`
		Formate  bool `json:"formate"`
	} `json:"response-config"`
	
	temp TempHTTP_Data
}

func (data *HTTP_Data) TempData() *TempHTTP_Data{
	return &data.temp
}

type TempHTTP_Data struct {
	Status_code int
	ResponseTime time.Duration
	ResponseSize int // In bytes
	Version struct {
		Major, Minor int
	}
	Headers []Attribute
	Body struct {
		File *os.File
		ContentType string
		Content string
	}
}
