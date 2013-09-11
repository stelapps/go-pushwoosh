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

func TestDevicesService_SetBadge(t *testing.T) {
	mux, server, client := sandbox()
	defer server.Close()

	mux.HandleFunc("/setBadge", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Request struct {
				Application string `json:"application"`
				HardwareId  string `json:"hwid"`
				Badge       int    `json:"badge"`
			} `json:"request"`
		}

		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			t.Errorf("Expected no error, found %s", err.Error())
		}

		if body.Request.Application != "testAppToken" {
			t.Errorf("Application = %s, want %s", body.Request.Application, "testAppToken")
		}

		if body.Request.HardwareId != "testHardwareId" {
			t.Errorf("HardwareId = %s, want %s", body.Request.HardwareId, "testHardwareId")
		}

		if body.Request.Badge != 14 {
			t.Errorf("Badge = %d, want %d", body.Request.Badge, 14)
		}

		res := Response{Status: 200, Message: "OK"}
		json.NewEncoder(w).Encode(res)
	})
	device := Device{HardwareId: "testHardwareId"}
	client.Application = "testAppToken"
	resp, err := client.Devices.SetBadge(device, 14)
	if err != nil {
		t.Errorf("Expected no error, found %s", err.Error())
	}
	want := Response{Status: 200, Message: "OK"}
	if !compareResponses(*resp, want) {
		t.Errorf("Response resp = %v, want %v", resp, want)
	}
}

func TestDevicesService_SetBadge_invalidApp(t *testing.T) {
	client := NewClient(nil)
	device := Device{HardwareId: "testHardwareId"}
	_, err := client.Devices.SetBadge(device, 0)
	if err == nil {
		t.Errorf("Expected an error")
	}
}

func TestDevicesService_SetBadge_invalidHardwareId(t *testing.T) {
	client := NewClient(nil)
	client.Application = "testAppToken"
	device := Device{}
	_, err := client.Devices.SetBadge(device, 0)
	if err == nil {
		t.Errorf("Expected an error")
	}
}
