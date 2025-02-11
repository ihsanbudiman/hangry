package handler

import (
	"hangry/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func ResponseError(ctx echo.Context, err error) error {
	errData, ok := err.(*utils.CustomError)

	if !ok {
		errData = utils.NewCustomError(err.Error(), nil, http.StatusInternalServerError)
	}

	return ctx.JSON(errData.StatusCode, errData)
}
