// Code generated by mockery v2.50.0. DO NOT EDIT.

package payment

import (
	context "context"

	grpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"

	payment "github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/payment"
)

// Payment is an autogenerated mock type for the Payment type
type Payment struct {
	mock.Mock
}

type Payment_Expecter struct {
	mock *mock.Mock
}

func (_m *Payment) EXPECT() *Payment_Expecter {
	return &Payment_Expecter{mock: &_m.Mock}
}

// CreatePayment provides a mock function with given fields: ctx, in, opts
func (_m *Payment) CreatePayment(ctx context.Context, in *payment.CreatePaymentRequest, opts ...grpc.CallOption) (*payment.CreatePaymentResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for CreatePayment")
	}

	var r0 *payment.CreatePaymentResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *payment.CreatePaymentRequest, ...grpc.CallOption) (*payment.CreatePaymentResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *payment.CreatePaymentRequest, ...grpc.CallOption) *payment.CreatePaymentResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*payment.CreatePaymentResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *payment.CreatePaymentRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Payment_CreatePayment_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreatePayment'
type Payment_CreatePayment_Call struct {
	*mock.Call
}

// CreatePayment is a helper method to define mock.On call
//   - ctx context.Context
//   - in *payment.CreatePaymentRequest
//   - opts ...grpc.CallOption
func (_e *Payment_Expecter) CreatePayment(ctx interface{}, in interface{}, opts ...interface{}) *Payment_CreatePayment_Call {
	return &Payment_CreatePayment_Call{Call: _e.mock.On("CreatePayment",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *Payment_CreatePayment_Call) Run(run func(ctx context.Context, in *payment.CreatePaymentRequest, opts ...grpc.CallOption)) *Payment_CreatePayment_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*payment.CreatePaymentRequest), variadicArgs...)
	})
	return _c
}

func (_c *Payment_CreatePayment_Call) Return(_a0 *payment.CreatePaymentResponse, _a1 error) *Payment_CreatePayment_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Payment_CreatePayment_Call) RunAndReturn(run func(context.Context, *payment.CreatePaymentRequest, ...grpc.CallOption) (*payment.CreatePaymentResponse, error)) *Payment_CreatePayment_Call {
	_c.Call.Return(run)
	return _c
}

// CreatePaymentChannel provides a mock function with given fields: ctx, in, opts
func (_m *Payment) CreatePaymentChannel(ctx context.Context, in *payment.CreatePaymentChannelRequest, opts ...grpc.CallOption) (*payment.CreatePaymentChannelResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for CreatePaymentChannel")
	}

	var r0 *payment.CreatePaymentChannelResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *payment.CreatePaymentChannelRequest, ...grpc.CallOption) (*payment.CreatePaymentChannelResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *payment.CreatePaymentChannelRequest, ...grpc.CallOption) *payment.CreatePaymentChannelResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*payment.CreatePaymentChannelResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *payment.CreatePaymentChannelRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Payment_CreatePaymentChannel_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreatePaymentChannel'
type Payment_CreatePaymentChannel_Call struct {
	*mock.Call
}

// CreatePaymentChannel is a helper method to define mock.On call
//   - ctx context.Context
//   - in *payment.CreatePaymentChannelRequest
//   - opts ...grpc.CallOption
func (_e *Payment_Expecter) CreatePaymentChannel(ctx interface{}, in interface{}, opts ...interface{}) *Payment_CreatePaymentChannel_Call {
	return &Payment_CreatePaymentChannel_Call{Call: _e.mock.On("CreatePaymentChannel",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *Payment_CreatePaymentChannel_Call) Run(run func(ctx context.Context, in *payment.CreatePaymentChannelRequest, opts ...grpc.CallOption)) *Payment_CreatePaymentChannel_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*payment.CreatePaymentChannelRequest), variadicArgs...)
	})
	return _c
}

func (_c *Payment_CreatePaymentChannel_Call) Return(_a0 *payment.CreatePaymentChannelResponse, _a1 error) *Payment_CreatePaymentChannel_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Payment_CreatePaymentChannel_Call) RunAndReturn(run func(context.Context, *payment.CreatePaymentChannelRequest, ...grpc.CallOption) (*payment.CreatePaymentChannelResponse, error)) *Payment_CreatePaymentChannel_Call {
	_c.Call.Return(run)
	return _c
}

// CreateRefund provides a mock function with given fields: ctx, in, opts
func (_m *Payment) CreateRefund(ctx context.Context, in *payment.CreateRefundRequest, opts ...grpc.CallOption) (*payment.CreateRefundResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for CreateRefund")
	}

	var r0 *payment.CreateRefundResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *payment.CreateRefundRequest, ...grpc.CallOption) (*payment.CreateRefundResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *payment.CreateRefundRequest, ...grpc.CallOption) *payment.CreateRefundResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*payment.CreateRefundResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *payment.CreateRefundRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Payment_CreateRefund_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateRefund'
type Payment_CreateRefund_Call struct {
	*mock.Call
}

// CreateRefund is a helper method to define mock.On call
//   - ctx context.Context
//   - in *payment.CreateRefundRequest
//   - opts ...grpc.CallOption
func (_e *Payment_Expecter) CreateRefund(ctx interface{}, in interface{}, opts ...interface{}) *Payment_CreateRefund_Call {
	return &Payment_CreateRefund_Call{Call: _e.mock.On("CreateRefund",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *Payment_CreateRefund_Call) Run(run func(ctx context.Context, in *payment.CreateRefundRequest, opts ...grpc.CallOption)) *Payment_CreateRefund_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*payment.CreateRefundRequest), variadicArgs...)
	})
	return _c
}

func (_c *Payment_CreateRefund_Call) Return(_a0 *payment.CreateRefundResponse, _a1 error) *Payment_CreateRefund_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Payment_CreateRefund_Call) RunAndReturn(run func(context.Context, *payment.CreateRefundRequest, ...grpc.CallOption) (*payment.CreateRefundResponse, error)) *Payment_CreateRefund_Call {
	_c.Call.Return(run)
	return _c
}

// GetPayment provides a mock function with given fields: ctx, in, opts
func (_m *Payment) GetPayment(ctx context.Context, in *payment.GetPaymentRequest, opts ...grpc.CallOption) (*payment.GetPaymentResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetPayment")
	}

	var r0 *payment.GetPaymentResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *payment.GetPaymentRequest, ...grpc.CallOption) (*payment.GetPaymentResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *payment.GetPaymentRequest, ...grpc.CallOption) *payment.GetPaymentResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*payment.GetPaymentResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *payment.GetPaymentRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Payment_GetPayment_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetPayment'
type Payment_GetPayment_Call struct {
	*mock.Call
}

// GetPayment is a helper method to define mock.On call
//   - ctx context.Context
//   - in *payment.GetPaymentRequest
//   - opts ...grpc.CallOption
func (_e *Payment_Expecter) GetPayment(ctx interface{}, in interface{}, opts ...interface{}) *Payment_GetPayment_Call {
	return &Payment_GetPayment_Call{Call: _e.mock.On("GetPayment",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *Payment_GetPayment_Call) Run(run func(ctx context.Context, in *payment.GetPaymentRequest, opts ...grpc.CallOption)) *Payment_GetPayment_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*payment.GetPaymentRequest), variadicArgs...)
	})
	return _c
}

func (_c *Payment_GetPayment_Call) Return(_a0 *payment.GetPaymentResponse, _a1 error) *Payment_GetPayment_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Payment_GetPayment_Call) RunAndReturn(run func(context.Context, *payment.GetPaymentRequest, ...grpc.CallOption) (*payment.GetPaymentResponse, error)) *Payment_GetPayment_Call {
	_c.Call.Return(run)
	return _c
}

// GetPaymentStatus provides a mock function with given fields: ctx, in, opts
func (_m *Payment) GetPaymentStatus(ctx context.Context, in *payment.GetPaymentStatusRequest, opts ...grpc.CallOption) (*payment.GetPaymentStatusResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetPaymentStatus")
	}

	var r0 *payment.GetPaymentStatusResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *payment.GetPaymentStatusRequest, ...grpc.CallOption) (*payment.GetPaymentStatusResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *payment.GetPaymentStatusRequest, ...grpc.CallOption) *payment.GetPaymentStatusResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*payment.GetPaymentStatusResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *payment.GetPaymentStatusRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Payment_GetPaymentStatus_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetPaymentStatus'
type Payment_GetPaymentStatus_Call struct {
	*mock.Call
}

// GetPaymentStatus is a helper method to define mock.On call
//   - ctx context.Context
//   - in *payment.GetPaymentStatusRequest
//   - opts ...grpc.CallOption
func (_e *Payment_Expecter) GetPaymentStatus(ctx interface{}, in interface{}, opts ...interface{}) *Payment_GetPaymentStatus_Call {
	return &Payment_GetPaymentStatus_Call{Call: _e.mock.On("GetPaymentStatus",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *Payment_GetPaymentStatus_Call) Run(run func(ctx context.Context, in *payment.GetPaymentStatusRequest, opts ...grpc.CallOption)) *Payment_GetPaymentStatus_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*payment.GetPaymentStatusRequest), variadicArgs...)
	})
	return _c
}

func (_c *Payment_GetPaymentStatus_Call) Return(_a0 *payment.GetPaymentStatusResponse, _a1 error) *Payment_GetPaymentStatus_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Payment_GetPaymentStatus_Call) RunAndReturn(run func(context.Context, *payment.GetPaymentStatusRequest, ...grpc.CallOption) (*payment.GetPaymentStatusResponse, error)) *Payment_GetPaymentStatus_Call {
	_c.Call.Return(run)
	return _c
}

// GetRefund provides a mock function with given fields: ctx, in, opts
func (_m *Payment) GetRefund(ctx context.Context, in *payment.GetRefundRequest, opts ...grpc.CallOption) (*payment.GetRefundResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetRefund")
	}

	var r0 *payment.GetRefundResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *payment.GetRefundRequest, ...grpc.CallOption) (*payment.GetRefundResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *payment.GetRefundRequest, ...grpc.CallOption) *payment.GetRefundResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*payment.GetRefundResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *payment.GetRefundRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Payment_GetRefund_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRefund'
type Payment_GetRefund_Call struct {
	*mock.Call
}

// GetRefund is a helper method to define mock.On call
//   - ctx context.Context
//   - in *payment.GetRefundRequest
//   - opts ...grpc.CallOption
func (_e *Payment_Expecter) GetRefund(ctx interface{}, in interface{}, opts ...interface{}) *Payment_GetRefund_Call {
	return &Payment_GetRefund_Call{Call: _e.mock.On("GetRefund",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *Payment_GetRefund_Call) Run(run func(ctx context.Context, in *payment.GetRefundRequest, opts ...grpc.CallOption)) *Payment_GetRefund_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*payment.GetRefundRequest), variadicArgs...)
	})
	return _c
}

func (_c *Payment_GetRefund_Call) Return(_a0 *payment.GetRefundResponse, _a1 error) *Payment_GetRefund_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Payment_GetRefund_Call) RunAndReturn(run func(context.Context, *payment.GetRefundRequest, ...grpc.CallOption) (*payment.GetRefundResponse, error)) *Payment_GetRefund_Call {
	_c.Call.Return(run)
	return _c
}

// ListPaymentChannels provides a mock function with given fields: ctx, in, opts
func (_m *Payment) ListPaymentChannels(ctx context.Context, in *payment.ListPaymentChannelsRequest, opts ...grpc.CallOption) (*payment.ListPaymentChannelsResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for ListPaymentChannels")
	}

	var r0 *payment.ListPaymentChannelsResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *payment.ListPaymentChannelsRequest, ...grpc.CallOption) (*payment.ListPaymentChannelsResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *payment.ListPaymentChannelsRequest, ...grpc.CallOption) *payment.ListPaymentChannelsResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*payment.ListPaymentChannelsResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *payment.ListPaymentChannelsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Payment_ListPaymentChannels_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListPaymentChannels'
type Payment_ListPaymentChannels_Call struct {
	*mock.Call
}

// ListPaymentChannels is a helper method to define mock.On call
//   - ctx context.Context
//   - in *payment.ListPaymentChannelsRequest
//   - opts ...grpc.CallOption
func (_e *Payment_Expecter) ListPaymentChannels(ctx interface{}, in interface{}, opts ...interface{}) *Payment_ListPaymentChannels_Call {
	return &Payment_ListPaymentChannels_Call{Call: _e.mock.On("ListPaymentChannels",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *Payment_ListPaymentChannels_Call) Run(run func(ctx context.Context, in *payment.ListPaymentChannelsRequest, opts ...grpc.CallOption)) *Payment_ListPaymentChannels_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*payment.ListPaymentChannelsRequest), variadicArgs...)
	})
	return _c
}

func (_c *Payment_ListPaymentChannels_Call) Return(_a0 *payment.ListPaymentChannelsResponse, _a1 error) *Payment_ListPaymentChannels_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Payment_ListPaymentChannels_Call) RunAndReturn(run func(context.Context, *payment.ListPaymentChannelsRequest, ...grpc.CallOption) (*payment.ListPaymentChannelsResponse, error)) *Payment_ListPaymentChannels_Call {
	_c.Call.Return(run)
	return _c
}

// PaymentNotify provides a mock function with given fields: ctx, in, opts
func (_m *Payment) PaymentNotify(ctx context.Context, in *payment.PaymentNotifyRequest, opts ...grpc.CallOption) (*payment.PaymentNotifyResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for PaymentNotify")
	}

	var r0 *payment.PaymentNotifyResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *payment.PaymentNotifyRequest, ...grpc.CallOption) (*payment.PaymentNotifyResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *payment.PaymentNotifyRequest, ...grpc.CallOption) *payment.PaymentNotifyResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*payment.PaymentNotifyResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *payment.PaymentNotifyRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Payment_PaymentNotify_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PaymentNotify'
type Payment_PaymentNotify_Call struct {
	*mock.Call
}

// PaymentNotify is a helper method to define mock.On call
//   - ctx context.Context
//   - in *payment.PaymentNotifyRequest
//   - opts ...grpc.CallOption
func (_e *Payment_Expecter) PaymentNotify(ctx interface{}, in interface{}, opts ...interface{}) *Payment_PaymentNotify_Call {
	return &Payment_PaymentNotify_Call{Call: _e.mock.On("PaymentNotify",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *Payment_PaymentNotify_Call) Run(run func(ctx context.Context, in *payment.PaymentNotifyRequest, opts ...grpc.CallOption)) *Payment_PaymentNotify_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*payment.PaymentNotifyRequest), variadicArgs...)
	})
	return _c
}

func (_c *Payment_PaymentNotify_Call) Return(_a0 *payment.PaymentNotifyResponse, _a1 error) *Payment_PaymentNotify_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Payment_PaymentNotify_Call) RunAndReturn(run func(context.Context, *payment.PaymentNotifyRequest, ...grpc.CallOption) (*payment.PaymentNotifyResponse, error)) *Payment_PaymentNotify_Call {
	_c.Call.Return(run)
	return _c
}

// RefundNotify provides a mock function with given fields: ctx, in, opts
func (_m *Payment) RefundNotify(ctx context.Context, in *payment.RefundNotifyRequest, opts ...grpc.CallOption) (*payment.RefundNotifyResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for RefundNotify")
	}

	var r0 *payment.RefundNotifyResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *payment.RefundNotifyRequest, ...grpc.CallOption) (*payment.RefundNotifyResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *payment.RefundNotifyRequest, ...grpc.CallOption) *payment.RefundNotifyResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*payment.RefundNotifyResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *payment.RefundNotifyRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Payment_RefundNotify_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RefundNotify'
type Payment_RefundNotify_Call struct {
	*mock.Call
}

// RefundNotify is a helper method to define mock.On call
//   - ctx context.Context
//   - in *payment.RefundNotifyRequest
//   - opts ...grpc.CallOption
func (_e *Payment_Expecter) RefundNotify(ctx interface{}, in interface{}, opts ...interface{}) *Payment_RefundNotify_Call {
	return &Payment_RefundNotify_Call{Call: _e.mock.On("RefundNotify",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *Payment_RefundNotify_Call) Run(run func(ctx context.Context, in *payment.RefundNotifyRequest, opts ...grpc.CallOption)) *Payment_RefundNotify_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*payment.RefundNotifyRequest), variadicArgs...)
	})
	return _c
}

func (_c *Payment_RefundNotify_Call) Return(_a0 *payment.RefundNotifyResponse, _a1 error) *Payment_RefundNotify_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Payment_RefundNotify_Call) RunAndReturn(run func(context.Context, *payment.RefundNotifyRequest, ...grpc.CallOption) (*payment.RefundNotifyResponse, error)) *Payment_RefundNotify_Call {
	_c.Call.Return(run)
	return _c
}

// UpdatePaymentChannel provides a mock function with given fields: ctx, in, opts
func (_m *Payment) UpdatePaymentChannel(ctx context.Context, in *payment.UpdatePaymentChannelRequest, opts ...grpc.CallOption) (*payment.UpdatePaymentChannelResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for UpdatePaymentChannel")
	}

	var r0 *payment.UpdatePaymentChannelResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *payment.UpdatePaymentChannelRequest, ...grpc.CallOption) (*payment.UpdatePaymentChannelResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *payment.UpdatePaymentChannelRequest, ...grpc.CallOption) *payment.UpdatePaymentChannelResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*payment.UpdatePaymentChannelResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *payment.UpdatePaymentChannelRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Payment_UpdatePaymentChannel_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdatePaymentChannel'
type Payment_UpdatePaymentChannel_Call struct {
	*mock.Call
}

// UpdatePaymentChannel is a helper method to define mock.On call
//   - ctx context.Context
//   - in *payment.UpdatePaymentChannelRequest
//   - opts ...grpc.CallOption
func (_e *Payment_Expecter) UpdatePaymentChannel(ctx interface{}, in interface{}, opts ...interface{}) *Payment_UpdatePaymentChannel_Call {
	return &Payment_UpdatePaymentChannel_Call{Call: _e.mock.On("UpdatePaymentChannel",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *Payment_UpdatePaymentChannel_Call) Run(run func(ctx context.Context, in *payment.UpdatePaymentChannelRequest, opts ...grpc.CallOption)) *Payment_UpdatePaymentChannel_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*payment.UpdatePaymentChannelRequest), variadicArgs...)
	})
	return _c
}

func (_c *Payment_UpdatePaymentChannel_Call) Return(_a0 *payment.UpdatePaymentChannelResponse, _a1 error) *Payment_UpdatePaymentChannel_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Payment_UpdatePaymentChannel_Call) RunAndReturn(run func(context.Context, *payment.UpdatePaymentChannelRequest, ...grpc.CallOption) (*payment.UpdatePaymentChannelResponse, error)) *Payment_UpdatePaymentChannel_Call {
	_c.Call.Return(run)
	return _c
}

// NewPayment creates a new instance of Payment. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPayment(t interface {
	mock.TestingT
	Cleanup(func())
}) *Payment {
	mock := &Payment{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
