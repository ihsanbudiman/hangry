// Code generated by MockGen. DO NOT EDIT.
// Source: ./promo_repository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	models "hangry/domain/models"
	repository "hangry/repository"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	gorm "gorm.io/gorm"
)

// MockPromoRepository is a mock of PromoRepository interface.
type MockPromoRepository struct {
	ctrl     *gomock.Controller
	recorder *MockPromoRepositoryMockRecorder
}

// MockPromoRepositoryMockRecorder is the mock recorder for MockPromoRepository.
type MockPromoRepositoryMockRecorder struct {
	mock *MockPromoRepository
}

// NewMockPromoRepository creates a new mock instance.
func NewMockPromoRepository(ctrl *gomock.Controller) *MockPromoRepository {
	mock := &MockPromoRepository{ctrl: ctrl}
	mock.recorder = &MockPromoRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPromoRepository) EXPECT() *MockPromoRepositoryMockRecorder {
	return m.recorder
}

// GetPromoByPromoID mocks base method.
func (m *MockPromoRepository) GetPromoByPromoID(ctx context.Context, tx *gorm.DB, promoID uint) (models.Promo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPromoByPromoID", ctx, tx, promoID)
	ret0, _ := ret[0].(models.Promo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPromoByPromoID indicates an expected call of GetPromoByPromoID.
func (mr *MockPromoRepositoryMockRecorder) GetPromoByPromoID(ctx, tx, promoID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPromoByPromoID", reflect.TypeOf((*MockPromoRepository)(nil).GetPromoByPromoID), ctx, tx, promoID)
}

// GetPromoByUserCart mocks base method.
func (m *MockPromoRepository) GetPromoByUserCart(ctx context.Context, tx *gorm.DB, input repository.GetPromoByUserCartInput) ([]models.Promo, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPromoByUserCart", ctx, tx, input)
	ret0, _ := ret[0].([]models.Promo)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetPromoByUserCart indicates an expected call of GetPromoByUserCart.
func (mr *MockPromoRepositoryMockRecorder) GetPromoByUserCart(ctx, tx, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPromoByUserCart", reflect.TypeOf((*MockPromoRepository)(nil).GetPromoByUserCart), ctx, tx, input)
}

// Save mocks base method.
func (m *MockPromoRepository) Save(ctx context.Context, tx *gorm.DB, promo *models.Promo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, tx, promo)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockPromoRepositoryMockRecorder) Save(ctx, tx, promo interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockPromoRepository)(nil).Save), ctx, tx, promo)
}

// SaveCities mocks base method.
func (m *MockPromoRepository) SaveCities(ctx context.Context, tx *gorm.DB, promoID uint, cities []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveCities", ctx, tx, promoID, cities)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveCities indicates an expected call of SaveCities.
func (mr *MockPromoRepositoryMockRecorder) SaveCities(ctx, tx, promoID, cities interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveCities", reflect.TypeOf((*MockPromoRepository)(nil).SaveCities), ctx, tx, promoID, cities)
}
