package user

import (
	"github.com/pantame/server/apperror"
	"github.com/pantame/server/entities"
	"github.com/pantame/server/storage"
	"gorm.io/gorm"
	"strings"
	"time"
)

func UpdateName(body *entities.User) error {
	user := entities.User{
		ID: body.ID,
	}
	err := storage.DB.Model(&user).Updates(entities.User{Name: body.Name, Change: time.Now().Unix()}).Error
	if err != nil {
		return apperror.NewError(err, 500, apperror.InternalError)
	}
	return nil
}

func exprLimitUsed(userID, size uint64, operation string) error {
	err := storage.DB.Model(&entities.User{
		ID: userID,
	}).Update("limit_used", gorm.Expr("limit_used "+operation+" ?", size)).Update("change", time.Now().Unix()).Error
	if err != nil {
		return apperror.NewError(err, 500, apperror.InternalError)
	}
	return nil
}

func AddToLimitUsed(userID, size uint64) error {
	return exprLimitUsed(userID, size, "+")
}

func SubtractFromLimitUsed(userID, size uint64) error {
	return exprLimitUsed(userID, size, "-")
}

func GetAccessPass(accessPass string) (*entities.AccessPass, error) {
	var ap entities.AccessPass
	err := storage.DB.Where("pass = ?", accessPass).First(&ap).Error
	if err != nil {
		if err.Error() != "record not found" {
			return nil, apperror.NewError(err, 500, apperror.InternalError)
		}
		return nil, apperror.NewError(err, 404, apperror.NotFound)
	}
	return &ap, nil
}

func GetUserByAccessPass(accessPass string) (*entities.User, error) {
	ap, err := GetAccessPass(accessPass)
	if err != nil {
		return nil, err
	}
	return GetUser(ap.UserID)
}

func GetUser(userID uint64) (*entities.User, error) {
	u := entities.User{
		ID: userID,
	}

	err := storage.DB.First(&u).Error
	if err != nil {
		if err.Error() != "record not found" {
			return nil, apperror.NewError(err, 500, apperror.InternalError)
		}
		return nil, apperror.NewError(err, 404, apperror.NotFound)
	}
	return &u, nil
}

func NewUser(user entities.User, accessPass, apType string) (*entities.User, error) {
	datetime := time.Now().Unix()

	u := &entities.User{
		Username: user.Username,
		Name:     user.Name,
		Key:      user.Key,
		Level:    user.Level,
		Limit:    user.Limit,
		Register: datetime,
		AccessPasses: []entities.AccessPass{
			{
				Pass:     accessPass,
				Type:     apType,
				Register: datetime,
			},
		},
	}
	err := storage.DB.Create(u).Error
	if err != nil {
		return nil, apperror.NewError(err, 500, apperror.InternalError)
	}
	return u, nil
}

func GetAccessPasses(userID uint64) ([]entities.AccessPass, error) {
	var accessPasses []entities.AccessPass

	err := storage.DB.Find(&accessPasses, "user_id = ?", userID).Error
	if err != nil {
		return nil, apperror.NewError(err, 500, apperror.InternalError)
	}

	return accessPasses, nil
}

func NewAccessPass(userID uint64, pass, apType string) error {
	accessPass := entities.AccessPass{
		UserID:   userID,
		Pass:     pass,
		Type:     apType,
		Register: time.Now().Unix(),
	}
	err := storage.DB.Create(&accessPass).Error
	if err != nil {
		if strings.Contains(err.Error(), "unique") {
			return apperror.NewError(err, 409, apperror.AlreadyExisting)
		} else if strings.Contains(err.Error(), "violates foreign key constraint") {
			return apperror.NewError(err, 404, apperror.UserNotFount)
		}
		return apperror.NewError(err, 500, apperror.InternalError)
	}

	return nil
}

func RemoveAccessPass(userID uint64, pass string) error {
	var n int64
	storage.DB.Table("access_passes").Where("user_id = ? AND pass != ?", userID, pass).Count(&n)

	if n == 0 {
		return apperror.NewError(nil, 401, "Você deve ter pelo menos um passe de acesso.")
	}

	r := storage.DB.Table("access_passes").Delete(entities.AccessPass{}, "user_id = ? AND pass = ?", userID, pass)
	if r.Error != nil {
		return apperror.NewError(r.Error, 500, apperror.InternalError)
	}

	if r.RowsAffected == 0 {
		return apperror.NewError(nil, 404, apperror.NotFound)
	}

	return nil
}
