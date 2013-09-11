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

func TestDevicesService_NearestZone(t *testing.T) {
	mux, server, client := sandbox()
	defer server.Close()

	mux.HandleFunc("/getNearestZone", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Request struct {
				Application string  `json:"application"`
				HardwareId  string  `json:"hwid"`
				Lat         float64 `json:"lat"`
				Lng         float64 `json:"lng"`
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

		if body.Request.Lat != 3.14 {
			t.Errorf("Lat = %f, want %f", body.Request.Lat, 3.14)
		}

		if body.Request.Lng != 299792458.0 {
			t.Errorf("Lng = %f, want %f", body.Request.Lat, 299792458.0)
		}
		var res ZoneResponse
		res.Status = 200
		res.Message = "OK"
		res.Info.Name = "neverhood"
		res.Info.Lat = 10.1110
		res.Info.Lng = 11.1011
		res.Info.Distance = 90000.0
		json.NewEncoder(w).Encode(res)
	})

	device := Device{HardwareId: "testHardwareId"}

	client.Application = "testAppToken"
	resp, err := client.Devices.NearestZone(device, 3.14, 299792458.0)
	if err != nil {
		t.Errorf("Expected no error, found %s", err.Error())
	}
	want := Response{Status: 200, Message: "OK"}
	if !compareResponses(resp.Response, want) {
		t.Errorf("Response resp = %v, want %v", resp, want)
	}
	if resp.Info.Name != "neverhood" {
		t.Errorf("Response name = %s, want %s", resp.Info.Name, "neverhood")
	}
	if resp.Info.Lat != 10.1110 {
		t.Errorf("Response name = %f, want %f", resp.Info.Lat, 10.1110)
	}
	if resp.Info.Lng != 11.1011 {
		t.Errorf("Response name = %f, want %f", resp.Info.Lng, 11.1011)
	}
	if resp.Info.Distance != 90000.0 {
		t.Errorf("Response name = %f, want %f", resp.Info.Distance, 90000.0)
	}
}

func TestDevicesService_NearestZone_invalidApp(t *testing.T) {
	client := NewClient(nil)
	device := deviceTagsTest{
		HardwareId: "testHardwareId",
	}
	_, err := client.Devices.NearestZone(device, 1.0, 2.0)
	if err == nil {
		t.Errorf("Expected an error")
	}
}

func TestDevicesService_NearestZone_invalidHardwareId(t *testing.T) {
	client := NewClient(nil)
	client.Application = "testAppToken"
	device := deviceTagsTest{}
	_, err := client.Devices.NearestZone(device, 1.0, 2.0)
	if err == nil {
		t.Errorf("Expected an error")
	}
}
