package utils

import "github.com/gofiber/fiber"

func JSON(ctx *fiber.Ctx, statusCode int, data interface{}) {
	ctx.Status(statusCode)
	err := ctx.JSON(data)
	if err != nil {
		ctx.JSON(err)
	}
}

func Error(ctx *fiber.Ctx, statusCode int, err error) {
	JSON(ctx, statusCode, struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	})
}
