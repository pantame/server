package session

import (
	"fmt"
	"github.com/pantame/server/apperror"
	"github.com/pantame/server/entities"
	ipQuery "github.com/pantame/server/external/ip"
	"github.com/pantame/server/models/ip"
	"github.com/pantame/server/models/user"
	"github.com/pantame/server/models/verification"
	"github.com/pantame/server/storage"
	"github.com/pantame/server/utils"
	"strings"
	"time"
)

type LoginSuccess struct {
	Token     string         `json:"token"`
	User      *entities.User `json:"user"`
}

type Login struct {
	AccessPass string `json:"access_pass" binding:"required"`
}

type LoginValidation struct {
	Login
	Code string        `json:"code" binding:"required"`
}

func (d *Login) GetAccessPass() string {
	return strings.TrimSpace(d.AccessPass)
}

func SaveSessionInfo(data LoginSuccess, accessPass string, IP string, userAgent string) error {
	s := entities.Session{
		UserID:     data.User.ID,
		Token:      data.Token,
		Ip:         IP,
		AccessPass: accessPass,
		UserAgent:  userAgent,
		Register:   time.Now().Unix(),
	}

	ipData, err := ipQuery.Query(IP)
	if err == nil {
		ip.InsertIP(ipData)
	}

	err = storage.DB.Create(&s).Error
	if err != nil {
		return apperror.NewError(err, 500, apperror.InternalError)
	}
	return nil
}

func New(accessPass string) error {
	_, err := user.GetUserByAccessPass(accessPass)
	if err != nil {
		return err
	}
	return verification.SendVerificationCode(accessPass)
}

func ValidateCode(body *LoginValidation) (*LoginSuccess, error) {
	if err := verification.ValidateVerificationCode(body.GetAccessPass(), body.Code); err != nil {
		return nil, err
	}

	hash, err := utils.NewUUID()
	if err != nil {
		return nil, apperror.NewError(err, 500, apperror.InternalError)
	}

	var accessPass entities.AccessPass

	exist := storage.DB.Where("pass = ?", body.GetAccessPass()).First(&accessPass)
	if exist.Error != nil && exist.Error.Error() != "record not found" {
		return nil, apperror.NewError(exist.Error, 500, apperror.InternalError)
	}

	if exist.RowsAffected == 0 {
		return nil, apperror.NewError(nil, 404, apperror.NoAccessPass)
	}

	u := &entities.User{
		ID: accessPass.UserID,
	}

	err = storage.DB.First(u).Error
	if err != nil {
		return nil, apperror.NewError(err, 500, apperror.InternalError)
	}

	u.AccessPasses = append(u.AccessPasses, accessPass)

	if u == nil || u.ID <= 0 {
		return nil, apperror.NewError(nil, 500, apperror.MissingInternalInfo)
	}

	go storage.Cache.Set(hash.String(), fmt.Sprintf("%d@%d", u.ID, u.Level), 0)

	return &LoginSuccess{
		Token:     hash.String(),
		User:      u,
	}, nil
}

func LogoutByToken(userId uint64, token string) error {
	storage.Cache.Delete(token)

	r := storage.DB.Table("sessions").Where("user_id = ? AND token = ? AND active = ?", userId, token, true).Updates(map[string]interface{}{"active": false, "change": time.Now().Unix()})
	if r.Error != nil {
		return apperror.NewError(r.Error, 500, apperror.InternalError)
	}

	if r.RowsAffected == 0 {
		return apperror.NewError(nil, 404, apperror.NotFound)
	}

	return nil
}

func LogoutByID(userId, id uint64) error {
	var session entities.Session
	r := storage.DB.First(&session, "id = ? AND user_id = ? AND active = ?", id, userId, true)
	if r.Error != nil {
		if r.Error.Error() == "record not found" {
			return apperror.NewError(r.Error, 404, apperror.NotFound)
		}
		return apperror.NewError(r.Error, 500, apperror.InternalError)
	}

	return LogoutByToken(userId, session.Token)
}

func GetAllSessionsByUserID(userID uint64, active string) ([]entities.Session, error) {
	var sessions []entities.Session

	var err error
	if len(active) == 0 {
		err = storage.DB.Table("sessions").Omit("token").Where("user_id = ?", userID).Order("register DESC").Scan(&sessions).Error
	} else {
		err = storage.DB.Table("sessions").Omit("token").Where("user_id = ? AND active = ?", userID, active).Order("register DESC").Scan(&sessions).Error
	}

	if err != nil {
		return nil, apperror.NewError(err, 500, apperror.InternalError)
	}

	return sessions, nil
}
