package verification

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pantame/server/apperror"
	"github.com/pantame/server/models/session"
	"github.com/pantame/server/models/verification"
	"github.com/pantame/server/utils/validators"
	"github.com/pantame/server/views"
)

func NewVerificationCode(c *fiber.Ctx) error {
	body := new(session.Login)

	if err := c.BodyParser(body); err != nil {
		return views.SendStatus(c, 400, apperror.InvalidData)
	}

	valid := validators.IsValidEmailByMX(body.GetAccessPass())
	if !valid {
		return views.SendStatus(c, 400, apperror.InvalidEmailAndPhone)
	}

	err := verification.SendVerificationCode(body.GetAccessPass())
	if err != nil {
		return views.SendError(c, err)
	}

	return views.SendSuccess(c)
}