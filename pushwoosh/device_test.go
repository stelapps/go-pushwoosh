// Copyright (c) 2013, Álvaro Vilanova Vidal
// Copyright (c) 2013, Stelapps (Appsales Dev S.L.)
// Use of this source code is governed by a BSD 2-Clause
// license that can be found in the LICENSE file.

package pushwoosh

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestDevicesService_Register(t *testing.T) {
	mux, server, client := sandbox()
	defer server.Close()

	device := Device{
		HardwareId: "testHardwareId",
		Language:   "en",
		PushToken:  "testPushToken",
		TimeZone:   3600,
		Type:       IOS,
	}

	mux.HandleFunc("/registerDevice", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Request struct {
				Application string     `json:"application"`
				PushToken   string     `json:"push_token"`
				Language    string     `json:"language,omitempty"`
				HardwareId  string     `json:"hwid"`
				TimeZone    int        `json:"timezone,omitempty"`
				Type        DeviceType `json:"device_type"`
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

		if body.Request.Language != device.Language {
			t.Errorf("Language = %s, want %s", body.Request.Language, device.Language)
		}

		if body.Request.PushToken != device.PushToken {
			t.Errorf("PushToken = %s, want %s", body.Request.PushToken, device.PushToken)
		}

		if body.Request.TimeZone != device.TimeZone {
			t.Errorf("TimeZone = %s, want %s", body.Request.TimeZone, device.TimeZone)
		}

		if body.Request.Type != device.Type {
			t.Errorf("Type = %s, want %s", body.Request.Type, device.Type)
		}

		res := Response{Status: 200, Message: "OK"}
		json.NewEncoder(w).Encode(res)
	})

	client.Application = "testAppToken"
	resp, err := client.Devices.Register(device)
	if err != nil {
		t.Errorf("Expected no error, found %s", err.Error())
	}
	want := Response{Status: 200, Message: "OK"}
	if !compareResponses(*resp, want) {
		t.Errorf("Response resp = %v, want %v", resp, want)
	}
}

func TestDevicesService_Register_invalidApp(t *testing.T) {
	client := NewClient(nil)
	device := Device{
		HardwareId: "testHardwareId",
		Language:   "en",
		PushToken:  "testPushToken",
		TimeZone:   0,
		Type:       IOS,
	}
	_, err := client.Devices.Register(device)
	if err == nil {
		t.Errorf("Expected an error")
	}
}

func TestDevicesService_Register_invalidHardwareId(t *testing.T) {
	client := NewClient(nil)
	client.Application = "testAppToken"
	device := Device{
		Language:  "en",
		PushToken: "testPushToken",
		TimeZone:  0,
		Type:      IOS,
	}
	_, err := client.Devices.Register(device)
	if err == nil {
		t.Errorf("Expected an error")
	}
}

func TestDevicesService_Register_invalidPushToken(t *testing.T) {
	client := NewClient(nil)
	client.Application = "testAppToken"
	device := Device{
		HardwareId: "testHardwareId",
		Language:   "en",
		TimeZone:   0,
		Type:       IOS,
	}
	_, err := client.Devices.Register(device)
	if err == nil {
		t.Errorf("Expected an error")
	}
}

func TestDevicesService_Register_invalidDeviceToken(t *testing.T) {
	client := NewClient(nil)
	client.Application = "testAppToken"
	device := Device{
		HardwareId: "testHardwareId",
		Language:   "en",
		PushToken:  "testPushToken",
		TimeZone:   0,
	}
	_, err := client.Devices.Register(device)
	if err == nil {
		t.Errorf("Expected an error")
	}
}

func TestDevicesService_Register_error(t *testing.T) {
	mux, server, client := sandbox()
	defer server.Close()

	mux.HandleFunc("/registerDevice", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(Response{Status: 210, Message: "foo"})
	})

	device := Device{
		HardwareId: "testHardwareId",
		PushToken:  "testPushToken",
		Type:       IOS,
	}
	client.Application = "testAppToken"
	resp, err := client.Devices.Register(device)
	if err == nil {
		t.Errorf("Expected an error")
	}
	want := Response{Status: 210, Message: "foo"}
	if !compareResponses(*resp, want) {
		t.Errorf("Response resp = %v, want %v", resp, want)
	}
}

func TestDevicesService_Unregister(t *testing.T) {
	mux, server, client := sandbox()
	defer server.Close()

	mux.HandleFunc("/unregisterDevice", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Request struct {
				Application string `json:"application"`
				HardwareId  string `json:"hwid"`
			} `json:"request"`
		}

		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			t.Errorf("Expected no error, found %s", err.Error())
		}

		if body.Request.Application != "foo" {
			t.Errorf("HardwareId = %s, want %s", body.Request.HardwareId, "foo")
		}

		if body.Request.HardwareId != "bar" {
			t.Errorf("HardwareId = %s, want %s", body.Request.HardwareId, "bar")
		}

		res := Response{Status: 200, Message: "OK"}
		json.NewEncoder(w).Encode(res)
	})

	client.Application = "foo"
	resp, err := client.Devices.Unregister("bar")
	if err != nil {
		t.Errorf("Expected no error, found %s", err.Error())
	}
	want := Response{Status: 200, Message: "OK"}
	if !compareResponses(*resp, want) {
		t.Errorf("Response resp = %v, want %v", resp, want)
	}
}

func TestDevicesService_Unregister_invalidApp(t *testing.T) {
	client := NewClient(nil)
	_, err := client.Devices.Unregister("bar")
	if err == nil {
		t.Errorf("Expected an error")
	}
}

func TestDevicesService_Unregister_invalidHardwareId(t *testing.T) {
	client := NewClient(nil)
	client.Application = "foo"
	_, err := client.Devices.Unregister("")
	if err == nil {
		t.Errorf("Expected an error")
	}
}

func TestDevicesService_Unegister_error(t *testing.T) {
	mux, server, client := sandbox()
	defer server.Close()

	mux.HandleFunc("/unregisterDevice", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(Response{Status: 210, Message: "foo"})
	})

	client.Application = "testAppToken"
	resp, err := client.Devices.Unregister("testHardwareId")
	if err == nil {
		t.Errorf("Expected an error")
	}
	want := Response{Status: 210, Message: "foo"}
	if !compareResponses(*resp, want) {
		t.Errorf("Response resp = %v, want %v", resp, want)
	}
}
