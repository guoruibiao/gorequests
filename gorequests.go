package gorequests

import (
	"net/http"
	"bytes"
			"io/ioutil"
	"strings"
		"encoding/base64"
	"encoding/json"
)

// wrapper of common web request.
type (
	// Request
	Request struct {
		/** GET/POST  **/
		method string
		url string
		/**  extra http header **/
		headers map[string]string

		/** Content-Type=application/x-www-form-urlencoded, values in form always be `map[string][string]`  **/
		form map[string]string

		/** multipart style post payload parameters. maybe `map[string]interface{}` **/
		body map[string]interface{}

		/**basic http authorization **/
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

// special for header with `Content-Type=application/json`, it should never be used if `Request.Form` called.
func (this *Request) Body(payload map[string]interface{}) *Request {
	if payload != nil {
		this.body = payload
	}
	return this
}

// special for header with `Content-Type=application/x-www-form-urlencoded`, it should never be used if `Request.Body` called.
func (this *Request)Form(form map[string]string) *Request {
	if form != nil {
		this.form = form
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
	var payload string

	// for Content-Type = application/x-www-form-urlencoded
	if strings.ToLower(this.method) == "post" && this.form != nil {
		this.headers["Content-Type"] = "application/x-www-form-urlencoded"
		if len(this.form) != 0 {
			for key, value := range this.form {
				body.WriteString(key + "=" + value + "&")
			}
			// trim the last `&`
			payload = strings.TrimRight(body.String(), "&")
		}
	}

	// for Content-Type = application/json
	if strings.ToLower(this.method) == "post" && this.body != nil {
		this.headers["Content-Type"] = "application/json"
		if len(this.body) != 0 {
			marshalBytes, err := json.Marshal(this.body)
			if err != nil {
				return nil, err
			}
			payload = string(marshalBytes)
		}
	}

	// init the http client.
	req, err := http.NewRequest(this.method, this.url, strings.NewReader(payload))
	if err != nil {
		return nil, err
	}

	// set the basic Authorization
	if this.auth.username != "" && this.auth.password != "" {
		req.SetBasicAuth(this.auth.username, this.auth.password)
		req.Header.Add("Authorization", "Basic " + base64.StdEncoding.EncodeToString([]byte(this.auth.username + ":" + this.auth.password)))
	}

	// apply headers to request.Header
	if this.headers != nil {
		for key, value := range this.headers {
			req.Header.Set(key, value)
		}
	}
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

