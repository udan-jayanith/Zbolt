package attr

type Attribute struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

type AttrCheck struct {
	Checked bool
	Key   string `json:"Key"`
	Value string `json:"Value"`
}