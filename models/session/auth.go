package session

import (
	"github.com/pantame/server/apperror"
	"github.com/pantame/server/entities"
	"github.com/pantame/server/models/user"
	"github.com/pantame/server/storage"
)

func GetSessionByToken(token string) (*entities.Session, error) {
	var s entities.Session
	err := storage.DB.Where("token = ?", token).First(&s).Error
	if err != nil {
		if err.Error() != "record not found" {
			return nil, apperror.NewError(err, 500, apperror.InternalError)
		}
		return nil, apperror.NewError(err, 404, apperror.NotFound)
	}

	return &s, nil
}

func CheckPermissionLevel(token string, level uint64) error {
	s, err := GetSessionByToken(token)
	if err != nil {
		return err
	}

	if !s.Active {
		return apperror.NewError(nil, 401, apperror.InvalidSession)
	}

	u, err := user.GetUser(s.UserID)
	if err != nil {
		return err
	}

	if u.Level < level {
		return apperror.NewError(err, 401, apperror.PermissionError)
	}
	return nil
}
