package http_widget

import (
	"strings"
)

type attribute struct {
	K, V string
}

type url_path_query struct {
	raw_path string
	List     []attribute
}

func (r *url_path_query) Path() string {
	p := r.raw_path
	for _, v := range r.List {
		p = strings.Replace(p, "{"+v.K+"}", v.V, 1)
	}
	return p
}

// This isn't robust but enough.
func Parse_url_path_query(url_path string) (url_path_query, error) {
	list := make([]attribute, 0, 4)

	var data attribute
	j := -1
	for i, char := range url_path {
		if char == '{' && j == -1 {
			j = i
		} else if char == '}' && j > -1 {
			data.K = url_path[j+1 : i]
			list = append(list, data)
			j = -1
		}
	}

	return url_path_query{
		raw_path: url_path,
		List:     list,
	}, nil
}
