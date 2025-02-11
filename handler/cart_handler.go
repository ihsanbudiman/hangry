package handler

import (
	"hangry/domain/dto"
	"hangry/generated"
	"hangry/utils"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
)

func validationAddCartRequest(req *generated.AddCartRequest) (dto.AddToCartInput, error) {
	err := validation.ValidateStruct(
		req,
		validation.Field(&req.ProductId, validation.Required),
		validation.Field(&req.Quantity, validation.Required, validation.Min(1)),
		validation.Field(&req.UserId, validation.Required),
	)

	if err != nil {
		return dto.AddToCartInput{}, err
	}

	return dto.AddToCartInput{
		ProductId: uint(req.ProductId),
		Quantity:  req.Quantity,
		UserId:    uint(req.UserId),
	}, nil

}

// PostAddCart implements generated.ServerInterface.
func (s *Server) PostAddCart(ctx echo.Context) error {
	req := generated.AddCartRequest{}
	if err := ctx.Bind(&req); err != nil {
		customErr := utils.NewCustomError(err.Error(), nil, http.StatusBadRequest)
		return ResponseError(ctx, customErr)
	}

	dto, err := validationAddCartRequest(&req)
	if err != nil {
		customErr := utils.NewCustomError("validation error", err, http.StatusBadRequest)
		return ResponseError(ctx, customErr)
	}

	err = s.cartUsecase.AddToCart(ctx.Request().Context(), dto)
	if err != nil {
		return ResponseError(ctx, err)
	}

	return ctx.JSON(http.StatusOK, generated.SuccessResponse{
		Message: "success",
	})
}

func validationRemoveCartRequest(req *generated.RemoveFromCartRequest) (dto.RemoveFromCartInput, error) {
	err := validation.ValidateStruct(
		req,
		validation.Field(&req.ProductId, validation.Required),
		validation.Field(&req.UserId, validation.Required),
	)

	if err != nil {
		return dto.RemoveFromCartInput{}, err
	}

	return dto.RemoveFromCartInput{
		ProductId: uint(req.ProductId),
		UserId:    uint(req.UserId),
	}, nil
}

// PostRemoveFromCart implements generated.ServerInterface.
func (s *Server) PostRemoveFromCart(ctx echo.Context) error {
	req := generated.RemoveFromCartRequest{}
	if err := ctx.Bind(&req); err != nil {
		customError := utils.NewCustomError(err.Error(), nil, http.StatusBadRequest)
		return ResponseError(ctx, customError)
	}

	dto, err := validationRemoveCartRequest(&req)
	if err != nil {
		customError := utils.NewCustomError("validation error", err, http.StatusBadRequest)
		return ResponseError(ctx, customError)
	}

	err = s.cartUsecase.RemoveFromCart(ctx.Request().Context(), dto)
	if err != nil {
		return ResponseError(ctx, err)
	}

	return ctx.JSON(http.StatusOK, generated.SuccessResponse{
		Message: "success",
	})
}
