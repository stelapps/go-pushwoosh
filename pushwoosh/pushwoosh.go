// Copyright (c) 2013, Álvaro Vilanova Vidal
// Copyright (c) 2013, Stelapps (Appsales Dev S.L.)
// Use of this source code is governed by a BSD 2-Clause
// license that can be found in the LICENSE file.

package pushwoosh

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"path"
)

const (
	Version          = "0.1.0"
	PushwooshVersion = "1.3"
)

type Client struct {
	Application string
	AuthToken   string
	// getaddrinfo may not be thread-safe on some systems. Enabling AddrInfo
	// cache ensures parallel bulk operations to work propertly without
	// getaddrinfo errors.
	// More info at https://code.google.com/p/go/issues/detail?id=3575
	CacheAddrInfo bool
	UserAgent     string

	Devices *DevicesService

	baseURL   *url.URL
	client    *http.Client
	endpoints []*url.URL
}

func (c *Client) BaseURL() *url.URL {
	return c.baseURL
}

func NewClient(client *http.Client) *Client {
	var c Client
	if client == nil {
		c.client = http.DefaultClient
	} else {
		c.client = client
	}
	baseURL, _ := url.Parse(defaultBaseURL())
	c.SetBaseURL(baseURL)
	c.UserAgent = defaultUserAgent()
	c.Devices = &DevicesService{&c}
	c.CacheAddrInfo = true
	return &c
}

func (c *Client) Do(req *http.Request, r interface{}) error {
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	err = json.NewDecoder(resp.Body).Decode(r)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		if rp, ok := r.(*Response); ok && err != nil {
			rp.Message = resp.Status
			rp.Response = resp
			rp.Status = resp.StatusCode
		}
		return errors.New(resp.Status)
	}

	rp, ok := r.(*Response)
	if err != nil {
		if resp.StatusCode != 200 && ok {
			rp.Message = resp.Status
			rp.Response = resp
			rp.Status = resp.StatusCode
			return errors.New(rp.Message)
		}
		return err
	}
	if ok {
		rp.Response = resp

		if rp.Status != 200 {
			return ErrorResponse(*rp)
		}
	}
	return nil
}

func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(path.Join(c.BaseURL().Path + urlStr))
	if err != nil {
		return nil, err
	}

	u := c.BaseURL().ResolveReference(rel)

	buffer := new(bytes.Buffer)
	if body != nil {
		err := json.NewEncoder(buffer).Encode(wrapRequestBody(body))
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buffer)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", c.UserAgent)
	c.CacheAddrInfo = true
	return req, nil
}

func (c *Client) SetBaseURL(url *url.URL) {
	c.baseURL = url
	c.endpoints = nil
}

func (c *Client) endpoint() *url.URL {
	if c.CacheAddrInfo {
		if len(c.endpoints) < 0 {
			endpoints, err := resolveHost(c.baseURL)
			if err != nil {
				return c.baseURL
			}
			c.endpoints = endpoints
		}
		return c.endpoints[rand.Intn(len(c.endpoints))]
	}
	return c.baseURL
}

type Response struct {
	Message  string         `json:"status_message"`
	Response *http.Response `json:"-"`
	Status   int            `json:"status_code"`
}

type ErrorResponse Response

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("(Code: %d) %s", e.Status, e.Message)
}

const (
	defaultBaseURLPattern   = "https://cp.pushwoosh.com/json/%s/"
	defaultUserAgentPattern = "go-pushwoosh/%s"
)

func defaultBaseURL() string {
	return fmt.Sprintf(defaultBaseURLPattern, PushwooshVersion)
}

func defaultUserAgent() string {
	return fmt.Sprintf(defaultUserAgentPattern, Version)
}

func resolveHost(fullURL *url.URL) ([]*url.URL, error) {
	ips, err := net.LookupIP(fullURL.Host)
	if err != nil {
		return []*url.URL{fullURL}, err
	}
	var urls []*url.URL
	for _, ip := range ips {
		ipUrl, err := url.Parse(fullURL.String())
		if err != nil {
			return []*url.URL{fullURL}, err
		}
		ipUrl.Host = ip.String()
		urls = append(urls, ipUrl)
	}
	return urls, nil
}

func wrapRequestBody(body interface{}) interface{} {
	return struct {
		Request interface{} `json:"request"`
	}{body}
}
