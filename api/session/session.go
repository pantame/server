package session

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pantame/server/apperror"
	"github.com/pantame/server/entities"
	"github.com/pantame/server/models/ip"
	"github.com/pantame/server/models/session"
	"github.com/pantame/server/utils/validators"
	"github.com/pantame/server/views"
	"time"
)

func NewSession(c *fiber.Ctx) error {
	body := new(session.Login)

	if err := c.BodyParser(body); err != nil {
		return views.SendStatus(c, 400, apperror.InvalidData)
	}

	valid := validators.IsValidEmailByMX(body.GetAccessPass())
	if !valid {
		return views.SendStatus(c, 400, apperror.InvalidEmailAndPhone)
	}

	err := session.New(body.GetAccessPass())
	if err != nil {
		return views.SendError(c, err)
	}

	return views.SendSuccess(c)
}

func Validate(c *fiber.Ctx) error {
	body := new(session.LoginValidation)

	if err := c.BodyParser(body); err != nil {
		return views.SendStatus(c, 400, apperror.InvalidData)
	}

	valid := validators.IsValidEmailByMX(body.GetAccessPass())
	if !valid {
		return views.SendStatus(c, 400, apperror.InvalidEmailAndPhone)
	}

	if !validators.IsValidString(body.Code, 6, 6) {
		return views.SendStatus(c, 400, apperror.InvalidCode)
	}

	s, err := session.ValidateCode(body)
	if err != nil {
		return views.SendError(c, err)
	}

	go session.SaveSessionInfo(*s, body.GetAccessPass(), c.IP(), c.Get("user-agent"))

	return views.SendDataSuccess(c, s)
}

func Logout(c *fiber.Ctx) error {
	userId := c.Locals("userId").(uint64)

	var token string
	var id uint64

	if len(c.Body()) == 0 {
		token = c.Get("token")
	} else {
		body := new(entities.Session)
		if err := c.BodyParser(body); err == nil {
			token = body.Token
			id = body.ID
		}
	}

	var err error
	if id != 0 {
		err = session.LogoutByID(userId, id)
	} else if len(token) != 0 {
		err = session.LogoutByToken(userId, token)
	} else {
		return views.SendStatus(c, 400, apperror.InvalidData)
	}

	if err != nil {
		return views.SendError(c, err)
	}

	return views.SendSuccess(c)
}

func GetAllActiveSessionsByUserID(c *fiber.Ctx) error {
	sessions, err := session.GetAllSessionsByUserID(c.Locals("userId").(uint64), "true")
	if err != nil {
		return views.SendError(c, err)
	}

	data := make([]map[string]interface{}, 0, len(sessions))

	for _, s := range sessions {
		currentSession := false
		if s.Token == c.Get("token") {
			currentSession = true
		}

		IPData := map[string]string{
			"ip": s.Ip,
		}

		IP, err := ip.GetIP("ip_date", s.Ip+"-"+time.Unix(s.Register, 0).Format("2006-01-02"))
		if err == nil {
			IPData = map[string]string{
				"ip":           s.Ip,
				"country":      IP.Country,
				"country_code": IP.CountryCode,
				"region":       IP.Region,
				"region_name":  IP.RegionName,
				"city":         IP.City,
				"district":     IP.District,
				"zip":          IP.Zip,
			}
		}

		data = append(data, map[string]interface{}{
			"id":              s.ID,
			"access_pass":     s.AccessPass,
			"ip":              IPData,
			"user_agent":      s.UserAgent,
			"current_session": currentSession,
			"register":        s.Register,
			"change":          s.Change,
		})
	}

	return views.SendDataSuccess(c, data)
}
