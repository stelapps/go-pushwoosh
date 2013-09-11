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

func TestDevicesService_PushStat(t *testing.T) {
	mux, server, client := sandbox()
	defer server.Close()

	mux.HandleFunc("/pushStat", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Request struct {
				Application string `json:"application"`
				HardwareId  string `json:"hwid"`
				Hash        string `json:"hash"`
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

		if body.Request.Hash != "myHash!" {
			t.Errorf("Hash = %s, want %s", body.Request.Hash, "myHash!")
		}

		res := Response{Status: 200, Message: "OK"}
		json.NewEncoder(w).Encode(res)
	})
	device := Device{HardwareId: "testHardwareId"}
	client.Application = "testAppToken"
	resp, err := client.Devices.PushStat(device, "myHash!")
	if err != nil {
		t.Errorf("Expected no error, found %s", err.Error())
	}
	want := Response{Status: 200, Message: "OK"}
	if !compareResponses(*resp, want) {
		t.Errorf("Response resp = %v, want %v", resp, want)
	}
}

func TestDevicesService_PushStat_invalidApp(t *testing.T) {
	client := NewClient(nil)
	device := Device{HardwareId: "testHardwareId"}
	_, err := client.Devices.PushStat(device, "myHash!")
	if err == nil {
		t.Errorf("Expected an error")
	}
}

func TestDevicesService_PushStat_invalidHardwareId(t *testing.T) {
	client := NewClient(nil)
	client.Application = "testAppToken"
	device := Device{}
	_, err := client.Devices.PushStat(device, "myHash!")
	if err == nil {
		t.Errorf("Expected an error")
	}
}

func TestDevicesService_PushStat_invalidHash(t *testing.T) {
	client := NewClient(nil)
	client.Application = "testAppToken"
	device := Device{HardwareId: "testHardwareId"}
	_, err := client.Devices.PushStat(device, "")
	if err == nil {
		t.Errorf("Expected an error")
	}
}
