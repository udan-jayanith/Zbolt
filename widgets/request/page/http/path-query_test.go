package http_widget_test

import(
	"testing"
	"API-Client/widgets/request/page/http"
)

func TestParse_url_path_query(t *testing.T){
	query, _ := http_widget.Parse_url_path_query("/{user}/{repo}")
	
	if string(query.List[0].K) != "user"{
		t.Fatal("Expected user but got", query.List[0].K)
	}
	query.List[0].V = "udan-jayanith"
	
	if string(query.List[1].K) != "repo"{
		t.Fatal("Expected repo but got", query.List[0].K)
	}
	query.List[1].V = "zbolt"
	
	if string(query.Path()) != "/udan-jayanith/zbolt" {
		t.Fatal("Unexpected output", string(query.Path()))
	}
}

