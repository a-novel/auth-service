// Code generated by mockery v2.20.0. DO NOT EDIT.

package servicesmocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// EmailExistsService is an autogenerated mock type for the EmailExistsService type
type EmailExistsService struct {
	mock.Mock
}

type EmailExistsService_Expecter struct {
	mock *mock.Mock
}

func (_m *EmailExistsService) EXPECT() *EmailExistsService_Expecter {
	return &EmailExistsService_Expecter{mock: &_m.Mock}
}

// EmailExists provides a mock function with given fields: ctx, email
func (_m *EmailExistsService) EmailExists(ctx context.Context, email string) (bool, error) {
	ret := _m.Called(ctx, email)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (bool, error)); ok {
		return rf(ctx, email)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) bool); ok {
		r0 = rf(ctx, email)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EmailExistsService_EmailExists_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EmailExists'
type EmailExistsService_EmailExists_Call struct {
	*mock.Call
}

// EmailExists is a helper method to define mock.On call
//   - ctx context.Context
//   - email string
func (_e *EmailExistsService_Expecter) EmailExists(ctx interface{}, email interface{}) *EmailExistsService_EmailExists_Call {
	return &EmailExistsService_EmailExists_Call{Call: _e.mock.On("EmailExists", ctx, email)}
}

func (_c *EmailExistsService_EmailExists_Call) Run(run func(ctx context.Context, email string)) *EmailExistsService_EmailExists_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *EmailExistsService_EmailExists_Call) Return(_a0 bool, _a1 error) *EmailExistsService_EmailExists_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *EmailExistsService_EmailExists_Call) RunAndReturn(run func(context.Context, string) (bool, error)) *EmailExistsService_EmailExists_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewEmailExistsService interface {
	mock.TestingT
	Cleanup(func())
}

// NewEmailExistsService creates a new instance of EmailExistsService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewEmailExistsService(t mockConstructorTestingTNewEmailExistsService) *EmailExistsService {
	mock := &EmailExistsService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
