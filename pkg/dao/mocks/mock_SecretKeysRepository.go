// Code generated by mockery v2.33.2. DO NOT EDIT.

package daomocks

import (
	context "context"
	ed25519 "crypto/ed25519"

	dao "github.com/a-novel/auth-service/pkg/dao"

	mock "github.com/stretchr/testify/mock"
)

// SecretKeysRepository is an autogenerated mock type for the SecretKeysRepository type
type SecretKeysRepository struct {
	mock.Mock
}

type SecretKeysRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *SecretKeysRepository) EXPECT() *SecretKeysRepository_Expecter {
	return &SecretKeysRepository_Expecter{mock: &_m.Mock}
}

// Delete provides a mock function with given fields: ctx, name
func (_m *SecretKeysRepository) Delete(ctx context.Context, name string) error {
	ret := _m.Called(ctx, name)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SecretKeysRepository_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type SecretKeysRepository_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx context.Context
//   - name string
func (_e *SecretKeysRepository_Expecter) Delete(ctx interface{}, name interface{}) *SecretKeysRepository_Delete_Call {
	return &SecretKeysRepository_Delete_Call{Call: _e.mock.On("Delete", ctx, name)}
}

func (_c *SecretKeysRepository_Delete_Call) Run(run func(ctx context.Context, name string)) *SecretKeysRepository_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *SecretKeysRepository_Delete_Call) Return(_a0 error) *SecretKeysRepository_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *SecretKeysRepository_Delete_Call) RunAndReturn(run func(context.Context, string) error) *SecretKeysRepository_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// List provides a mock function with given fields: ctx
func (_m *SecretKeysRepository) List(ctx context.Context) ([]*dao.SecretKeyModel, error) {
	ret := _m.Called(ctx)

	var r0 []*dao.SecretKeyModel
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]*dao.SecretKeyModel, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []*dao.SecretKeyModel); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*dao.SecretKeyModel)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SecretKeysRepository_List_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'List'
type SecretKeysRepository_List_Call struct {
	*mock.Call
}

// List is a helper method to define mock.On call
//   - ctx context.Context
func (_e *SecretKeysRepository_Expecter) List(ctx interface{}) *SecretKeysRepository_List_Call {
	return &SecretKeysRepository_List_Call{Call: _e.mock.On("List", ctx)}
}

func (_c *SecretKeysRepository_List_Call) Run(run func(ctx context.Context)) *SecretKeysRepository_List_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *SecretKeysRepository_List_Call) Return(_a0 []*dao.SecretKeyModel, _a1 error) *SecretKeysRepository_List_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SecretKeysRepository_List_Call) RunAndReturn(run func(context.Context) ([]*dao.SecretKeyModel, error)) *SecretKeysRepository_List_Call {
	_c.Call.Return(run)
	return _c
}

// Read provides a mock function with given fields: ctx, name
func (_m *SecretKeysRepository) Read(ctx context.Context, name string) (*dao.SecretKeyModel, error) {
	ret := _m.Called(ctx, name)

	var r0 *dao.SecretKeyModel
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*dao.SecretKeyModel, error)); ok {
		return rf(ctx, name)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *dao.SecretKeyModel); ok {
		r0 = rf(ctx, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dao.SecretKeyModel)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SecretKeysRepository_Read_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Read'
type SecretKeysRepository_Read_Call struct {
	*mock.Call
}

// Read is a helper method to define mock.On call
//   - ctx context.Context
//   - name string
func (_e *SecretKeysRepository_Expecter) Read(ctx interface{}, name interface{}) *SecretKeysRepository_Read_Call {
	return &SecretKeysRepository_Read_Call{Call: _e.mock.On("Read", ctx, name)}
}

func (_c *SecretKeysRepository_Read_Call) Run(run func(ctx context.Context, name string)) *SecretKeysRepository_Read_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *SecretKeysRepository_Read_Call) Return(_a0 *dao.SecretKeyModel, _a1 error) *SecretKeysRepository_Read_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SecretKeysRepository_Read_Call) RunAndReturn(run func(context.Context, string) (*dao.SecretKeyModel, error)) *SecretKeysRepository_Read_Call {
	_c.Call.Return(run)
	return _c
}

// Write provides a mock function with given fields: ctx, key, name
func (_m *SecretKeysRepository) Write(ctx context.Context, key ed25519.PrivateKey, name string) (*dao.SecretKeyModel, error) {
	ret := _m.Called(ctx, key, name)

	var r0 *dao.SecretKeyModel
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ed25519.PrivateKey, string) (*dao.SecretKeyModel, error)); ok {
		return rf(ctx, key, name)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ed25519.PrivateKey, string) *dao.SecretKeyModel); ok {
		r0 = rf(ctx, key, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dao.SecretKeyModel)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, ed25519.PrivateKey, string) error); ok {
		r1 = rf(ctx, key, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SecretKeysRepository_Write_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Write'
type SecretKeysRepository_Write_Call struct {
	*mock.Call
}

// Write is a helper method to define mock.On call
//   - ctx context.Context
//   - key ed25519.PrivateKey
//   - name string
func (_e *SecretKeysRepository_Expecter) Write(ctx interface{}, key interface{}, name interface{}) *SecretKeysRepository_Write_Call {
	return &SecretKeysRepository_Write_Call{Call: _e.mock.On("Write", ctx, key, name)}
}

func (_c *SecretKeysRepository_Write_Call) Run(run func(ctx context.Context, key ed25519.PrivateKey, name string)) *SecretKeysRepository_Write_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(ed25519.PrivateKey), args[2].(string))
	})
	return _c
}

func (_c *SecretKeysRepository_Write_Call) Return(_a0 *dao.SecretKeyModel, _a1 error) *SecretKeysRepository_Write_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SecretKeysRepository_Write_Call) RunAndReturn(run func(context.Context, ed25519.PrivateKey, string) (*dao.SecretKeyModel, error)) *SecretKeysRepository_Write_Call {
	_c.Call.Return(run)
	return _c
}

// NewSecretKeysRepository creates a new instance of SecretKeysRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSecretKeysRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *SecretKeysRepository {
	mock := &SecretKeysRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
