package mocks

import "github.com/hpcloud/cf-usb/lib/mgmt/cc_integration/httpclient"
import "github.com/stretchr/testify/mock"

type HttpClient struct {
	mock.Mock
}

func (_m *HttpClient) Request(request httpclient.Request) ([]byte, error) {
	ret := _m.Called(request)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(httpclient.Request) []byte); ok {
		r0 = rf(request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(httpclient.Request) error); ok {
		r1 = rf(request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
