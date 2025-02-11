// Code generated by MockGen. DO NOT EDIT.
// Source: ./order.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	dto "hangry/domain/dto"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockOrderUsecase is a mock of OrderUsecase interface.
type MockOrderUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockOrderUsecaseMockRecorder
}

// MockOrderUsecaseMockRecorder is the mock recorder for MockOrderUsecase.
type MockOrderUsecaseMockRecorder struct {
	mock *MockOrderUsecase
}

// NewMockOrderUsecase creates a new mock instance.
func NewMockOrderUsecase(ctrl *gomock.Controller) *MockOrderUsecase {
	mock := &MockOrderUsecase{ctrl: ctrl}
	mock.recorder = &MockOrderUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrderUsecase) EXPECT() *MockOrderUsecaseMockRecorder {
	return m.recorder
}

// CreateOrder mocks base method.
func (m *MockOrderUsecase) CreateOrder(ctx context.Context, dto dto.OrderInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrder", ctx, dto)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateOrder indicates an expected call of CreateOrder.
func (mr *MockOrderUsecaseMockRecorder) CreateOrder(ctx, dto interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrder", reflect.TypeOf((*MockOrderUsecase)(nil).CreateOrder), ctx, dto)
}
