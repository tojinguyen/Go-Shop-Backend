// Code generated by mockery v2.53.4. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/toji-dev/go-shop/internal/services/user-service/internal/domain"
	dto "github.com/toji-dev/go-shop/internal/services/user-service/internal/dto"

	mock "github.com/stretchr/testify/mock"
)

// ShipperRepository is an autogenerated mock type for the ShipperRepository type
type ShipperRepository struct {
	mock.Mock
}

type ShipperRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *ShipperRepository) EXPECT() *ShipperRepository_Expecter {
	return &ShipperRepository_Expecter{mock: &_m.Mock}
}

// CreateShipper provides a mock function with given fields: ctx, shipper
func (_m *ShipperRepository) CreateShipper(ctx context.Context, shipper *domain.Shipper) (*domain.Shipper, error) {
	ret := _m.Called(ctx, shipper)

	if len(ret) == 0 {
		panic("no return value specified for CreateShipper")
	}

	var r0 *domain.Shipper
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Shipper) (*domain.Shipper, error)); ok {
		return rf(ctx, shipper)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Shipper) *domain.Shipper); ok {
		r0 = rf(ctx, shipper)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Shipper)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *domain.Shipper) error); ok {
		r1 = rf(ctx, shipper)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ShipperRepository_CreateShipper_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateShipper'
type ShipperRepository_CreateShipper_Call struct {
	*mock.Call
}

// CreateShipper is a helper method to define mock.On call
//   - ctx context.Context
//   - shipper *domain.Shipper
func (_e *ShipperRepository_Expecter) CreateShipper(ctx interface{}, shipper interface{}) *ShipperRepository_CreateShipper_Call {
	return &ShipperRepository_CreateShipper_Call{Call: _e.mock.On("CreateShipper", ctx, shipper)}
}

func (_c *ShipperRepository_CreateShipper_Call) Run(run func(ctx context.Context, shipper *domain.Shipper)) *ShipperRepository_CreateShipper_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*domain.Shipper))
	})
	return _c
}

func (_c *ShipperRepository_CreateShipper_Call) Return(_a0 *domain.Shipper, _a1 error) *ShipperRepository_CreateShipper_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ShipperRepository_CreateShipper_Call) RunAndReturn(run func(context.Context, *domain.Shipper) (*domain.Shipper, error)) *ShipperRepository_CreateShipper_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteShipper provides a mock function with given fields: ctx, userID
func (_m *ShipperRepository) DeleteShipper(ctx context.Context, userID string) error {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteShipper")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ShipperRepository_DeleteShipper_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteShipper'
type ShipperRepository_DeleteShipper_Call struct {
	*mock.Call
}

// DeleteShipper is a helper method to define mock.On call
//   - ctx context.Context
//   - userID string
func (_e *ShipperRepository_Expecter) DeleteShipper(ctx interface{}, userID interface{}) *ShipperRepository_DeleteShipper_Call {
	return &ShipperRepository_DeleteShipper_Call{Call: _e.mock.On("DeleteShipper", ctx, userID)}
}

func (_c *ShipperRepository_DeleteShipper_Call) Run(run func(ctx context.Context, userID string)) *ShipperRepository_DeleteShipper_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *ShipperRepository_DeleteShipper_Call) Return(_a0 error) *ShipperRepository_DeleteShipper_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ShipperRepository_DeleteShipper_Call) RunAndReturn(run func(context.Context, string) error) *ShipperRepository_DeleteShipper_Call {
	_c.Call.Return(run)
	return _c
}

// GetShipperByUserID provides a mock function with given fields: ctx, userID
func (_m *ShipperRepository) GetShipperByUserID(ctx context.Context, userID string) (*domain.Shipper, error) {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for GetShipperByUserID")
	}

	var r0 *domain.Shipper
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*domain.Shipper, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *domain.Shipper); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Shipper)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ShipperRepository_GetShipperByUserID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetShipperByUserID'
type ShipperRepository_GetShipperByUserID_Call struct {
	*mock.Call
}

// GetShipperByUserID is a helper method to define mock.On call
//   - ctx context.Context
//   - userID string
func (_e *ShipperRepository_Expecter) GetShipperByUserID(ctx interface{}, userID interface{}) *ShipperRepository_GetShipperByUserID_Call {
	return &ShipperRepository_GetShipperByUserID_Call{Call: _e.mock.On("GetShipperByUserID", ctx, userID)}
}

func (_c *ShipperRepository_GetShipperByUserID_Call) Run(run func(ctx context.Context, userID string)) *ShipperRepository_GetShipperByUserID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *ShipperRepository_GetShipperByUserID_Call) Return(_a0 *domain.Shipper, _a1 error) *ShipperRepository_GetShipperByUserID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ShipperRepository_GetShipperByUserID_Call) RunAndReturn(run func(context.Context, string) (*domain.Shipper, error)) *ShipperRepository_GetShipperByUserID_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateShipper provides a mock function with given fields: ctx, userID, updateRequest
func (_m *ShipperRepository) UpdateShipper(ctx context.Context, userID string, updateRequest *dto.ShipperUpdateRequest) (*domain.Shipper, error) {
	ret := _m.Called(ctx, userID, updateRequest)

	if len(ret) == 0 {
		panic("no return value specified for UpdateShipper")
	}

	var r0 *domain.Shipper
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *dto.ShipperUpdateRequest) (*domain.Shipper, error)); ok {
		return rf(ctx, userID, updateRequest)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, *dto.ShipperUpdateRequest) *domain.Shipper); ok {
		r0 = rf(ctx, userID, updateRequest)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Shipper)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, *dto.ShipperUpdateRequest) error); ok {
		r1 = rf(ctx, userID, updateRequest)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ShipperRepository_UpdateShipper_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateShipper'
type ShipperRepository_UpdateShipper_Call struct {
	*mock.Call
}

// UpdateShipper is a helper method to define mock.On call
//   - ctx context.Context
//   - userID string
//   - updateRequest *dto.ShipperUpdateRequest
func (_e *ShipperRepository_Expecter) UpdateShipper(ctx interface{}, userID interface{}, updateRequest interface{}) *ShipperRepository_UpdateShipper_Call {
	return &ShipperRepository_UpdateShipper_Call{Call: _e.mock.On("UpdateShipper", ctx, userID, updateRequest)}
}

func (_c *ShipperRepository_UpdateShipper_Call) Run(run func(ctx context.Context, userID string, updateRequest *dto.ShipperUpdateRequest)) *ShipperRepository_UpdateShipper_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(*dto.ShipperUpdateRequest))
	})
	return _c
}

func (_c *ShipperRepository_UpdateShipper_Call) Return(_a0 *domain.Shipper, _a1 error) *ShipperRepository_UpdateShipper_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ShipperRepository_UpdateShipper_Call) RunAndReturn(run func(context.Context, string, *dto.ShipperUpdateRequest) (*domain.Shipper, error)) *ShipperRepository_UpdateShipper_Call {
	_c.Call.Return(run)
	return _c
}

// NewShipperRepository creates a new instance of ShipperRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewShipperRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *ShipperRepository {
	mock := &ShipperRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
