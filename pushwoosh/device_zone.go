// Copyright (c) 2013, Álvaro Vilanova Vidal
// Copyright (c) 2013, Stelapps (Appsales Dev S.L.)
// Use of this source code is governed by a BSD 2-Clause
// license that can be found in the LICENSE file.

package pushwoosh

import (
	"errors"
)

type ZoneResponse struct {
	Response
	Info struct {
		Name     string  `json:"name,omitempty"`
		Lat      float64 `json:"lat,omitempty"`
		Lng      float64 `json:"lng,omitempty"`
		Distance float64 `json:"distance,omitempty"`
	} `json:"response,omitempty"`
}

func (s DevicesService) NearestZone(device Identifiable, lat, lng float64) (*ZoneResponse, error) {
	if len(s.client.Application) <= 0 {
		return nil, errors.New("Application token is required")
	}
	if len(device.DeviceId()) <= 0 {
		return nil, errors.New("Device Hardware ID is required")
	}
	body := newDeviceNearestZoneBody(s.client.Application, device, lat, lng)
	req, err := s.client.NewRequest("POST", "/getNearestZone", body)
	if err != nil {
		return nil, err
	}
	resp := new(ZoneResponse)
	err = s.client.Do(req, resp)
	return resp, err
}

func newDeviceNearestZoneBody(app string, device Identifiable, lat, lng float64) interface{} {
	return struct {
		Application string  `json:"application"`
		HardwareId  string  `json:"hwid"`
		Lat         float64 `json:"lat"`
		Lng         float64 `json:"lng"`
	}{app, device.DeviceId(), lat, lng}
}
