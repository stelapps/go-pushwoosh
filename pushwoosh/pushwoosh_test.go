// Copyright (c) 2013, Álvaro Vilanova Vidal
// Copyright (c) 2013, Stelapps (Appsales Dev S.L.)
// Use of this source code is governed by a BSD 2-Clause
// license that can be found in the LICENSE file.

package pushwoosh

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func TestNewClient(t *testing.T) {
	c := NewClient(nil)
	if c.BaseURL().String() != defaultBaseURL() {
		t.Errorf("BaseURL = %v, want %v", c.BaseURL().String(), defaultBaseURL)
	}
	if c.UserAgent != defaultUserAgent() {
		t.Errorf("UserAgent = %v, want %v", c.UserAgent, defaultUserAgent())
	}
}

func TestNewRequest(t *testing.T) {
	c := NewClient(nil)

	inURL, outURL := "/foo", defaultBaseURL()+"foo"
	inBody := struct {
		TestField string `json:"testfield"`
	}{"test value"}
	outBody, _ := json.Marshal(wrapRequestBody(inBody))
	req, _ := c.NewRequest("POST", inURL, inBody)

	if req.URL.String() != outURL {
		t.Errorf("URL = %v, want %v", req.URL, outURL)
	}

	body, _ := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if bytes.Compare(bytes.Trim(body, " \n\r"), bytes.Trim(outBody, " \n\r")) != 0 {
		t.Errorf("Body = %s, want %s", body, outBody)
	}

	userAgent := req.Header.Get("User-Agent")
	if c.UserAgent != defaultUserAgent() {
		t.Errorf("User-Agent = %v, want %v", userAgent, c.UserAgent)
	}
}

func TestNewRequest_invalidJSON(t *testing.T) {
	c := NewClient(nil)
	_, err := c.NewRequest("GET", "/", &struct{ InvalidField map[int]int }{})

	if err == nil {
		t.Error("Expected error to be returned.")
	}
	if err, ok := err.(*json.UnsupportedTypeError); !ok {
		t.Errorf("Expected a JSON error; got %#v.", err)
	}
}

func TestDo(t *testing.T) {
	mux, server, client := sandbox()
	defer server.Close()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Request method = %v, want POST", r.Method)
		}
		json.NewEncoder(w).Encode(Response{Message: "foo", Status: 200})
	})

	req, _ := client.NewRequest("POST", "/", nil)
	var resp Response
	err := client.Do(req, &resp)

	if err != nil {
		t.Errorf("Expected no error, found %s", err.Error())
	}

	if resp.Response == nil {
		t.Errorf("Expected a response reference, found %v", resp.Response)
	}

	want := Response{Message: "foo", Status: 200}
	if !compareResponses(resp, want) {
		t.Errorf("Response resp = %v, want %v", resp, want)
	}
}

func TestDo_httpError(t *testing.T) {
	mux, server, client := sandbox()
	defer server.Close()
	message := "foo"

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, message, http.StatusInternalServerError)
	})

	req, _ := client.NewRequest("POST", "/", nil)
	var resp Response
	err := client.Do(req, &resp)

	if err == nil {
		t.Errorf("Expected an error")
	}

	if resp.Status != http.StatusInternalServerError || resp.Message == message {
		t.Errorf("Response resp = %v, want Status = %d and Message != %s",
			resp, http.StatusInternalServerError, message)
	}
}

func TestDo_httpStatusError(t *testing.T) {
	mux, server, client := sandbox()
	defer server.Close()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body, _ := json.Marshal(Response{Status: 514, Message: "foo"})
		http.Error(w, bytes.NewBuffer(body).String(), http.StatusInternalServerError)
	})

	req, _ := client.NewRequest("POST", "/", nil)
	var resp Response
	err := client.Do(req, &resp)

	if err == nil {
		t.Errorf("Expected an error")
	}

	want := Response{Status: 514, Message: "foo"}
	if !compareResponses(resp, want) {
		t.Errorf("Response resp = %v, want = %v", resp, want)
	}
}

func TestDo_invalidJSON(t *testing.T) {
	mux, server, client := sandbox()
	defer server.Close()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		foo := struct{ Bar map[int]int }{}
		json.NewEncoder(w).Encode(foo)
	})

	req, _ := client.NewRequest("POST", "/", nil)
	var resp Response
	err := client.Do(req, &resp)

	if err == nil {
		t.Errorf("Expected an error")
	}
	var want Response
	if !reflect.DeepEqual(resp, want) {
		t.Errorf("Response resp = %v, want %v", resp, want)
	}
}

func compareResponses(a, b Response) bool {
	return a.Message == b.Message && a.Status == b.Status
}

func sandbox() (*http.ServeMux, *httptest.Server, *Client) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	client := NewClient(nil)

	u, _ := url.Parse(server.URL)
	client.CacheAddrInfo = false
	client.SetBaseURL(u)

	return mux, server, client
}
