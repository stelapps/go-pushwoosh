// Copyright (c) 2013, Álvaro Vilanova Vidal
// Copyright (c) 2013, Stelapps (Appsales Dev S.L.)
// Use of this source code is governed by a BSD 2-Clause
// license that can be found in the LICENSE file.

package pushwoosh

import (
	"errors"
)

func (s DevicesService) SetBadge(device Identifiable, value int) (*Response, error) {
	if len(s.client.Application) <= 0 {
		return nil, errors.New("Application token is required")
	}
	if len(device.DeviceId()) <= 0 {
		return nil, errors.New("Device Hardware ID is required")
	}
	body := newSetDeviceBadgeBody(s.client.Application, device, value)
	req, err := s.client.NewRequest("POST", "/setBadge", body)
	if err != nil {
		return nil, err
	}
	resp := new(Response)
	err = s.client.Do(req, resp)
	return resp, err
}

func newSetDeviceBadgeBody(app string, device Identifiable, value int) interface{} {
	return struct {
		Application string `json:"application"`
		HardwareId  string `json:"hwid"`
		Badge       int    `json:"badge"`
	}{app, device.DeviceId(), value}
}
