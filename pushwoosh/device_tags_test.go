// Copyright (c) 2013, Álvaro Vilanova Vidal
// Copyright (c) 2013, Stelapps (Appsales Dev S.L.)
// Use of this source code is governed by a BSD 2-Clause
// license that can be found in the LICENSE file.

package pushwoosh

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
)

func TestDevicesService_SetTags(t *testing.T) {
	mux, server, client := sandbox()
	defer server.Close()

	device := deviceTagsTest{
		HardwareId: "testHardwareId",
		Foo:        14,
		Bar:        "fourteen",
	}

	mux.HandleFunc("/setTags", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Request struct {
				Application string                 `json:"application"`
				HardwareId  string                 `json:"hwid"`
				Tags        map[string]interface{} `json:"tags"`
			} `json:"request"`
		}

		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			t.Errorf("Expected no error, found %s", err.Error())
		}

		if body.Request.Application != "testAppToken" {
			t.Errorf("Application = %s, want %s", body.Request.Application, "testAppToken")
		}

		if body.Request.HardwareId != device.HardwareId {
			t.Errorf("HardwareId = %s, want %s", body.Request.HardwareId, device.HardwareId)
		}

		v, ok := body.Request.Tags["foo!"]
		if !ok {
			t.Error(`Expected tag "foo!". Not found`)
		}
		if v != 14.0 {
			t.Errorf(`Tags["foo!"] = %v, want %f`, v, 14.0)
		}
		v, ok = body.Request.Tags["bar"]
		if !ok {
			t.Error(`Expected tag "bar". Not found`)
		}
		if v != "fourteen" {
			t.Errorf(`Tags["bar"] = %v, want %s`, v, "fourteen")
		}
		res := Response{Status: 200, Message: "OK"}
		json.NewEncoder(w).Encode(res)
	})

	client.Application = "testAppToken"
	resp, err := client.Devices.SetTags(device)
	if err != nil {
		t.Errorf("Expected no error, found %s", err.Error())
	}
	want := Response{Status: 200, Message: "OK"}
	if !compareResponses(resp.Response, want) {
		t.Errorf("Response resp = %v, want %v", resp, want)
	}
}

func TestDevicesService_SetTags_invalidApp(t *testing.T) {
	client := NewClient(nil)
	device := deviceTagsTest{
		HardwareId: "testHardwareId",
	}
	_, err := client.Devices.SetTags(device)
	if err == nil {
		t.Errorf("Expected an error")
	}
}

func TestDevicesService_SetTags_invalidHardwareId(t *testing.T) {
	client := NewClient(nil)
	client.Application = "testAppToken"
	device := deviceTagsTest{}
	_, err := client.Devices.SetTags(device)
	if err == nil {
		t.Errorf("Expected an error")
	}
}

func TestDevicesService_SetTags_skippedTags(t *testing.T) {
	mux, server, client := sandbox()
	defer server.Close()

	device := deviceTagsTest{
		HardwareId: "testHardwareId",
		Foo:        14,
		Bar:        "fourteen",
	}

	mux.HandleFunc("/setTags", func(w http.ResponseWriter, r *http.Request) {
		var res TagsResponse
		res.Status = 200
		res.Message = "OK"
		res.Info.Skipped = []SkippedTag{SkippedTag{
			Tag:    "foo!",
			Reason: "Invalid name",
		}}
		json.NewEncoder(w).Encode(res)
	})

	client.Application = "testAppToken"
	resp, err := client.Devices.SetTags(device)
	if err != nil {
		t.Errorf("Expected no error, found %s", err.Error())
	}
	want := Response{Status: 200, Message: "OK"}
	if !compareResponses(resp.Response, want) {
		t.Errorf("Response resp = %v, want %v", resp, want)
	}
	swant := []SkippedTag{SkippedTag{Tag: "foo!", Reason: "Invalid name"}}
	if !reflect.DeepEqual(resp.Info.Skipped, swant) {
		t.Errorf("Response.Skipped resp = %v, want %v", resp.Info.Skipped, swant)
	}
}

type deviceTagsTest struct {
	HardwareId string
	Foo        int    `tag:"foo!"`
	Bar        string `tag:"bar"`
}

func (device deviceTagsTest) DeviceId() string {
	return device.HardwareId
}
