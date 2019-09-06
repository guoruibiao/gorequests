package gorequests

import (
	"testing"
	"encoding/json"
)

func TestRequest_DoRequest_GET(t *testing.T) {
	url := ""
	url += "?name=guoruibiao&age=25&address=北京朝阳"
	response, err := NewRequest("GET", url).DoRequest()
	if err != nil {
		t.Error(err)
	}
	t.Log(response.Content())
}

func TestRequest_DoRequest_POST(t *testing.T) {
	url := ""
	payload := map[string]string{
		"name": "guoruibiao",
		"age": "25",
		"school": "大连理工大学",
	}
	response, err := NewRequest("POST", url).Body(payload).DoRequest()
	if err != nil {
		t.Error(err)
	}
	//t.Log(response.Content())
	type user struct {
		Username string `json:"name"`
		Userage int `json:"age"`
		Userschool string `json:"school"`
	}
	content, _ := response.Content()
	t.Log(content)
	u := &user{}
	err = json.Unmarshal([]byte(content), u)
	if err != nil {
		t.Error(err)
	}
	t.Log(u)
}

