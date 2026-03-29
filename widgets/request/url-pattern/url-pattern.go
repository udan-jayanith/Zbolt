package url_pattern

import (
	"strings"
)

type Attribute struct {
	Checked bool
	
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

type Pattern struct {
	raw_path string
	Values   map[string]int
	List     []Attribute
}

func (r *Pattern) Path() string {
	p := r.raw_path
	for _, attr := range r.List {
		p = strings.Replace(p, "{"+attr.Key+"}", attr.Value, 1)
	}
	return p
}

func (r *Pattern) Set(key, value string) {
	r.List[r.Values[key]].Value = value
}

func ParsePattern(pattern string) (Pattern, error) {
	list := make([]Attribute, 0, 4)
	value := make(map[string]int, 4)
	
	i := -1
	var idx int
	var attr Attribute
	
	for j, char := range pattern {
		if char == '{' && i == -1 {
			i = j
		} else if char == '}' && i > -1 {
			attr.Key = pattern[i+1:j] 
			value[attr.Key] = idx
			list = append(list, attr)
			idx++
			i = -1
		}
	}

	return Pattern{
		raw_path: pattern,
		Values:   value,
		List: list,
	}, nil
}
