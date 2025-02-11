package handler

import (
	"hangry/domain/dto"
	"hangry/generated"
	"hangry/utils"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
)

func validationCreatePromoRequest(req *generated.CreatePromoRequest) (dto.CreatePromoInput, error) {
	err := validation.ValidateStruct(
		req,
		// name required
		validation.Field(&req.Name, validation.Required),
		// segmentation required and check enum
		validation.Field(&req.Segmentation, validation.Required, validation.In(generated.CreatePromoRequestSegmentationALL, generated.CreatePromoRequestSegmentationCITY, generated.CreatePromoRequestSegmentationLOYALUSER, generated.CreatePromoRequestSegmentationNEWUSER)),
		// startDate
		validation.Field(&req.StartDate, validation.Required),
		// endDate required and greater than startDate
		validation.Field(&req.EndDate, validation.Required, validation.Min(req.StartDate)),
		// type required and check enum
		validation.Field(&req.Type, validation.Required, validation.In(generated.CreatePromoRequestTypePERCENTAGEDISCOUNT, generated.CreatePromoRequestTypeBUYXGETYFREE)),
		// buyItemCount if type is BUYXGETYFREE required and greater than 0
		validation.Field(&req.BuyItemCount, validation.When(req.Type == generated.CreatePromoRequestTypeBUYXGETYFREE, validation.Required, validation.Min(1))),
		// freeItemCount if type is BUYXGETYFREE required and greater than 0
		validation.Field(&req.FreeItemCount, validation.When(req.Type == generated.CreatePromoRequestTypeBUYXGETYFREE, validation.Required, validation.Min(1))),
		// BuyProductId if type is BUYXGETYFREE required and greater than 0
		validation.Field(&req.BuyProductId, validation.When(req.Type == generated.CreatePromoRequestTypeBUYXGETYFREE, validation.Required, validation.Min(1))),
		// FreeProductId if type is BUYXGETYFREE required and greater than 0
		validation.Field(&req.FreeProductId, validation.When(req.Type == generated.CreatePromoRequestTypeBUYXGETYFREE, validation.Required, validation.Min(1))),
		// MaxUsageLimit if it exists required and greater than 0
		validation.Field(&req.MaxUsageLimit, validation.When(req.MaxUsageLimit != nil, validation.Min(1))),
		// DiscountValue if type is PERCENTAGEDISCOUNT required and greater than 0
		validation.Field(&req.DiscountValue, validation.When(req.Type == generated.CreatePromoRequestTypePERCENTAGEDISCOUNT, validation.Required, validation.Min(float32(1))), validation.Max(float32(100))),
		// MaxDiscountAmount if type is PERCENTAGEDISCOUNT and it exists required and greater than 0
		validation.Field(&req.MaxDiscountAmount, validation.When(req.Type == generated.CreatePromoRequestTypePERCENTAGEDISCOUNT && req.MaxDiscountAmount != nil, validation.Min(float32(1)))),
		// MinOrderAmount if it exists required and greater than 0
		validation.Field(&req.MinOrderAmount, validation.When(req.Type == generated.CreatePromoRequestTypePERCENTAGEDISCOUNT && req.MinOrderAmount != nil, validation.Min(float32(1)))),
		// Cities if segmentation is CITY required and not empty
		validation.Field(&req.Cities, validation.When(req.Segmentation == generated.CreatePromoRequestSegmentationCITY, validation.Required, validation.Length(1, 0))),
	)

	if err != nil {
		return dto.CreatePromoInput{}, err
	}

	dto := dto.CreatePromoInput{
		Name:         req.Name,
		Segmentation: string(req.Segmentation),
		Type:         string(req.Type),
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
	}

	if req.Description != nil {
		dto.Description = req.Description
	}

	if req.MinOrderAmount != nil {
		minOrderAmount := float64(*req.MinOrderAmount)
		dto.MinOrderAmount = &minOrderAmount
	}

	if req.DiscountValue != nil {
		discountValue := float64(*req.DiscountValue)
		dto.DiscountValue = &discountValue
	}

	if req.MaxDiscountAmount != nil {
		maxDiscountAmount := float64(*req.MaxDiscountAmount)
		dto.MaxDiscountAmount = &maxDiscountAmount
	}

	if req.BuyProductId != nil {
		buyProductId := int(*req.BuyProductId)
		dto.BuyProductId = &buyProductId
	}

	if req.FreeProductId != nil {
		freeProductId := int(*req.FreeProductId)
		dto.FreeProductId = &freeProductId
	}

	if req.BuyItemCount != nil {
		buyItemCount := int(*req.BuyItemCount)
		dto.BuyItemCount = &buyItemCount
	}

	if req.FreeItemCount != nil {
		freeItemCount := int(*req.FreeItemCount)
		dto.FreeItemCount = &freeItemCount
	}

	if req.MaxUsageLimit != nil {
		maxUsageLimit := int(*req.MaxUsageLimit)
		dto.MaxUsageLimit = &maxUsageLimit
	}

	if req.Cities != nil {
		dto.Cities = *req.Cities
	}

	return dto, nil
}

// PostPromo implements generated.ServerInterface.
func (s *Server) PostPromo(ctx echo.Context) error {
	req := generated.CreatePromoRequest{}
	if err := ctx.Bind(&req); err != nil {
		customError := utils.NewCustomError(err.Error(), nil, http.StatusBadRequest)
		return ResponseError(ctx, customError)
	}

	dto, err := validationCreatePromoRequest(&req)
	if err != nil {
		customErr := utils.NewCustomError("validation error", err, http.StatusBadRequest)
		return ResponseError(ctx, customErr)
	}

	promoId, err := s.promoUsecase.CreatePromo(ctx.Request().Context(), dto)
	if err != nil {
		return ResponseError(ctx, err)
	}

	return ctx.JSON(http.StatusOK, utils.NewResponse("promo created", echo.Map{"promoId": promoId}, nil))
}

func validationExtendPromoRequest(req *generated.ExtendPromoRequest) (dto.ExtendPromoInput, error) {
	err := validation.ValidateStruct(
		req,
		// endDate required and greater than startDate
		validation.Field(&req.EndDate,
			validation.Required,
			validation.When(!req.StartDate.IsZero(), validation.Min(req.StartDate)),
			validation.When(req.StartDate.IsZero(), validation.Min(time.Now())),
		),
	)

	dto := dto.ExtendPromoInput{
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}

	return dto, err
}

// PostPromoIdExtend implements generated.ServerInterface.
func (s *Server) PostPromoIdExtend(ctx echo.Context, id int) error {
	req := generated.ExtendPromoRequest{}
	if err := ctx.Bind(&req); err != nil {
		customError := utils.NewCustomError(err.Error(), nil, http.StatusBadRequest)
		return ResponseError(ctx, customError)
	}

	dto, err := validationExtendPromoRequest(&req)
	if err != nil {
		return ResponseError(ctx, err)
	}

	dto.ID = uint(id)

	err = s.promoUsecase.ExtendPromo(ctx.Request().Context(), dto)
	if err != nil {
		return ResponseError(ctx, err)
	}

	return ctx.JSON(http.StatusOK, utils.NewResponse("promo extended", nil, nil))
}

func validationGetPromoRequest(req *generated.GetGetPromoParams) (dto.GetPromoInput, error) {
	err := validation.ValidateStruct(
		req,
		validation.Field(&req.UserId, validation.Required),
	)

	if err != nil {
		return dto.GetPromoInput{}, err
	}

	if req.Page == nil || (req.Page != nil && *req.Page < 1) {
		page := 1
		req.Page = &page
	}

	if req.PerPage == nil || (req.PerPage != nil && *req.PerPage < 1) {
		perPage := 10
		req.PerPage = &perPage
	}

	return dto.GetPromoInput{
		UserId:  uint(req.UserId),
		Page:    *req.Page,
		PerPage: *req.PerPage,
	}, err
}

// GetGetPromo implements generated.ServerInterface.
func (s *Server) GetGetPromo(ctx echo.Context, params generated.GetGetPromoParams) error {
	dto, err := validationGetPromoRequest(&params)
	if err != nil {
		customError := utils.NewCustomError("validation error", err, http.StatusBadRequest)
		return ResponseError(ctx, customError)
	}

	promos, total, err := s.promoUsecase.GetPromo(ctx.Request().Context(), dto)
	if err != nil {
		return ResponseError(ctx, err)
	}

	meta := utils.BuildMeta(dto.Page, dto.PerPage, int(total))

	return ctx.JSON(http.StatusOK, utils.NewResponse("promo list", promos, meta))
}
