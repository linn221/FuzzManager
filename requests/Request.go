package requests

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Request struct {
	// request
	Base          string
	Prams         map[string]string
	Method        string
	RequestHeader map[string]string
	RequestBody   []byte

	// response
	Status         int
	ResponseHeader map[string]string
	ResponseBody   []byte
	Latency        time.Duration
}

func NewRequest(base string) *Request {
	return &Request{
		Base:           base,
		Prams:          make(map[string]string),
		RequestHeader:  make(map[string]string),
		ResponseHeader: make(map[string]string),
	}
}

// func (r *Request) Send() error {
// 	// Build the query parameters
// 	params := url.Values{}
// 	for k, v := range r.Prams {
// 		params.Add(k, v)
// 	}

// 	// Parse base URL
// 	reqURL, err := url.Parse(r.Base)
// 	if err != nil {
// 		log.Fatalf("Invalid URL: %v", err)
// 	}

// 	// Attach query params
// 	reqURL.RawQuery = params.Encode()

// 	// Create the GET request
// 	req, err := http.NewRequest("GET", reqURL.String(), nil)
// 	if err != nil {
// 		log.Fatalf("Failed to create request: %v", err)
// 	}
// }

func (r *Request) Clone() *Request {
	requestHeaders := make(map[string]string, len(r.RequestHeader))
	params := make(map[string]string, len(r.Prams))
	for k, v := range r.RequestHeader {
		requestHeaders[k] = v
	}
	for k, v := range r.Prams {
		params[k] = v
	}
	requestBody := make([]byte, len(r.RequestBody))
	copy(requestBody, r.RequestBody)

	nr := NewRequest(r.Base)
	nr.Base = r.Base
	nr.Method = r.Method
	nr.Prams = params
	nr.RequestBody = requestBody
	nr.RequestHeader = requestHeaders
	return nr
}

func (r *Request) StdRequest() (*http.Request, error) {
	// parse url
	url, err := url.Parse(r.Base)
	if err != nil {
		return nil, err
	}

	if len(r.Prams) > 0 {
		q := url.Query()
		for k, v := range r.Prams {
			q.Set(k, v)
		}
		url.RawQuery = q.Encode()
	}

	// Create request
	req, err := http.NewRequest(r.Method, url.String(), bytes.NewReader(r.RequestBody))
	if err != nil {
		return nil, err
	}

	// Set headers
	for k, v := range r.RequestHeader {
		req.Header.Set(k, v)
	}
	return req, nil
}

func (r *Request) Send() error {

	// Send the request and measure latency
	req, err := r.StdRequest()
	if err != nil {
		return err
	}
	client := &http.Client{}
	start := time.Now()
	resp, err := client.Do(req)
	r.Latency = time.Since(start)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Save response status
	r.Status = resp.StatusCode

	// Save headers
	r.ResponseHeader = map[string]string{}
	for k, v := range resp.Header {
		if len(v) > 0 {
			r.ResponseHeader[k] = v[0]
		}
	}

	// Save body
	r.ResponseBody, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return nil
}

type FuzzFunc func(*Request)
