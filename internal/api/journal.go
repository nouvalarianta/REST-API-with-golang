package api

import (
	"context"
	"golang-yt/domain"
	"golang-yt/dto"
	"golang-yt/internal/util"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type journalApi struct {
	journalService domain.JournalService
}

func NewJournal(app *fiber.App, journalService domain.JournalService, authMid fiber.Handler) {
	ja := journalApi{
		journalService: journalService,
	}

	app.Get("/journals", authMid, ja.Index)
	app.Post("/journals", authMid, ja.Create)
	app.Put("/journals/:id", authMid, ja.Update)
}

func (ja journalApi) Index(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	customerId := ctx.Query("customer_id")
	status := ctx.Query("status")
	res, err := ja.journalService.Index(c, domain.JournalSearch{
		CustomerId: customerId,
		Status:     status,
	})

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.CreateResponseError(err.Error()))
	}

	return ctx.Status(http.StatusOK).JSON(dto.CreateResponseSuccess(res))
}

func (ja journalApi) Create(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.CreateJournalRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.SendStatus(http.StatusUnprocessableEntity)
	}

	fails := util.Validate(req)
	if len(fails) > 0 {
		return ctx.Status(http.StatusBadRequest).JSON(dto.CreateResponseErrorData("validate error", fails))
	}

	err := ja.journalService.Create(c, req)

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.CreateResponseError(err.Error()))
	}

	return ctx.Status(http.StatusCreated).JSON(dto.CreateResponseSuccess(""))
}

func (ja journalApi) Update(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	id := ctx.Params("id")

	err := ja.journalService.Return(c, dto.ReturnJournalRequest{JournaID: id})

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dto.CreateResponseError(err.Error()))
	}

	return ctx.Status(http.StatusCreated).JSON(dto.CreateResponseSuccess(""))
}
