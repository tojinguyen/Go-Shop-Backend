// Code generated by mockery v2.53.4. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// EmailService is an autogenerated mock type for the EmailService type
type EmailService struct {
	mock.Mock
}

type EmailService_Expecter struct {
	mock *mock.Mock
}

func (_m *EmailService) EXPECT() *EmailService_Expecter {
	return &EmailService_Expecter{mock: &_m.Mock}
}

// SendEmail provides a mock function with given fields: to, subject, body
func (_m *EmailService) SendEmail(to []string, subject string, body string) error {
	ret := _m.Called(to, subject, body)

	if len(ret) == 0 {
		panic("no return value specified for SendEmail")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func([]string, string, string) error); ok {
		r0 = rf(to, subject, body)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// EmailService_SendEmail_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SendEmail'
type EmailService_SendEmail_Call struct {
	*mock.Call
}

// SendEmail is a helper method to define mock.On call
//   - to []string
//   - subject string
//   - body string
func (_e *EmailService_Expecter) SendEmail(to interface{}, subject interface{}, body interface{}) *EmailService_SendEmail_Call {
	return &EmailService_SendEmail_Call{Call: _e.mock.On("SendEmail", to, subject, body)}
}

func (_c *EmailService_SendEmail_Call) Run(run func(to []string, subject string, body string)) *EmailService_SendEmail_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]string), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *EmailService_SendEmail_Call) Return(_a0 error) *EmailService_SendEmail_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *EmailService_SendEmail_Call) RunAndReturn(run func([]string, string, string) error) *EmailService_SendEmail_Call {
	_c.Call.Return(run)
	return _c
}

// SendHTMLEmail provides a mock function with given fields: to, subject, htmlBody
func (_m *EmailService) SendHTMLEmail(to []string, subject string, htmlBody string) error {
	ret := _m.Called(to, subject, htmlBody)

	if len(ret) == 0 {
		panic("no return value specified for SendHTMLEmail")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func([]string, string, string) error); ok {
		r0 = rf(to, subject, htmlBody)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// EmailService_SendHTMLEmail_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SendHTMLEmail'
type EmailService_SendHTMLEmail_Call struct {
	*mock.Call
}

// SendHTMLEmail is a helper method to define mock.On call
//   - to []string
//   - subject string
//   - htmlBody string
func (_e *EmailService_Expecter) SendHTMLEmail(to interface{}, subject interface{}, htmlBody interface{}) *EmailService_SendHTMLEmail_Call {
	return &EmailService_SendHTMLEmail_Call{Call: _e.mock.On("SendHTMLEmail", to, subject, htmlBody)}
}

func (_c *EmailService_SendHTMLEmail_Call) Run(run func(to []string, subject string, htmlBody string)) *EmailService_SendHTMLEmail_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]string), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *EmailService_SendHTMLEmail_Call) Return(_a0 error) *EmailService_SendHTMLEmail_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *EmailService_SendHTMLEmail_Call) RunAndReturn(run func([]string, string, string) error) *EmailService_SendHTMLEmail_Call {
	_c.Call.Return(run)
	return _c
}

// SendTemplateEmail provides a mock function with given fields: to, subject, templateName, data
func (_m *EmailService) SendTemplateEmail(to []string, subject string, templateName string, data interface{}) error {
	ret := _m.Called(to, subject, templateName, data)

	if len(ret) == 0 {
		panic("no return value specified for SendTemplateEmail")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func([]string, string, string, interface{}) error); ok {
		r0 = rf(to, subject, templateName, data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// EmailService_SendTemplateEmail_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SendTemplateEmail'
type EmailService_SendTemplateEmail_Call struct {
	*mock.Call
}

// SendTemplateEmail is a helper method to define mock.On call
//   - to []string
//   - subject string
//   - templateName string
//   - data interface{}
func (_e *EmailService_Expecter) SendTemplateEmail(to interface{}, subject interface{}, templateName interface{}, data interface{}) *EmailService_SendTemplateEmail_Call {
	return &EmailService_SendTemplateEmail_Call{Call: _e.mock.On("SendTemplateEmail", to, subject, templateName, data)}
}

func (_c *EmailService_SendTemplateEmail_Call) Run(run func(to []string, subject string, templateName string, data interface{})) *EmailService_SendTemplateEmail_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]string), args[1].(string), args[2].(string), args[3].(interface{}))
	})
	return _c
}

func (_c *EmailService_SendTemplateEmail_Call) Return(_a0 error) *EmailService_SendTemplateEmail_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *EmailService_SendTemplateEmail_Call) RunAndReturn(run func([]string, string, string, interface{}) error) *EmailService_SendTemplateEmail_Call {
	_c.Call.Return(run)
	return _c
}

// NewEmailService creates a new instance of EmailService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewEmailService(t interface {
	mock.TestingT
	Cleanup(func())
}) *EmailService {
	mock := &EmailService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
