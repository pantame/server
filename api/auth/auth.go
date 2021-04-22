package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pantame/server/apperror"
	"github.com/pantame/server/storage"
	"github.com/pantame/server/views"
	"strconv"
	"strings"
)

func Auth(c *fiber.Ctx) error {
	token := c.Get("token")
	if len(token) != 36 {
		return views.SendStatus(c, 401, apperror.Unauthorized)
	}

	d, err := storage.Cache.Get(token)
	if err != nil {
		return views.SendStatus(c, 401, apperror.Unauthorized)
	}

	info := strings.Split(d, "@")
	if len(info) != 2 {
		return views.SendStatus(c, 500, apperror.InternalError)
	}

	userID, err := strconv.ParseUint(info[0], 10, 64)
	if err != nil {
		return views.SendStatus(c, 500, apperror.InternalError)
	}

	level, err := strconv.ParseUint(info[1], 10, 64)
	if err != nil {
		return views.SendStatus(c, 500, apperror.InternalError)
	}

	c.Locals("userId", userID)
	c.Locals("level", level)

	return c.Next()
}
