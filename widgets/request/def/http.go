package def

import (
	attr "API-Client/widgets/request/attributes"
	url_utils "API-Client/widgets/request/url-utils"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type HTTP_Request_Body struct {
	ContentType ContentType `json:"content-type"`
	Content     string      `json:"content"`
}

type URL struct {
	BaseURL string `json:"base-url"` // Everything before the path.

	Path struct {
		RawPath string `json:"raw-path"`
		Pattern struct {
			Pattern    string           `json:"pattern"`
			Attributes []attr.Attribute `json:"attributes"`
		} `json:"pattern"`
	} `json:"path"` // Both path and pattern can't exists at once.
}

func (u *URL) IsPattern() bool {
	return len(u.Path.Pattern.Attributes) > 0
}

// EncodedPath returns the encoded path
func (u *URL) EncodedPath() string {
	if !u.IsPattern() {
		return u.Path.RawPath
	}

	pattern, _ := url_utils.ParsePattern(u.Path.Pattern.Pattern)
	for _, attr := range u.Path.Pattern.Attributes {
		pattern.Set(attr.Key, attr.Value)
	}
	return pattern.Path()
}

// RawPath returns the raw-path if exists otherwise returns the raw-path-pattern without being encoded.
func (u *URL) RawPath() string {
	if u.IsPattern() {
		return u.Path.Pattern.Pattern
	}
	return u.Path.RawPath
}

func (u *URL) SetPattern(pattern string, attributes []attr.Attribute) {
	u.Path.RawPath = ""
	u.Path.Pattern.Pattern = pattern
	u.Path.Pattern.Attributes = attributes
}

func (u *URL) SetPath(path string) {
	u.Path.RawPath = path
	u.Path.Pattern.Pattern = ""
	u.Path.Pattern.Attributes = []attr.Attribute{}
}

type HTTP_Data struct {
	Method string `json:"method"` // HTTP method

	URL URL `json:"url"`

	Parameters []attr.AttrCheck  `json:"parameters"`
	Headers    []attr.AttrCheck  `json:"headers"`
	Body       HTTP_Request_Body `json:"body"` // Filepath of Content

	RequestConfig struct {
		AutoWrap bool `json:"auto-wrap"`
		Formate  bool `json:"formate"`
	} `json:"request-config"`

	ResponseConfig struct {
		AutoWrap bool `json:"auto-wrap"`
		Formate  bool `json:"formate"`
	} `json:"response-config"`

	selected_request_tab int

	request struct {
		is_fetching, canceled     bool
		response_data_update_time time.Time
		err                       error
		cancel                    chan struct{}
		m                         sync.Mutex
		response_data             HTTP_Response_Data
	}
	response_data HTTP_Response_Data
}

func (data *HTTP_Data) SetSelectedRequestTab(index int) {
	data.selected_request_tab = index
}

func (data *HTTP_Data) SelectedRequestTab() int {
	return data.selected_request_tab
}

/*
Adapted from Golang net/http package.
example: username=edger&age=20
*/
// TODO: Make this a separate public function
func (data *HTTP_Data) EncodedParameters() string {
	return url_utils.EncodeParameters(data.Parameters)
}

// GetUrl return the full url.
func (data *HTTP_Data) FullURL() *url.URL {
	u, _ := url.Parse(data.URL.BaseURL)
	u.RawPath = data.URL.EncodedPath()
	u.RawQuery = data.EncodedParameters()
	return u
}

func (data *HTTP_Data) update_response_data() {
	if !data.request.m.TryLock() {
		return
	}
	defer data.request.m.Unlock()

	data.response_data = data.request.response_data
}

func (data *HTTP_Data) ResponseData() *HTTP_Response_Data {
	if data.request.is_fetching {
		data.update_response_data()
	}
	return &data.response_data
}

func (data *HTTP_Data) set_req_headers(req *http.Request) {
	req_headers_mapped := make(map[string]string, len(req.Header))
	for key, vals := range req.Header {
		req_headers_mapped[key] = strings.Join(vals, ",")
	}

	for i, header := range data.Headers {
		if len(req_headers_mapped) == 0 {
			break
		}
		val, ok := req_headers_mapped[header.Key]
		if ok {
			header.Checked = true
			header.Value = val
			data.Headers[i] = header
			delete(req_headers_mapped, header.Key)
		}
	}

	for k, v := range req_headers_mapped {
		data.Headers = append([]attr.AttrCheck{
			{
				Checked: true,
				Key:     k,
				Value:   v,
			},
		}, data.Headers...)
	}

	if data.Body.ContentType != "" {
		var content_type_found bool
		for i, header := range data.Headers {
			if header.Key == "Content-Type" {
				header.Value = string(data.Body.ContentType)
				header.Checked = true
				data.Headers[i] = header
				content_type_found = true
				break
			}
		}

		if !content_type_found {
			data.Headers = append([]attr.AttrCheck{
				{
					Checked: true,
					Key:     "Content-Type",
					Value:   string(data.Body.ContentType),
				},
			}, data.Headers...)
		}
	}

	for _, header := range data.Headers {
		if !header.Checked {
			continue
		}
		req.Header.Set(header.Key, header.Value)
	}
}

func (data *HTTP_Data) GrabRequestErr() error {
	err := data.request.err
	data.request.err = nil
	return err
}

func (data *HTTP_Data) close_request() {
	data.update_response_data()

	data.request.response_data = HTTP_Response_Data{}
	close(data.request.cancel)
	data.request.is_fetching = false
	data.request.canceled = false
}

// Do performs the http request
// Response data can be revised through ResponseData method
// Calling Do updates Headers so headers must be update in the HTTP_widget
func (data *HTTP_Data) Do() bool {
	if data.request.is_fetching {
		panic("Request is already being requested")
	}
	data.request.is_fetching = true
	data.request.err = nil
	data.request.cancel = make(chan struct{}, 1)

	method := strings.ToUpper(data.Method)
	var body io.Reader
	if method == "POST" || method == "PUT" || method == "PATCH" {
		body = strings.NewReader(data.Body.Content)
	}

	req, err := http.NewRequest(method, data.FullURL().String(), body)
	if err != nil {
		data.request.err = err
		data.close_request()
		return false
	}
	data.set_req_headers(req)
	go data.do(req)

	return true
}

func (data *HTTP_Data) update_res_data_with(res_data HTTP_Response_Data) {
	locked := data.request.m.TryLock()
	if !locked {
		return
	}
	defer data.request.m.Unlock()

	body_content_copied := make([]byte, len(res_data.Body.content))
	copy(body_content_copied, res_data.Body.content)

	headers_copied := make([]attr.AttrCheck, len(res_data.Headers))
	copy(headers_copied, res_data.Headers)

	res_data.SelectedResponseTab = data.request.response_data.SelectedResponseTab

	data.request.response_data = res_data
	data.request.response_data.Headers = headers_copied
	data.request.response_data.Body.content = body_content_copied
}

func (data *HTTP_Data) do(req *http.Request) {
	res_data := HTTP_Response_Data{}
	response_time := time.Now()
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		data.request.err = err
		data.close_request()
		return
	}
	defer res.Body.Close()

	res_data.Status_code = res.StatusCode
	res_data.Version = Version{
		Major: res.ProtoMajor,
		Minor: res.ProtoMinor,
	}
	res_data.Body.ContentType = ContentType(res.Header.Get("Content-Type"))
	res_data.Headers = http_headers_to_attr_check(res.Header)
	res_data.ResponseTime = time.Since(response_time)
	data.update_res_data_with(res_data)

	body_content := make([]byte, 0, 1024*2)
	p := make([]byte, 1024)
	update_time := time.Now()
res_body_reader_loop:
	for {
		n, err := res.Body.Read(p)
		if err != nil {
			data.request.err = err
			break
		}

		select {
		case <-data.request.cancel:
			break res_body_reader_loop
		default:
		}

		body_content = append(body_content, p[:n]...)
		res_data.ResponseSize = len(body_content)
		res_data.ResponseTime = time.Since(response_time)
		if time.Since(update_time).Milliseconds() >= 500 {
			data.update_res_data_with(res_data)
		}
	}
	data.close_request()
}

func (data *HTTP_Data) CancelRequest() error {
	if data.request.canceled {
		return errors.New("Request is already being canceled")
	}
	data.request.cancel <- struct{}{}
	data.request.canceled = true
	return nil
}

type HTTP_Response_Body struct {
	// TODO:
	//File         *os.File
	//Path         string
	//IsFileClosed bool

	ContentType ContentType
	content     []byte
}

func (b *HTTP_Response_Body) Content() []byte {
	return b.content
}

type Version struct {
	Major, Minor int
}

type HTTP_Response_Data struct {
	Status_code  int
	ResponseTime time.Duration
	ResponseSize int // In bytes // TODO: change this to Content-Length since header size is quite hard calculate
	Version      Version
	Headers      []attr.AttrCheck
	Body         HTTP_Response_Body

	SelectedResponseTab int
}
