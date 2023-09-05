// Code generated by mockery v2.20.0. DO NOT EDIT.

package servicesmocks

import (
	context "context"

	models "github.com/a-novel/auth-service/pkg/models"
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// LoginService is an autogenerated mock type for the LoginService type
type LoginService struct {
	mock.Mock
}

type LoginService_Expecter struct {
	mock *mock.Mock
}

func (_m *LoginService) EXPECT() *LoginService_Expecter {
	return &LoginService_Expecter{mock: &_m.Mock}
}

// Login provides a mock function with given fields: ctx, email, password, now
func (_m *LoginService) Login(ctx context.Context, email string, password string, now time.Time) (*models.UserTokenStatus, error) {
	ret := _m.Called(ctx, email, password, now)

	var r0 *models.UserTokenStatus
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, time.Time) (*models.UserTokenStatus, error)); ok {
		return rf(ctx, email, password, now)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, time.Time) *models.UserTokenStatus); ok {
		r0 = rf(ctx, email, password, now)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.UserTokenStatus)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, time.Time) error); ok {
		r1 = rf(ctx, email, password, now)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LoginService_Login_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Login'
type LoginService_Login_Call struct {
	*mock.Call
}

// Login is a helper method to define mock.On call
//   - ctx context.Context
//   - email string
//   - password string
//   - now time.Time
func (_e *LoginService_Expecter) Login(ctx interface{}, email interface{}, password interface{}, now interface{}) *LoginService_Login_Call {
	return &LoginService_Login_Call{Call: _e.mock.On("Login", ctx, email, password, now)}
}

func (_c *LoginService_Login_Call) Run(run func(ctx context.Context, email string, password string, now time.Time)) *LoginService_Login_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(time.Time))
	})
	return _c
}

func (_c *LoginService_Login_Call) Return(_a0 *models.UserTokenStatus, _a1 error) *LoginService_Login_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *LoginService_Login_Call) RunAndReturn(run func(context.Context, string, string, time.Time) (*models.UserTokenStatus, error)) *LoginService_Login_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewLoginService interface {
	mock.TestingT
	Cleanup(func())
}

// NewLoginService creates a new instance of LoginService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewLoginService(t mockConstructorTestingTNewLoginService) *LoginService {
	mock := &LoginService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
