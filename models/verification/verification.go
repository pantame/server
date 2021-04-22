package verification

import (
	"fmt"
	"github.com/pantame/server/apperror"
	"github.com/pantame/server/storage"
	"github.com/pantame/server/utils"
	"time"
)

func getAttempts(accessPass string) int {
	num, err := storage.Cache.GetInt(fmt.Sprintf("att-%s", accessPass))
	if err != nil {
		return -1
	}
	return num
}

func setAttempts(accessPass string, att int) {
	if att < 0 {
		storage.Cache.Set(fmt.Sprintf("att-%s", accessPass), 1, 3600*time.Second)
	} else {
		storage.Cache.IncrBy(fmt.Sprintf("att-%s", accessPass), 1)
	}
}

func labelForValCode(accessPass string) string {
	return "cod-" + accessPass
}

func SendVerificationCode(accessPass string) error {
	attempts := getAttempts(accessPass)
	// Controla a quantidade máxima de tentativas.
	if attempts >= 4 {
		return apperror.NewError(nil, 403, apperror.TooManyRequest)
	}

	go setAttempts(accessPass, attempts)

	code, err := utils.SendEmailVerificationCode("User", accessPass)
	if err != nil {
		return apperror.NewError(err, 500, "Não foi possível enviar o e-mail.")
	}

	err = storage.Cache.Set(labelForValCode(accessPass), code, 300*time.Second)
	if err != nil {
		return apperror.NewError(err, 500, apperror.InternalError)
	}
	return nil
}

func ValidateVerificationCode(accessPass, code string) error {
	r, err := storage.Cache.Get(labelForValCode(accessPass))
	if err != nil {
		return apperror.NewError(err, 428, apperror.PreconditionRequired)
	}

	if r != fmt.Sprint(code) {
		return apperror.NewError(nil, 401, apperror.InvalidCode)
	}

	go storage.Cache.Delete(labelForValCode(accessPass))

	return nil
}