// Code generated by mockery v2.20.0. DO NOT EDIT.

package servicesmocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// UpdateEmailService is an autogenerated mock type for the UpdateEmailService type
type UpdateEmailService struct {
	mock.Mock
}

type UpdateEmailService_Expecter struct {
	mock *mock.Mock
}

func (_m *UpdateEmailService) EXPECT() *UpdateEmailService_Expecter {
	return &UpdateEmailService_Expecter{mock: &_m.Mock}
}

// UpdateEmail provides a mock function with given fields: ctx, tokenRaw, newEmail, now
func (_m *UpdateEmailService) UpdateEmail(ctx context.Context, tokenRaw string, newEmail string, now time.Time) (func() error, error) {
	ret := _m.Called(ctx, tokenRaw, newEmail, now)

	var r0 func() error
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, time.Time) (func() error, error)); ok {
		return rf(ctx, tokenRaw, newEmail, now)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, time.Time) func() error); ok {
		r0 = rf(ctx, tokenRaw, newEmail, now)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(func() error)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, time.Time) error); ok {
		r1 = rf(ctx, tokenRaw, newEmail, now)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateEmailService_UpdateEmail_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateEmail'
type UpdateEmailService_UpdateEmail_Call struct {
	*mock.Call
}

// UpdateEmail is a helper method to define mock.On call
//   - ctx context.Context
//   - tokenRaw string
//   - newEmail string
//   - now time.Time
func (_e *UpdateEmailService_Expecter) UpdateEmail(ctx interface{}, tokenRaw interface{}, newEmail interface{}, now interface{}) *UpdateEmailService_UpdateEmail_Call {
	return &UpdateEmailService_UpdateEmail_Call{Call: _e.mock.On("UpdateEmail", ctx, tokenRaw, newEmail, now)}
}

func (_c *UpdateEmailService_UpdateEmail_Call) Run(run func(ctx context.Context, tokenRaw string, newEmail string, now time.Time)) *UpdateEmailService_UpdateEmail_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(time.Time))
	})
	return _c
}

func (_c *UpdateEmailService_UpdateEmail_Call) Return(_a0 func() error, _a1 error) *UpdateEmailService_UpdateEmail_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UpdateEmailService_UpdateEmail_Call) RunAndReturn(run func(context.Context, string, string, time.Time) (func() error, error)) *UpdateEmailService_UpdateEmail_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewUpdateEmailService interface {
	mock.TestingT
	Cleanup(func())
}

// NewUpdateEmailService creates a new instance of UpdateEmailService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUpdateEmailService(t mockConstructorTestingTNewUpdateEmailService) *UpdateEmailService {
	mock := &UpdateEmailService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}