package views

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pantame/server/apperror"
)

const Success = "OK"

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type ResponseData struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SendStatus(c *fiber.Ctx, status int, msg string) error {
	return c.Status(status).JSON(Response{
		Status:  status,
		Message: msg,
	})
}

func SendSuccess(c *fiber.Ctx) error {
	return c.Status(200).JSON(Response{
		Status:  200,
		Message: Success,
	})
}

func SendDataSuccess(c *fiber.Ctx, data interface{}) error {
	return c.Status(200).JSON(ResponseData{
		Status:  200,
		Message: Success,
		Data:    data,
	})
}

func SendError(c *fiber.Ctx, err error) error {
	appErr, ok := err.(apperror.AppError)
	if !ok {
		return c.Status(500).JSON(Response{
			Status:  500,
			Message: apperror.InternalError,
		})
	}

	return c.Status(appErr.StatusCode()).JSON(Response{
		Status:  appErr.StatusCode(),
		Message: appErr.Message(),
	})
}
