// Code generated by mockery v2.20.0. DO NOT EDIT.

package daomocks

import (
	context "context"

	dao "github.com/a-novel/auth-service/pkg/dao"
	mock "github.com/stretchr/testify/mock"

	time "time"

	uuid "github.com/google/uuid"
)

// ProfileRepository is an autogenerated mock type for the ProfileRepository type
type ProfileRepository struct {
	mock.Mock
}

type ProfileRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *ProfileRepository) EXPECT() *ProfileRepository_Expecter {
	return &ProfileRepository_Expecter{mock: &_m.Mock}
}

// GetProfile provides a mock function with given fields: ctx, id
func (_m *ProfileRepository) GetProfile(ctx context.Context, id uuid.UUID) (*dao.ProfileModel, error) {
	ret := _m.Called(ctx, id)

	var r0 *dao.ProfileModel
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (*dao.ProfileModel, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) *dao.ProfileModel); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dao.ProfileModel)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProfileRepository_GetProfile_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetProfile'
type ProfileRepository_GetProfile_Call struct {
	*mock.Call
}

// GetProfile is a helper method to define mock.On call
//   - ctx context.Context
//   - id uuid.UUID
func (_e *ProfileRepository_Expecter) GetProfile(ctx interface{}, id interface{}) *ProfileRepository_GetProfile_Call {
	return &ProfileRepository_GetProfile_Call{Call: _e.mock.On("GetProfile", ctx, id)}
}

func (_c *ProfileRepository_GetProfile_Call) Run(run func(ctx context.Context, id uuid.UUID)) *ProfileRepository_GetProfile_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *ProfileRepository_GetProfile_Call) Return(_a0 *dao.ProfileModel, _a1 error) *ProfileRepository_GetProfile_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ProfileRepository_GetProfile_Call) RunAndReturn(run func(context.Context, uuid.UUID) (*dao.ProfileModel, error)) *ProfileRepository_GetProfile_Call {
	_c.Call.Return(run)
	return _c
}

// GetProfileBySlug provides a mock function with given fields: ctx, slug
func (_m *ProfileRepository) GetProfileBySlug(ctx context.Context, slug string) (*dao.ProfileModel, error) {
	ret := _m.Called(ctx, slug)

	var r0 *dao.ProfileModel
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*dao.ProfileModel, error)); ok {
		return rf(ctx, slug)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *dao.ProfileModel); ok {
		r0 = rf(ctx, slug)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dao.ProfileModel)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, slug)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProfileRepository_GetProfileBySlug_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetProfileBySlug'
type ProfileRepository_GetProfileBySlug_Call struct {
	*mock.Call
}

// GetProfileBySlug is a helper method to define mock.On call
//   - ctx context.Context
//   - slug string
func (_e *ProfileRepository_Expecter) GetProfileBySlug(ctx interface{}, slug interface{}) *ProfileRepository_GetProfileBySlug_Call {
	return &ProfileRepository_GetProfileBySlug_Call{Call: _e.mock.On("GetProfileBySlug", ctx, slug)}
}

func (_c *ProfileRepository_GetProfileBySlug_Call) Run(run func(ctx context.Context, slug string)) *ProfileRepository_GetProfileBySlug_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *ProfileRepository_GetProfileBySlug_Call) Return(_a0 *dao.ProfileModel, _a1 error) *ProfileRepository_GetProfileBySlug_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ProfileRepository_GetProfileBySlug_Call) RunAndReturn(run func(context.Context, string) (*dao.ProfileModel, error)) *ProfileRepository_GetProfileBySlug_Call {
	_c.Call.Return(run)
	return _c
}

// SlugExists provides a mock function with given fields: ctx, slug
func (_m *ProfileRepository) SlugExists(ctx context.Context, slug string) (bool, error) {
	ret := _m.Called(ctx, slug)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (bool, error)); ok {
		return rf(ctx, slug)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) bool); ok {
		r0 = rf(ctx, slug)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, slug)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProfileRepository_SlugExists_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SlugExists'
type ProfileRepository_SlugExists_Call struct {
	*mock.Call
}

// SlugExists is a helper method to define mock.On call
//   - ctx context.Context
//   - slug string
func (_e *ProfileRepository_Expecter) SlugExists(ctx interface{}, slug interface{}) *ProfileRepository_SlugExists_Call {
	return &ProfileRepository_SlugExists_Call{Call: _e.mock.On("SlugExists", ctx, slug)}
}

func (_c *ProfileRepository_SlugExists_Call) Run(run func(ctx context.Context, slug string)) *ProfileRepository_SlugExists_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *ProfileRepository_SlugExists_Call) Return(_a0 bool, _a1 error) *ProfileRepository_SlugExists_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ProfileRepository_SlugExists_Call) RunAndReturn(run func(context.Context, string) (bool, error)) *ProfileRepository_SlugExists_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: ctx, data, id, now
func (_m *ProfileRepository) Update(ctx context.Context, data *dao.ProfileModelCore, id uuid.UUID, now time.Time) (*dao.ProfileModel, error) {
	ret := _m.Called(ctx, data, id, now)

	var r0 *dao.ProfileModel
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *dao.ProfileModelCore, uuid.UUID, time.Time) (*dao.ProfileModel, error)); ok {
		return rf(ctx, data, id, now)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *dao.ProfileModelCore, uuid.UUID, time.Time) *dao.ProfileModel); ok {
		r0 = rf(ctx, data, id, now)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dao.ProfileModel)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *dao.ProfileModelCore, uuid.UUID, time.Time) error); ok {
		r1 = rf(ctx, data, id, now)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProfileRepository_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type ProfileRepository_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx context.Context
//   - data *dao.ProfileModelCore
//   - id uuid.UUID
//   - now time.Time
func (_e *ProfileRepository_Expecter) Update(ctx interface{}, data interface{}, id interface{}, now interface{}) *ProfileRepository_Update_Call {
	return &ProfileRepository_Update_Call{Call: _e.mock.On("Update", ctx, data, id, now)}
}

func (_c *ProfileRepository_Update_Call) Run(run func(ctx context.Context, data *dao.ProfileModelCore, id uuid.UUID, now time.Time)) *ProfileRepository_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*dao.ProfileModelCore), args[2].(uuid.UUID), args[3].(time.Time))
	})
	return _c
}

func (_c *ProfileRepository_Update_Call) Return(_a0 *dao.ProfileModel, _a1 error) *ProfileRepository_Update_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ProfileRepository_Update_Call) RunAndReturn(run func(context.Context, *dao.ProfileModelCore, uuid.UUID, time.Time) (*dao.ProfileModel, error)) *ProfileRepository_Update_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewProfileRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewProfileRepository creates a new instance of ProfileRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewProfileRepository(t mockConstructorTestingTNewProfileRepository) *ProfileRepository {
	mock := &ProfileRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
