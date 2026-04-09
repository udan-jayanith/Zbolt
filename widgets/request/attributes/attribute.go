package attr

type Attribute struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

type AttrCheck struct {
	Checked bool
	Key     string `json:"Key"`
	Value   string `json:"Value"`
}

// MergeAttrList merges a into b
func MergeAttrList(a []Attribute, b []Attribute) []Attribute {
	query_mapped := make(map[string]string, len(a))
	for _, attr := range a {
		query_mapped[attr.Key] = attr.Value
	}

	for i, attr := range b {
		val, ok := query_mapped[attr.Key]
		if ok {
			attr.Value = val
			b[i] = attr
		}
	}

	return b
}

// MergeAttrCheckList merges a into b
func MergeAttrCheckList(a []AttrCheck, b []AttrCheck) []AttrCheck {
	type attr_check struct {
		val     string
		checked bool
	}

	query_mapped := make(map[string]attr_check, len(a))
	for _, attr := range a {
		query_mapped[attr.Key] = attr_check{
			val:     attr.Value,
			checked: attr.Checked,
		}
	}

	for i, attr := range b {
		attr_check, ok := query_mapped[attr.Key]
		if ok {
			attr.Value = attr_check.val
			attr.Checked = attr_check.checked
			b[i] = attr
		}
	}

	return b
}
