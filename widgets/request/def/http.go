package def

import (
	attr "github.com/udan-jayanith/Zbolt/widgets/request"
	url_pattern "github.com/udan-jayanith/Zbolt/widgets/request/url-pattern"
	"net/url"
	"os"
	"strings"
	"time"
)

type HTTP_Request_Body struct {
	FilePath   string `json:"filepath"`
	IsFileOpen bool

	ContentType ContentType `json:"content-type"`
	Content     string      `json:"content"`
}

type HTTP_Data struct {
	Method string `json:"method"` // HTTP method

	URL struct {
		BaseURL string `json:"base-url"` // Everything before the path.

		Path struct {
			RawPath string `json:"raw-path"`
			Pattern struct {
				Pattern    string           `json:"pattern"`
				Attributes []attr.Attribute `json:"attributes"`
			} `json:"pattern"`
		} `json:"path"` // Both path and pattern can't exists at once.
	} `json:"url"`

	Parameters []attr.AttrCheck  `json:"parameters"`
	Headers    []attr.AttrCheck  `json:"headers"`
	Body       HTTP_Request_Body `json:"body"` // Filepath of Content

	ResponseConfig struct {
		AutoWrap bool `json:"auto-wrap"`
		Formate  bool `json:"formate"`
	} `json:"response-config"`

	response_data HTTP_Response_Data
}

func (data *HTTP_Data) path() string {
	if data.URL.Path.RawPath != "" {
		return data.URL.Path.RawPath
	}

	pattern, _ := url_pattern.ParsePattern(data.URL.Path.Pattern.Pattern)
	for _, attr := range data.URL.Path.Pattern.Attributes {
		pattern.Set(attr.Key, attr.Value)
	}
	return pattern.Path()
}

/*
Adapted from Golang net/http package.
*/
func (data *HTTP_Data) EncodedParameters() string {

	parameters := data.Parameters
	if len(parameters) == 0 {
		return ""
	}
	var buf strings.Builder
	// This assumes key and values is about length of 5 each.
	buf.Grow(len(parameters) * 10)

	for _, attr := range parameters {
		key := url.QueryEscape(attr.Key)
		value := url.QueryEscape(attr.Value)

		if buf.Len() > 0 {
			buf.WriteByte('&')
		}

		buf.WriteString(key)
		buf.WriteByte('=')
		buf.WriteString(value)
	}
	return buf.String()
}

func (data *HTTP_Data) GetUrl() *url.URL {
	u, _ := url.Parse(data.URL.BaseURL)
	u.RawPath = data.path()
	u.RawQuery = data.EncodedParameters()
	return u
}

func (data *HTTP_Data) ResponseData() *HTTP_Response_Data {
	return &data.response_data
}

func (data *HTTP_Data) Do() {

}

// UpdateResponseData updates HTTP_Response_Data if http response data isn't locked
func (data *HTTP_Data) UpdateResponseData() {
}

type HTTP_Response_Body struct {
	File         *os.File
	Path         string
	IsFileClosed bool

	ContentType ContentType
	Content     string
}

type HTTP_Response_Data struct {
	Status_code  int
	ResponseTime time.Duration
	ResponseSize int // In bytes
	Version      struct {
		Major, Minor int
	}
	Headers []attr.AttrCheck
	Body    HTTP_Response_Body
}
