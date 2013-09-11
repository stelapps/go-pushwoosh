// Copyright (c) 2013, Álvaro Vilanova Vidal
// Copyright (c) 2013, Stelapps (Appsales Dev S.L.)
// Use of this source code is governed by a BSD 2-Clause
// license that can be found in the LICENSE file.

package pushwoosh

import (
	"errors"
	"reflect"
)

type TagsResponse struct {
	Response
	Info struct {
		Skipped []SkippedTag `json:"skipped,omitempty"`
	} `json:"response,omitempty"`
}

type SkippedTag struct {
	Tag    string `json:"tag"`
	Reason string `json:"reason"`
}

func (s DevicesService) SetTags(device Identifiable) (*TagsResponse, error) {
	if len(s.client.Application) <= 0 {
		return nil, errors.New("Application token is required")
	}
	if len(device.DeviceId()) <= 0 {
		return nil, errors.New("Device Hardware ID is required")
	}
	body, err := newSetDeviceTagsBody(s.client.Application, device)
	if err != nil {
		return nil, err
	}
	req, err := s.client.NewRequest("POST", "/setTags", body)
	if err != nil {
		return nil, err
	}
	resp := new(TagsResponse)
	err = s.client.Do(req, resp)
	return resp, err
}

func getTags(d interface{}) (map[string]interface{}, error) {
	v := reflect.ValueOf(d)
	t := v.Type()
	tags := map[string]interface{}{}
	if v.Kind() != reflect.Struct {
		return nil, errors.New("Tags can only be taken from a struct")
	}

	for i := 0; i < v.NumField(); i++ {
		value := v.Field(i)
		field := t.Field(i)
		name := field.Tag.Get("tag")
		if len(name) > 0 && value.CanInterface() {
			tags[name] = value.Interface()
		}
	}

	return tags, nil
}

func newSetDeviceTagsBody(app string, device Identifiable) (interface{}, error) {
	tags, err := getTags(device)
	if err != nil {
		return nil, err
	}
	return struct {
		Application string      `json:"application"`
		HardwareId  string      `json:"hwid"`
		Tags        interface{} `json:"tags"`
	}{app, device.DeviceId(), tags}, nil
}
