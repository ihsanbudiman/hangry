package usecase

import (
	"context"
	"hangry/domain/dto"
	"hangry/repository"
	"testing"
)

func Test_cartUsecase_RemoveFromCart(t *testing.T) {
	type fields struct {
		transaction       repository.TransactionRepository
		cartRepository    repository.CartRepository
		productRepository repository.ProductRepository
	}
	type args struct {
		ctx context.Context
		dto dto.RemoveFromCartInput
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cartUsecase{
				transaction:       tt.fields.transaction,
				cartRepository:    tt.fields.cartRepository,
				productRepository: tt.fields.productRepository,
			}
			if err := c.RemoveFromCart(tt.args.ctx, tt.args.dto); (err != nil) != tt.wantErr {
				t.Errorf("cartUsecase.RemoveFromCart() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
