package controller

import (
	"context"
	"github.com/labstack/echo/v4"
)

func Ctx(ctx echo.Context) context.Context {
	return ctx.Request().Context()
}
