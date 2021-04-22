package ip

import (
	"github.com/pantame/server/apperror"
	"github.com/pantame/server/entities"
	"github.com/pantame/server/storage"
	"strings"
	"time"
)

func GenerateIPDate(IP string) string {
	return IP + "-" + time.Now().Format("2006-01-02")
}

func InsertIP(IP *entities.IPData) error {
	err := storage.DB.Create(&IP).Error
	if err != nil {
		if !strings.Contains(err.Error(), "unique") {
			return apperror.NewError(err, 500, apperror.InternalError)
		}
	}
	return nil
}

func GetIP(where, IP string) (*entities.IPData, error) {
	var ip entities.IPData

	r := storage.DB.First(&ip, where + " = ?", IP)
	if r.Error != nil {
		if r.Error.Error() == "record not found" {
			return nil, apperror.NewError(r.Error, 404, apperror.NotFound)
		}
		return nil, apperror.NewError(r.Error, 500, apperror.InternalError)
	}

	return &ip, nil
}
