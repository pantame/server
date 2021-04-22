package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pantame/server/apperror"
	"github.com/pantame/server/entities"
	"github.com/pantame/server/models/session"
	"github.com/pantame/server/models/user"
	"github.com/pantame/server/models/verification"
	"github.com/pantame/server/utils/validators"
	"github.com/pantame/server/views"
)

func UpdateNameUser(c *fiber.Ctx) error {
	body := new(entities.User)
	body.ID = c.Locals("userId").(uint64)

	if err := c.BodyParser(body); err != nil {
		return views.SendStatus(c, 400, apperror.InvalidData)
	}

	if len(body.Name) > 120 {
		return views.SendStatus(c, 400, apperror.InvalidData)
	}

	err := user.UpdateName(body)
	if err != nil {
		return views.SendError(c, err)
	}
	return views.SendSuccess(c)
}

func GetUser(c *fiber.Ctx) error {
	userID := c.Locals("userId").(uint64)

	u, err := user.GetUser(userID)
	if err != nil {
		return views.SendError(c, err)
	}

	return views.SendDataSuccess(c, u)
}

func GetAccessPasses(c *fiber.Ctx) error {
	userID := c.Locals("userId").(uint64)

	accessPasses, err := user.GetAccessPasses(userID)
	if err != nil {
		return views.SendError(c, err)
	}

	return views.SendDataSuccess(c, accessPasses)
}

func NewAccessPass(c *fiber.Ctx) error {
	body := new(session.LoginValidation)

	if err := c.BodyParser(body); err != nil {
		return views.SendStatus(c, 400, apperror.InvalidData)
	}

	userID := c.Locals("userId").(uint64)

	ok := validators.IsValidEmailByMX(body.GetAccessPass())
	if !ok {
		return views.SendStatus(c, 400, apperror.InvalidEmailAndPhone)
	}

	if !validators.IsValidString(body.Code, 6, 6) {
		return views.SendStatus(c, 400, apperror.InvalidCode)
	}

	if err := verification.ValidateVerificationCode(body.GetAccessPass(), body.Code); err != nil {
		return views.SendError(c, err)
	}

	err := user.NewAccessPass(userID, body.GetAccessPass(), "mail")
	if err != nil {
		return views.SendError(c, err)
	}

	return views.SendSuccess(c)
}

func RemoveAccessPass(c *fiber.Ctx) error {
	body := new(session.Login)

	if err := c.BodyParser(body); err != nil {
		return views.SendStatus(c, 400, apperror.InvalidData)
	}

	userID := c.Locals("userId").(uint64)

	ok := validators.IsValidEmailByMX(body.GetAccessPass())
	if !ok {
		return views.SendStatus(c, 400, apperror.InvalidEmailAndPhone)
	}

	err := user.RemoveAccessPass(userID, body.GetAccessPass())
	if err != nil {
		return views.SendError(c, err)
	}

	return views.SendSuccess(c)
}