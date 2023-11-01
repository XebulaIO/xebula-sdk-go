package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/tidwall/sjson"
)

type HttpMethod string

const (
	HTTP_GET    HttpMethod = "GET"
	HTTP_POST   HttpMethod = "POST"
	HTTP_PUT    HttpMethod = "PUT"
	HTTP_PATCH  HttpMethod = "PATCH"
	HTTP_DELETE HttpMethod = "DELETE"
)

type HttpClient struct {
	client    *http.Client
	BasicAuth struct {
		Username string
		Password string
	}
}

type Header struct {
	Key   string
	Value string
}

func BuildClient() *http.Client {

	transport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   60 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	return &http.Client{Transport: transport}
}

func NewHttpClient() *HttpClient {
	return &HttpClient{client: BuildClient()}
}

func (h *HttpClient) SetBasicAuth(username string, password string) {
	h.BasicAuth.Username = username
	h.BasicAuth.Password = password
}

func (c *HttpClient) Get(rawurl string, out interface{}, headers ...Header) error {
	return c.do(rawurl, "GET", nil, out, headers...)
}

// helper function for making an http POST request.
func (c *HttpClient) Post(rawurl string, in, out interface{}, headers ...Header) error {
	return c.do(rawurl, "POST", in, out, headers...)
}

// helper function for making an http PUT request.
func (c *HttpClient) Put(rawurl string, in, out interface{}, headers ...Header) error {
	return c.do(rawurl, "PUT", in, out, headers...)
}

// helper function for making an http PATCH request.
func (c *HttpClient) Patch(rawurl string, in, out interface{}, headers ...Header) error {
	return c.do(rawurl, "PATCH", in, out, headers...)
}

// helper function for making an http DELETE request.
func (c *HttpClient) Delete(rawurl string, in, out interface{}, headers ...Header) error {
	return c.do(rawurl, "DELETE", in, out, headers...)
}

func (c *HttpClient) Do(rawurl string, method HttpMethod, in, out interface{}, headers ...Header) error {
	switch method {
	case HTTP_GET:
		return c.Get(rawurl, out, headers...)
	case HTTP_POST:
		return c.Post(rawurl, in, out, headers...)
	case HTTP_PUT:
		return c.Put(rawurl, in, out, headers...)
	case HTTP_PATCH:
		return c.Patch(rawurl, in, out, headers...)
	case HTTP_DELETE:
		return c.Delete(rawurl, in, out, headers...)
	default:
		return fmt.Errorf("invalid HTTP method: %s", method)
	}
}

func (c *HttpClient) do(rawurl, method string, in, out interface{}, headers ...Header) error {

	var (
		code = http.StatusTeapot
	)

	resp, err := c.open(rawurl, method, in, out, headers...)
	if err != nil {
		if resp != nil && resp.StatusCode > 0 {
			code = resp.StatusCode
		}
		return &HttpError{Status: code, Description: err.Error(), Method: method, URL: rawurl}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return &HttpError{Status: resp.StatusCode, Description: string(body), Method: method, URL: rawurl}
	}

	if buff, ok := out.(*bytes.Buffer); ok {
		_, err = buff.Write(body)
		return err
	}

	if out != nil {
		body, _ = sjson.SetBytes(body, "$cookies$", resp.Cookies())
		return json.NewDecoder(bytes.NewBuffer(body)).Decode(out)
	}

	return nil
}

func (c *HttpClient) open(rawurl, method string, in, out interface{}, headers ...Header) (*http.Response, error) {
	uri, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}

	// creates a new http request to bitbucket.
	req, err := http.NewRequest(method, uri.String(), nil)
	if err != nil {
		return nil, err
	}

	for _, h := range headers {
		req.Header.Add(h.Key, h.Value)
	}

	if c.BasicAuth.Username != "" || c.BasicAuth.Password != "" {
		req.SetBasicAuth(c.BasicAuth.Username, c.BasicAuth.Password)
	}

	// if we are posting or putting data, we need to
	// write it to the body of the request.
	if in != nil {
		rc, ok := in.(io.ReadCloser)
		if ok {
			req.Body = rc
			req.Header.Set("Content-Type", "plain/text")
		} else {
			inJson, err := json.Marshal(in)
			if err != nil {
				return nil, err
			}

			buf := bytes.NewBuffer(inJson)
			req.Body = ioutil.NopCloser(buf)

			req.ContentLength = int64(len(inJson))
			req.Header.Set("Content-Length", strconv.Itoa(len(inJson)))
			req.Header.Set("Content-Type", "application/json")
		}
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

type HttpError struct {
	Status      int    `json:"-"`
	Description string `json:"error_description"`
	Method      string `json:"method"`
	URL         string `json:"url"`
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("%s %s returned [%d]: %s", e.Method, e.URL, e.Status, e.Description)
}
