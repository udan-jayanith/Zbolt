package http_widget

type url_query_data struct {
	Value      string
	start, end int
}

type url_path_query struct {
	raw_path []byte
	List     map[string]url_query_data
}

func (r *url_path_query) Path() string {
	return ""
}

func Parse_url_path_query(url_path string) (url_path_query, error) {
	list := make(map[string]url_query_data, 3)

	var data url_query_data
	var opening bool
	for i, char := range url_path {
		if char == '{' && !opening {
			data.start = i
			opening = true
		} else if char == '}' && opening {
			data.end = i
			list[url_path[data.start+1:data.end]] = data
			opening = false
		}
	}

	return url_path_query{
		raw_path: []byte(url_path),
		List:     list,
	}, nil
}
