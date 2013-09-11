// Copyright (c) 2013, Álvaro Vilanova Vidal
// Copyright (c) 2013, Stelapps (Appsales Dev S.L.)
// Use of this source code is governed by a BSD 2-Clause
// license that can be found in the LICENSE file.

package pushwoosh

import (
	"errors"
)

type Device struct {
	HardwareId string
	Language   string
	PushToken  string
	TimeZone   int
	Type       DeviceType
}

type DeviceType int

const (
	IOS          DeviceType = 1
	BlackBerry   DeviceType = 2
	Android      DeviceType = 3
	Nokia        DeviceType = 4
	WindowsPhone DeviceType = 5
	OSX          DeviceType = 7
)

type Identifiable interface {
	DeviceId() string
}

type LanguageRegistrable interface {
	DeviceLanguage() string
}

type Registrable interface {
	Identifiable
	DevicePushToken() string
	DeviceType() DeviceType
}

type TimeZoneRegistrable interface {
	DeviceTimeZone() int
}

func (c *Client) RegisterDevice(d Registrable) error {
	return nil
}

func (d Device) DeviceId() string {
	return d.HardwareId
}

func (d Device) DeviceLanguage() string {
	return d.Language
}

func (d Device) DevicePushToken() string {
	return d.PushToken
}

func (d Device) DeviceTimeZone() int {
	return d.TimeZone
}

func (d Device) DeviceType() DeviceType {
	return d.Type
}

type DevicesService struct {
	client *Client
}

func (s DevicesService) Register(device Registrable) (*Response, error) {
	if len(s.client.Application) <= 0 {
		return nil, errors.New("Application token is required")
	}
	if err := checkDevice(device); err != nil {
		return nil, err
	}
	body := newRegisterDeviceBody(s.client.Application, device)
	req, err := s.client.NewRequest("POST", "/registerDevice", body)
	if err != nil {
		return nil, err
	}
	resp := new(Response)
	err = s.client.Do(req, resp)
	return resp, err
}

func (s DevicesService) Unregister(hardwareId string) (*Response, error) {
	if len(s.client.Application) <= 0 {
		return nil, errors.New("Application token is required")
	}
	if len(hardwareId) <= 0 {
		return nil, errors.New("Device Hardware ID is required")
	}

	body := struct {
		Application string `json:"application"`
		HardwareId  string `json:"hwid"`
	}{s.client.Application, hardwareId}
	req, err := s.client.NewRequest("POST", "/unregisterDevice", body)
	if err != nil {
		return nil, err
	}
	resp := new(Response)
	err = s.client.Do(req, resp)
	return resp, err
}

func checkDevice(device Registrable) error {
	if len(device.DeviceId()) <= 0 {
		return errors.New("Device Hardware ID is required")
	}
	if len(device.DevicePushToken()) <= 0 {
		return errors.New("Device Push Token is required")
	}
	if device.DeviceType() == 0 {
		return errors.New("Device Type is required")
	}
	return nil
}

func newRegisterDeviceBody(app string, device Registrable) interface{} {
	body := struct {
		Application string     `json:"application"`
		PushToken   string     `json:"push_token"`
		Language    string     `json:"language,omitempty"`
		HardwareId  string     `json:"hwid"`
		TimeZone    int        `json:"timezone,omitempty"`
		Type        DeviceType `json:"device_type"`
	}{}
	body.Application = app
	body.HardwareId = device.DeviceId()
	body.PushToken = device.DevicePushToken()
	body.Type = device.DeviceType()
	timeZoneAspect, ok := device.(TimeZoneRegistrable)
	if ok {
		body.TimeZone = timeZoneAspect.DeviceTimeZone()
	}
	langugageAspect, ok := device.(LanguageRegistrable)
	if ok {
		body.Language = langugageAspect.DeviceLanguage()
	}
	return body
}
