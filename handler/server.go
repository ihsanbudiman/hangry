package handler

import (
	"hangry/generated"
	"hangry/usecase"

	"github.com/labstack/echo/v4"
)

type Server struct {
	cartUsecase  usecase.CartUsecase
	promoUsecase usecase.PromoUsecase
	orderUsecase usecase.OrderUsecase
}

// GetHealth implements generated.ServerInterface.
func (s *Server) GetHealth(ctx echo.Context) error {
	return ctx.JSON(200, "OK")
}

func NewServer(
	cartUsecase usecase.CartUsecase,
	promoUsecase usecase.PromoUsecase,
	orderUsecase usecase.OrderUsecase,
) generated.ServerInterface {
	return &Server{
		cartUsecase,
		promoUsecase,
		orderUsecase,
	}
}
