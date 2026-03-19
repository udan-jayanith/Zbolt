package http_widget_test

import(
	"testing"
	"API-Client/widgets/request/page/http"
)

func TestParse_url_path_query(t *testing.T){
	query, _ := http_widget.Parse_url_path_query("/{user}/{repo}")
	_, ok := query.List["user"]
	if !ok {
		t.Fatal("Expected user but not found")
	}
	
	_, ok = query.List["repo"]
	if !ok {
		t.Fatal("Expected repo but not found")
	}
}

