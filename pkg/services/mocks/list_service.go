// Code generated by mockery v2.20.0. DO NOT EDIT.

package servicesmocks

import (
	context "context"

	models "github.com/a-novel/auth-service/pkg/models"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// ListService is an autogenerated mock type for the ListService type
type ListService struct {
	mock.Mock
}

type ListService_Expecter struct {
	mock *mock.Mock
}

func (_m *ListService) EXPECT() *ListService_Expecter {
	return &ListService_Expecter{mock: &_m.Mock}
}

// List provides a mock function with given fields: ctx, ids
func (_m *ListService) List(ctx context.Context, ids []uuid.UUID) ([]*models.UserPreview, error) {
	ret := _m.Called(ctx, ids)

	var r0 []*models.UserPreview
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []uuid.UUID) ([]*models.UserPreview, error)); ok {
		return rf(ctx, ids)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []uuid.UUID) []*models.UserPreview); ok {
		r0 = rf(ctx, ids)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.UserPreview)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []uuid.UUID) error); ok {
		r1 = rf(ctx, ids)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListService_List_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'List'
type ListService_List_Call struct {
	*mock.Call
}

// List is a helper method to define mock.On call
//   - ctx context.Context
//   - ids []uuid.UUID
func (_e *ListService_Expecter) List(ctx interface{}, ids interface{}) *ListService_List_Call {
	return &ListService_List_Call{Call: _e.mock.On("List", ctx, ids)}
}

func (_c *ListService_List_Call) Run(run func(ctx context.Context, ids []uuid.UUID)) *ListService_List_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]uuid.UUID))
	})
	return _c
}

func (_c *ListService_List_Call) Return(_a0 []*models.UserPreview, _a1 error) *ListService_List_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ListService_List_Call) RunAndReturn(run func(context.Context, []uuid.UUID) ([]*models.UserPreview, error)) *ListService_List_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewListService interface {
	mock.TestingT
	Cleanup(func())
}

// NewListService creates a new instance of ListService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewListService(t mockConstructorTestingTNewListService) *ListService {
	mock := &ListService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
