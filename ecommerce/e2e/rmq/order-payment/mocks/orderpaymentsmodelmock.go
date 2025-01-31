// Code generated by mockery v2.51.1. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/model"
	mock "github.com/stretchr/testify/mock"

	sql "database/sql"

	time "time"
)

// OrderPaymentsModel is an autogenerated mock type for the OrderPaymentsModel type
type OrderPaymentsModel struct {
	mock.Mock
}

type OrderPaymentsModel_Expecter struct {
	mock *mock.Mock
}

func (_m *OrderPaymentsModel) EXPECT() *OrderPaymentsModel_Expecter {
	return &OrderPaymentsModel_Expecter{mock: &_m.Mock}
}

// CreatePayment provides a mock function with given fields: ctx, payment
func (_m *OrderPaymentsModel) CreatePayment(ctx context.Context, payment *model.OrderPayments) (uint64, error) {
	ret := _m.Called(ctx, payment)

	if len(ret) == 0 {
		panic("no return value specified for CreatePayment")
	}

	var r0 uint64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.OrderPayments) (uint64, error)); ok {
		return rf(ctx, payment)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.OrderPayments) uint64); ok {
		r0 = rf(ctx, payment)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.OrderPayments) error); ok {
		r1 = rf(ctx, payment)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OrderPaymentsModel_CreatePayment_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreatePayment'
type OrderPaymentsModel_CreatePayment_Call struct {
	*mock.Call
}

// CreatePayment is a helper method to define mock.On call
//   - ctx context.Context
//   - payment *model.OrderPayments
func (_e *OrderPaymentsModel_Expecter) CreatePayment(ctx interface{}, payment interface{}) *OrderPaymentsModel_CreatePayment_Call {
	return &OrderPaymentsModel_CreatePayment_Call{Call: _e.mock.On("CreatePayment", ctx, payment)}
}

func (_c *OrderPaymentsModel_CreatePayment_Call) Run(run func(ctx context.Context, payment *model.OrderPayments)) *OrderPaymentsModel_CreatePayment_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.OrderPayments))
	})
	return _c
}

func (_c *OrderPaymentsModel_CreatePayment_Call) Return(_a0 uint64, _a1 error) *OrderPaymentsModel_CreatePayment_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *OrderPaymentsModel_CreatePayment_Call) RunAndReturn(run func(context.Context, *model.OrderPayments) (uint64, error)) *OrderPaymentsModel_CreatePayment_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function with given fields: ctx, id
func (_m *OrderPaymentsModel) Delete(ctx context.Context, id uint64) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// OrderPaymentsModel_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type OrderPaymentsModel_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx context.Context
//   - id uint64
func (_e *OrderPaymentsModel_Expecter) Delete(ctx interface{}, id interface{}) *OrderPaymentsModel_Delete_Call {
	return &OrderPaymentsModel_Delete_Call{Call: _e.mock.On("Delete", ctx, id)}
}

func (_c *OrderPaymentsModel_Delete_Call) Run(run func(ctx context.Context, id uint64)) *OrderPaymentsModel_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint64))
	})
	return _c
}

func (_c *OrderPaymentsModel_Delete_Call) Return(_a0 error) *OrderPaymentsModel_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *OrderPaymentsModel_Delete_Call) RunAndReturn(run func(context.Context, uint64) error) *OrderPaymentsModel_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// FindByOrderId provides a mock function with given fields: ctx, orderId
func (_m *OrderPaymentsModel) FindByOrderId(ctx context.Context, orderId uint64) (*model.OrderPayments, error) {
	ret := _m.Called(ctx, orderId)

	if len(ret) == 0 {
		panic("no return value specified for FindByOrderId")
	}

	var r0 *model.OrderPayments
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64) (*model.OrderPayments, error)); ok {
		return rf(ctx, orderId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint64) *model.OrderPayments); ok {
		r0 = rf(ctx, orderId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.OrderPayments)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint64) error); ok {
		r1 = rf(ctx, orderId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OrderPaymentsModel_FindByOrderId_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindByOrderId'
type OrderPaymentsModel_FindByOrderId_Call struct {
	*mock.Call
}

// FindByOrderId is a helper method to define mock.On call
//   - ctx context.Context
//   - orderId uint64
func (_e *OrderPaymentsModel_Expecter) FindByOrderId(ctx interface{}, orderId interface{}) *OrderPaymentsModel_FindByOrderId_Call {
	return &OrderPaymentsModel_FindByOrderId_Call{Call: _e.mock.On("FindByOrderId", ctx, orderId)}
}

func (_c *OrderPaymentsModel_FindByOrderId_Call) Run(run func(ctx context.Context, orderId uint64)) *OrderPaymentsModel_FindByOrderId_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint64))
	})
	return _c
}

func (_c *OrderPaymentsModel_FindByOrderId_Call) Return(_a0 *model.OrderPayments, _a1 error) *OrderPaymentsModel_FindByOrderId_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *OrderPaymentsModel_FindByOrderId_Call) RunAndReturn(run func(context.Context, uint64) (*model.OrderPayments, error)) *OrderPaymentsModel_FindByOrderId_Call {
	_c.Call.Return(run)
	return _c
}

// FindByStatusAndTime provides a mock function with given fields: ctx, status, startTime, endTime
func (_m *OrderPaymentsModel) FindByStatusAndTime(ctx context.Context, status int64, startTime time.Time, endTime time.Time) ([]*model.OrderPayments, error) {
	ret := _m.Called(ctx, status, startTime, endTime)

	if len(ret) == 0 {
		panic("no return value specified for FindByStatusAndTime")
	}

	var r0 []*model.OrderPayments
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, time.Time, time.Time) ([]*model.OrderPayments, error)); ok {
		return rf(ctx, status, startTime, endTime)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64, time.Time, time.Time) []*model.OrderPayments); ok {
		r0 = rf(ctx, status, startTime, endTime)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.OrderPayments)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64, time.Time, time.Time) error); ok {
		r1 = rf(ctx, status, startTime, endTime)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OrderPaymentsModel_FindByStatusAndTime_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindByStatusAndTime'
type OrderPaymentsModel_FindByStatusAndTime_Call struct {
	*mock.Call
}

// FindByStatusAndTime is a helper method to define mock.On call
//   - ctx context.Context
//   - status int64
//   - startTime time.Time
//   - endTime time.Time
func (_e *OrderPaymentsModel_Expecter) FindByStatusAndTime(ctx interface{}, status interface{}, startTime interface{}, endTime interface{}) *OrderPaymentsModel_FindByStatusAndTime_Call {
	return &OrderPaymentsModel_FindByStatusAndTime_Call{Call: _e.mock.On("FindByStatusAndTime", ctx, status, startTime, endTime)}
}

func (_c *OrderPaymentsModel_FindByStatusAndTime_Call) Run(run func(ctx context.Context, status int64, startTime time.Time, endTime time.Time)) *OrderPaymentsModel_FindByStatusAndTime_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64), args[2].(time.Time), args[3].(time.Time))
	})
	return _c
}

func (_c *OrderPaymentsModel_FindByStatusAndTime_Call) Return(_a0 []*model.OrderPayments, _a1 error) *OrderPaymentsModel_FindByStatusAndTime_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *OrderPaymentsModel_FindByStatusAndTime_Call) RunAndReturn(run func(context.Context, int64, time.Time, time.Time) ([]*model.OrderPayments, error)) *OrderPaymentsModel_FindByStatusAndTime_Call {
	_c.Call.Return(run)
	return _c
}

// FindOne provides a mock function with given fields: ctx, id
func (_m *OrderPaymentsModel) FindOne(ctx context.Context, id uint64) (*model.OrderPayments, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for FindOne")
	}

	var r0 *model.OrderPayments
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64) (*model.OrderPayments, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint64) *model.OrderPayments); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.OrderPayments)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OrderPaymentsModel_FindOne_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindOne'
type OrderPaymentsModel_FindOne_Call struct {
	*mock.Call
}

// FindOne is a helper method to define mock.On call
//   - ctx context.Context
//   - id uint64
func (_e *OrderPaymentsModel_Expecter) FindOne(ctx interface{}, id interface{}) *OrderPaymentsModel_FindOne_Call {
	return &OrderPaymentsModel_FindOne_Call{Call: _e.mock.On("FindOne", ctx, id)}
}

func (_c *OrderPaymentsModel_FindOne_Call) Run(run func(ctx context.Context, id uint64)) *OrderPaymentsModel_FindOne_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint64))
	})
	return _c
}

func (_c *OrderPaymentsModel_FindOne_Call) Return(_a0 *model.OrderPayments, _a1 error) *OrderPaymentsModel_FindOne_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *OrderPaymentsModel_FindOne_Call) RunAndReturn(run func(context.Context, uint64) (*model.OrderPayments, error)) *OrderPaymentsModel_FindOne_Call {
	_c.Call.Return(run)
	return _c
}

// FindOneByPaymentNo provides a mock function with given fields: ctx, paymentNo
func (_m *OrderPaymentsModel) FindOneByPaymentNo(ctx context.Context, paymentNo string) (*model.OrderPayments, error) {
	ret := _m.Called(ctx, paymentNo)

	if len(ret) == 0 {
		panic("no return value specified for FindOneByPaymentNo")
	}

	var r0 *model.OrderPayments
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*model.OrderPayments, error)); ok {
		return rf(ctx, paymentNo)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.OrderPayments); ok {
		r0 = rf(ctx, paymentNo)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.OrderPayments)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, paymentNo)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OrderPaymentsModel_FindOneByPaymentNo_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindOneByPaymentNo'
type OrderPaymentsModel_FindOneByPaymentNo_Call struct {
	*mock.Call
}

// FindOneByPaymentNo is a helper method to define mock.On call
//   - ctx context.Context
//   - paymentNo string
func (_e *OrderPaymentsModel_Expecter) FindOneByPaymentNo(ctx interface{}, paymentNo interface{}) *OrderPaymentsModel_FindOneByPaymentNo_Call {
	return &OrderPaymentsModel_FindOneByPaymentNo_Call{Call: _e.mock.On("FindOneByPaymentNo", ctx, paymentNo)}
}

func (_c *OrderPaymentsModel_FindOneByPaymentNo_Call) Run(run func(ctx context.Context, paymentNo string)) *OrderPaymentsModel_FindOneByPaymentNo_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *OrderPaymentsModel_FindOneByPaymentNo_Call) Return(_a0 *model.OrderPayments, _a1 error) *OrderPaymentsModel_FindOneByPaymentNo_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *OrderPaymentsModel_FindOneByPaymentNo_Call) RunAndReturn(run func(context.Context, string) (*model.OrderPayments, error)) *OrderPaymentsModel_FindOneByPaymentNo_Call {
	_c.Call.Return(run)
	return _c
}

// Insert provides a mock function with given fields: ctx, data
func (_m *OrderPaymentsModel) Insert(ctx context.Context, data *model.OrderPayments) (sql.Result, error) {
	ret := _m.Called(ctx, data)

	if len(ret) == 0 {
		panic("no return value specified for Insert")
	}

	var r0 sql.Result
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.OrderPayments) (sql.Result, error)); ok {
		return rf(ctx, data)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.OrderPayments) sql.Result); ok {
		r0 = rf(ctx, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(sql.Result)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.OrderPayments) error); ok {
		r1 = rf(ctx, data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OrderPaymentsModel_Insert_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Insert'
type OrderPaymentsModel_Insert_Call struct {
	*mock.Call
}

// Insert is a helper method to define mock.On call
//   - ctx context.Context
//   - data *model.OrderPayments
func (_e *OrderPaymentsModel_Expecter) Insert(ctx interface{}, data interface{}) *OrderPaymentsModel_Insert_Call {
	return &OrderPaymentsModel_Insert_Call{Call: _e.mock.On("Insert", ctx, data)}
}

func (_c *OrderPaymentsModel_Insert_Call) Run(run func(ctx context.Context, data *model.OrderPayments)) *OrderPaymentsModel_Insert_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.OrderPayments))
	})
	return _c
}

func (_c *OrderPaymentsModel_Insert_Call) Return(_a0 sql.Result, _a1 error) *OrderPaymentsModel_Insert_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *OrderPaymentsModel_Insert_Call) RunAndReturn(run func(context.Context, *model.OrderPayments) (sql.Result, error)) *OrderPaymentsModel_Insert_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: ctx, data
func (_m *OrderPaymentsModel) Update(ctx context.Context, data *model.OrderPayments) error {
	ret := _m.Called(ctx, data)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.OrderPayments) error); ok {
		r0 = rf(ctx, data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// OrderPaymentsModel_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type OrderPaymentsModel_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx context.Context
//   - data *model.OrderPayments
func (_e *OrderPaymentsModel_Expecter) Update(ctx interface{}, data interface{}) *OrderPaymentsModel_Update_Call {
	return &OrderPaymentsModel_Update_Call{Call: _e.mock.On("Update", ctx, data)}
}

func (_c *OrderPaymentsModel_Update_Call) Run(run func(ctx context.Context, data *model.OrderPayments)) *OrderPaymentsModel_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.OrderPayments))
	})
	return _c
}

func (_c *OrderPaymentsModel_Update_Call) Return(_a0 error) *OrderPaymentsModel_Update_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *OrderPaymentsModel_Update_Call) RunAndReturn(run func(context.Context, *model.OrderPayments) error) *OrderPaymentsModel_Update_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateStatus provides a mock function with given fields: ctx, paymentNo, status, payTime
func (_m *OrderPaymentsModel) UpdateStatus(ctx context.Context, paymentNo string, status int64, payTime time.Time) error {
	ret := _m.Called(ctx, paymentNo, status, payTime)

	if len(ret) == 0 {
		panic("no return value specified for UpdateStatus")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, int64, time.Time) error); ok {
		r0 = rf(ctx, paymentNo, status, payTime)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// OrderPaymentsModel_UpdateStatus_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateStatus'
type OrderPaymentsModel_UpdateStatus_Call struct {
	*mock.Call
}

// UpdateStatus is a helper method to define mock.On call
//   - ctx context.Context
//   - paymentNo string
//   - status int64
//   - payTime time.Time
func (_e *OrderPaymentsModel_Expecter) UpdateStatus(ctx interface{}, paymentNo interface{}, status interface{}, payTime interface{}) *OrderPaymentsModel_UpdateStatus_Call {
	return &OrderPaymentsModel_UpdateStatus_Call{Call: _e.mock.On("UpdateStatus", ctx, paymentNo, status, payTime)}
}

func (_c *OrderPaymentsModel_UpdateStatus_Call) Run(run func(ctx context.Context, paymentNo string, status int64, payTime time.Time)) *OrderPaymentsModel_UpdateStatus_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(int64), args[3].(time.Time))
	})
	return _c
}

func (_c *OrderPaymentsModel_UpdateStatus_Call) Return(_a0 error) *OrderPaymentsModel_UpdateStatus_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *OrderPaymentsModel_UpdateStatus_Call) RunAndReturn(run func(context.Context, string, int64, time.Time) error) *OrderPaymentsModel_UpdateStatus_Call {
	_c.Call.Return(run)
	return _c
}

// NewOrderPaymentsModel creates a new instance of OrderPaymentsModel. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewOrderPaymentsModel(t interface {
	mock.TestingT
	Cleanup(func())
}) *OrderPaymentsModel {
	mock := &OrderPaymentsModel{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
