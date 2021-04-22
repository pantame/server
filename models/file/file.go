package file

import (
	"github.com/pantame/server/apperror"
	"github.com/pantame/server/entities"
	"github.com/pantame/server/models/user"
	"github.com/pantame/server/storage"
	"github.com/pantame/server/utils"
	"gorm.io/gorm"
	"time"
)

func NewMetadata(userID uint64, metadata *entities.File) (string, error) {
	u, err := user.GetUser(userID)
	if err != nil {
		return "", err
	}

	if metadata.Size > u.AvailableLimit() {
		return "", apperror.NewError(nil, 400, "Você atingiu o limite de armazenamento.")
	}

	go user.AddToLimitUsed(userID, metadata.Size)

	hash, err := utils.NewUUID()
	if err != nil {
		return "", apperror.NewError(err, 500, apperror.InternalError)
	}

	file := entities.File{
		UUID:       hash.String(),
		UserID:     userID,
		Key:        metadata.Key,
		Metadata:   metadata.Metadata,
		Size:       metadata.Size,
		TotalParts: metadata.CalcTotalParts(),
		Register:   time.Now().Unix(),
	}

	err = storage.DB.Create(&file).Error
	if err != nil {
		return "", apperror.NewError(err, 500, apperror.InternalError)
	}
	return hash.String(), nil
}

func SumPart(UUID string) error {
	err := storage.DB.Table("files").Where("uuid = ?", UUID).Updates(map[string]interface{}{"parts_sent": gorm.Expr("parts_sent + ?", 1)}).Update("change", time.Now().Unix()).Error
	if err != nil {
		return apperror.NewError(err, 500, apperror.InternalError)
	}
	return nil
}

func GetMetadata(userID uint64, UUID string) (*entities.File, error) {
	var file entities.File

	err := storage.DB.Where("uuid = ? AND user_id = ?", UUID, userID).First(&file).Error
	if err != nil {
		if err.Error() == "record not found" {
			return nil, apperror.NewError(err, 404, apperror.NotFound)
		}
		return nil, apperror.NewError(err, 500, apperror.InternalError)
	}
	return &file, nil
}

func GetAllMetadata(userID uint64) ([]entities.File, error) {
	var files []entities.File

	err := storage.DB.Where("user_id = ?", userID).Find(&files).Error
	if err != nil {
		return nil, apperror.NewError(err, 500, apperror.InternalError)
	}
	return files, nil
}
