package dto

import (
	"github.com/gofiber/fiber/v2"
)

func RespondWithError(c *fiber.Ctx, statusCode int, errorMsg string) error {
	return c.Status(statusCode).JSON(Response{
		Success: false,
		Error:   errorMsg,
	})
}

func RespondWithJSON(c *fiber.Ctx, statusCode int, data interface{}) error {
	return c.Status(statusCode).JSON(Response{
		Success: true,
		Data:    data,
	})
}
