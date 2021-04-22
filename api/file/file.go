package file

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/pantame/server/apperror"
	"github.com/pantame/server/config"
	"github.com/pantame/server/entities"
	"github.com/pantame/server/models/file"
	"github.com/pantame/server/views"
	"os"
	"path"
	"strconv"
	"strings"
)

func GetMetadata(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(uint64)

	m, err := file.GetMetadata(userID, ctx.Params("uuid"))
	if err != nil {
		return views.SendError(ctx, err)
	}

	return views.SendDataSuccess(ctx, m)
}

func GetAllMetadata(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(uint64)

	m, err := file.GetAllMetadata(userID)
	if err != nil {
		return views.SendError(ctx, err)
	}

	return views.SendDataSuccess(ctx, m)
}

func GetPart(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(uint64)

	m, err := file.GetMetadata(userID, ctx.Params("uuid"))
	if err != nil {
		return views.SendError(ctx, err)
	}

	n, err := strconv.ParseUint(ctx.Params("n"), 10, 64)
	if err != nil {
		return views.SendError(ctx, apperror.NewError(err, 400, apperror.InvalidData))
	}

	if !m.OwnsAllParts() {
		return views.SendError(ctx, apperror.NewError(nil, 400, apperror.IncompleteFileOnServer))
	}

	if n > m.TotalParts {
		return views.SendError(ctx, apperror.NewError(nil, 400, apperror.InvalidData))
	}

	filePath := path.Join(config.Paths().SaveDir, m.UUID, ctx.Params("n"))

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return views.SendError(ctx, apperror.NewError(err, 404, apperror.NotFound))
	}

	return ctx.Download(filePath)
}

func PostMetadata(ctx *fiber.Ctx) error {
	body := new(entities.File)

	if err := ctx.BodyParser(body); err != nil {
		return views.SendStatus(ctx, 400, apperror.InvalidData)
	}

	if len(body.Key) == 0 || len(body.Metadata) == 0 || body.Size == 0 {
		return views.SendStatus(ctx, 400, apperror.InvalidData)
	}

	userID, ok := ctx.Locals("userID").(uint64)
	if !ok {
		return views.SendStatus(ctx, 400, apperror.InvalidData)
	}

	id, err := file.NewMetadata(userID, body)
	if err != nil {
		return views.SendError(ctx, err)
	}

	data := struct {
		UUID string `json:"uuid"`
	}{
		UUID: id,
	}

	return views.SendDataSuccess(ctx, data)
}

func PostPart(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("userID").(uint64)
	if !ok {
		return views.SendStatus(ctx, 400, apperror.InvalidData)
	}

	m, err := file.GetMetadata(userID, ctx.FormValue("uuid"))
	if err != nil {
		return views.SendError(ctx, err)
	}

	if m.OwnsAllParts() {
		return views.SendError(ctx, apperror.NewError(nil, 403, apperror.AllParts))
	}

	f, err := ctx.FormFile("file")
	if err != nil {
		return views.SendError(ctx, apperror.NewError(err, 500, apperror.InternalError))
	}

	if !strings.Contains(f.Header.Get("content-type"), "octet-stream") {
		return views.SendStatus(ctx, 400, apperror.InvalidData)
	}

	if f.Size > 2500000 {
		return views.SendStatus(ctx, 400, apperror.MaxFileSize)
	}

	saveDir := path.Join(config.Paths().SaveDir, m.UUID)

	_, err = os.Stat(saveDir)
	if os.IsNotExist(err) {
		os.Mkdir(saveDir, os.ModePerm)
	}

	filename := path.Join(saveDir, fmt.Sprintf("%x", m.PartsSent))

	err = ctx.SaveFile(f, filename)
	if err != nil {
		return views.SendError(ctx, apperror.NewError(err, 500, apperror.InternalError))
	}

	err = file.SumPart(m.UUID)
	if err != nil {
		os.Remove(filename)
		return views.SendError(ctx, err)
	}
	return views.SendSuccess(ctx)
}
