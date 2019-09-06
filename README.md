# gorequests
just a wrapper of web request.

## how to use it

```
go get -u github.com/guoruibiao/gorequests
```

Get
```go
url := "http://vps.XXXXXXX.com/XXXXXXX/201909/test-requests-go.php"
	url += "?name=guoruibiao&age=25&address=北京朝阳"
	response, err := gorequests.NewRequest("GET", url).DoRequest()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.Content())
```

POST
```go
url := "http://vps.XXXXX.com/XXXXXX/201909/test-requests-go.php"
	payload := map[string]string{
		"name": "guoruibiao",
		"age": "25",
		"school": "大连理工大学",
	}
	response, err := gorequests.NewRequest("POST", url).Body(payload).DoRequest()
	if err != nil {
		fmt.Println(err)
	}
	//t.Log(response.Content())
	type user struct {
		Username string `json:"name"`
		Userage int `json:"age"`
		Userschool string `json:"school"`
	}
	content, _ := response.Content()
	fmt.Println(content)
	u := &user{}
	err = json.Unmarshal([]byte(content), u)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(u)
```