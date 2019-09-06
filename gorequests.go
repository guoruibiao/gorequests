package gorequests

import (
	"net/http"
	"bytes"
			"io/ioutil"
	"strings"
	"fmt"
	)

// wrapper of common web request.
type (
	// Request
	Request struct {
		method string
		url string
		headers map[string]string
		body map[string]string
		auth BasicAuth
	}
	// Response
	Response struct {
		*http.Response
		content string
	}
	// BasicAuth
	BasicAuth struct {
		username string
		password string
	}
)

func NewRequest(method, url string) *Request {
	return &Request{
		method: method,
		url: url,
		headers: make(map[string]string),
	}
}

func (this *Request) Headers(headers map[string]string) *Request {
	if headers != nil {
		for key, value := range headers {
			this.headers[key] = value
		}
	}
	return this
}

func (this *Request) Body(payload map[string]string) *Request {
	if payload != nil {
		this.body = payload
	}
	return this
}

func (this *Request) BasicAuth(username, password string) *Request {
	if username!="" && password!="" {
		this.auth = BasicAuth{
			username: username,
			password: password,
		}
	}
	return this
}


func (this *Request) DoRequest() (*Response, error) {
	// handle body data
	body := bytes.NewBuffer(nil)
	payload := ""
	if len(this.body) != 0 {
		for key, value := range this.body {
			body.WriteString(key + "=" + value + "&")
		}
		// trim the last `&`
		payload = strings.TrimRight(body.String(), "&")
	}
	// the third format must be `xxx=yyy&zzz=aaa&...`
	req, err := http.NewRequest(this.method, this.url, strings.NewReader(payload))
	if err != nil {
		return nil, err
	}
	// POST need to set the header['Content-Type'] = 'application/x-www-form-urlencoded'
	if strings.ToLower(this.method) == "post" && this.body != nil {
		this.headers["Content-Type"] = "application/x-www-form-urlencoded"
	}
	// apply headers to request.header
	if this.headers != nil {
		for key, value := range this.headers {
			// todo what the difference of Add and Set ?
			req.Header.Set(key, value)
		}
	}
	fmt.Println(req.Header)
	// create proxy to handle request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	// anything seems all right
	return &Response{resp, ""}, nil
}

func (this *Response) Content() (string, error) {
	bytes, err := ioutil.ReadAll(this.Body)
	if err != nil {
		return "", err
	}
	defer this.Body.Close()
	this.content = string(bytes)
	return string(bytes), nil
}

func (this *Response) StatusCode() int {
	return this.Response.StatusCode
}

