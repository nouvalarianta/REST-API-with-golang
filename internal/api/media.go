package api

import (
	"context"
	"golang-yt/domain"
	"golang-yt/dto"
	"golang-yt/internal/config"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type mediaApi struct {
	mediaService domain.MediaService
	cnf *config.Config
}

func Newmedia(app *fiber.App, mediaService domain.MediaService, authMid fiber.Handler, cnf *config.Config) {
	ma := mediaApi{
		mediaService: mediaService,
		cnf: cnf,
	}

	media := app.Group("/media", authMid)
 
	media.Post("/", ma.Create)
	media.Static("/", cnf.Storage.BasePath)
}
 
func (ma mediaApi) Create(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10 * time.Second)
	defer cancel()

	file, err := ctx.FormFile("media")
	if err != nil {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	filename := uuid.NewString() + file.Filename
	path := filepath.Join(ma.cnf.Storage.BasePath, filename)
	err = ctx.SaveFile(file, path)

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.CreateResponseError(err.Error()))
	}

	res, err := ma.mediaService.Create(c, dto.CreateMediaRequest{
		 Path: filename,
	})

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.CreateResponseError(err.Error()))
	}

	return ctx.Status(http.StatusCreated).JSON(dto.CreateResponseSuccess(res))
}