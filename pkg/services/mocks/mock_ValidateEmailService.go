// Code generated by mockery v2.33.2. DO NOT EDIT.

package servicesmocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	time "time"

	uuid "github.com/google/uuid"
)

// ValidateEmailService is an autogenerated mock type for the ValidateEmailService type
type ValidateEmailService struct {
	mock.Mock
}

type ValidateEmailService_Expecter struct {
	mock *mock.Mock
}

func (_m *ValidateEmailService) EXPECT() *ValidateEmailService_Expecter {
	return &ValidateEmailService_Expecter{mock: &_m.Mock}
}

// ValidateEmail provides a mock function with given fields: ctx, id, code, now
func (_m *ValidateEmailService) ValidateEmail(ctx context.Context, id uuid.UUID, code string, now time.Time) error {
	ret := _m.Called(ctx, id, code, now)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, string, time.Time) error); ok {
		r0 = rf(ctx, id, code, now)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ValidateEmailService_ValidateEmail_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ValidateEmail'
type ValidateEmailService_ValidateEmail_Call struct {
	*mock.Call
}

// ValidateEmail is a helper method to define mock.On call
//   - ctx context.Context
//   - id uuid.UUID
//   - code string
//   - now time.Time
func (_e *ValidateEmailService_Expecter) ValidateEmail(ctx interface{}, id interface{}, code interface{}, now interface{}) *ValidateEmailService_ValidateEmail_Call {
	return &ValidateEmailService_ValidateEmail_Call{Call: _e.mock.On("ValidateEmail", ctx, id, code, now)}
}

func (_c *ValidateEmailService_ValidateEmail_Call) Run(run func(ctx context.Context, id uuid.UUID, code string, now time.Time)) *ValidateEmailService_ValidateEmail_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(string), args[3].(time.Time))
	})
	return _c
}

func (_c *ValidateEmailService_ValidateEmail_Call) Return(_a0 error) *ValidateEmailService_ValidateEmail_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ValidateEmailService_ValidateEmail_Call) RunAndReturn(run func(context.Context, uuid.UUID, string, time.Time) error) *ValidateEmailService_ValidateEmail_Call {
	_c.Call.Return(run)
	return _c
}

// NewValidateEmailService creates a new instance of ValidateEmailService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewValidateEmailService(t interface {
	mock.TestingT
	Cleanup(func())
}) *ValidateEmailService {
	mock := &ValidateEmailService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}