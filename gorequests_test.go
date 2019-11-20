package gorequests

import (
	"testing"
	"encoding/json"
	"fmt"
)
var (
	url = "http://xxxxxx/201911/tmp.php";
)
func TestRequest_DoRequest_GET(t *testing.T) {
	url += "?name=guoruibiao&age=25&address=北京朝阳"
	response, err := NewRequest("GET", url).DoRequest()
	if err != nil {
		t.Error(err)
	}
	t.Log(response.Content())
}

func TestRequest_DoRequest_POST(t *testing.T) {
	payload := map[string]string{
		"name": "guoruibiao",
		"age": "25",
		"school": "大连理工大学",
	}
	response, err := NewRequest("POST", url).Form(payload).DoRequest()
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

func TestRequest_DoRequest_APPJSON(t *testing.T) {
	content := "Content-Type=application/json"
	payload := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": content,
		},
		"at":  map[string]interface{}{
			"atMobiles": nil,
			"isAtAll": true,
		},
	}
	resp, err := NewRequest("POST", "https://oapi.dingtalk.com/robot/send?access_token=b716e1f39b7fc7afbea04b23909f4bb79db65a117d589f886d1757").Body(payload).DoRequest()
	if err != nil {
		t.Error(err)
	}
	content, _ = resp.Content()
	fmt.Println(content)
}
