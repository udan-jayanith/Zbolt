package http_widget_test

import(
	"testing"
	"API-Client/widgets/request/page/http"
)

func TestParse_url_path_query(t *testing.T){
	query, _ := http_widget.Parse_url_path_query("/{user}/{repo}")
	_ , ok := query.List["user"]
	if !ok {
		t.Fatal("Expected user but not found")
	}
	query.Set("user", "udan-jayanith")
	
	_, ok = query.List["repo"]
	if !ok {
		t.Fatal("Expected repo but not found")
	}
	query.Set("repo", "zbolt")
	
	if string(query.Path()) != "/udan-jayanith/zbolt" {
		t.Fatal("Unexpected output")
	}
}

