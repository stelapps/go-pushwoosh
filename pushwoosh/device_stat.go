// Copyright (c) 2013, Álvaro Vilanova Vidal
// Copyright (c) 2013, Stelapps (Appsales Dev S.L.)
// Use of this source code is governed by a BSD 2-Clause
// license that can be found in the LICENSE file.

package pushwoosh

import (
	"errors"
)

func (s DevicesService) PushStat(device Identifiable, hash string) (*Response, error) {
	if len(s.client.Application) <= 0 {
		return nil, errors.New("Application token is required")
	}
	if len(device.DeviceId()) <= 0 {
		return nil, errors.New("Device Hardware ID is required")
	}
	if len(hash) <= 0 {
		return nil, errors.New("Hash is required")
	}
	body := newDevicePushStatBody(s.client.Application, device, hash)
	req, err := s.client.NewRequest("POST", "/pushStat", body)
	if err != nil {
		return nil, err
	}
	resp := new(Response)
	err = s.client.Do(req, resp)
	return resp, err
}

func newDevicePushStatBody(app string, device Identifiable, hash string) interface{} {
	return struct {
		Application string `json:"application"`
		HardwareId  string `json:"hwid"`
		Hash        string `json:"hash"`
	}{app, device.DeviceId(), hash}
}
