package handler

import (
	"hangry/domain/dto"
	"hangry/generated"
	"hangry/utils"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
)

func validateCreateOrderRequest(req *generated.PostOrderJSONRequestBody) (dto.OrderInput, error) {

	err := validation.ValidateStruct(
		req,
		validation.Field(&req.UserId, validation.Required),
	)

	if err != nil {
		return dto.OrderInput{}, err
	}

	promoIds := []uint{}
	if req.PromoIds != nil {
		for _, promoId := range *req.PromoIds {
			promoIds = append(promoIds, uint(promoId))
		}
	}

	return dto.OrderInput{
		UserId:   uint(req.UserId),
		PromoIds: promoIds,
	}, nil
}

// PostOrder implements generated.ServerInterface.
func (s *Server) PostOrder(ctx echo.Context) error {
	req := generated.PostOrderJSONRequestBody{}
	if err := ctx.Bind(&req); err != nil {
		customError := utils.NewCustomError("failed to bind request", err, http.StatusBadRequest)
		return ResponseError(ctx, customError)
	}

	dto, err := validateCreateOrderRequest(&req)
	if err != nil {
		customError := utils.NewCustomError("validation error", err, http.StatusBadRequest)
		return ResponseError(ctx, customError)
	}

	err = s.orderUsecase.CreateOrder(ctx.Request().Context(), dto)
	if err != nil {
		return ResponseError(ctx, err)
	}

	return ctx.NoContent(http.StatusOK)
}
