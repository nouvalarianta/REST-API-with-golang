package api

import (
	"context"
	"errors"
	"golang-yt/domain"
	"golang-yt/dto"
	"golang-yt/internal/util"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type journalApi struct {
	journalService domain.JournalService
}

func NewJournal(app *fiber.App, journalService domain.JournalService, authMid fiber.Handler) {
	ja := journalApi{
		journalService: journalService,
	}

	journal := app.Group("/journals", authMid)

	journal.Get("/", ja.Index)
	journal.Post("/", ja.Create)
	journal.Put("/:id", ja.Update)
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
		return ctx.Status(http.StatusUnprocessableEntity).
			JSON(dto.CreateResponseError("invalid request body"))
	}

	fails := util.Validate(req)
	if len(fails) > 0 {
		return ctx.Status(http.StatusBadRequest).JSON(dto.CreateResponseErrorData("validate error", fails))
	}

	if err := ja.journalService.Create(c, req); err != nil {
		if errors.Is(err, domain.BookNotFound) || errors.Is(err, domain.BookStockNotFound) {
			return ctx.Status(http.StatusNotFound).JSON(dto.CreateResponseError(err.Error()))
		}
		if errors.Is(err, domain.BookAlreadyBorrowed) {
			return ctx.Status(http.StatusBadRequest).JSON(dto.CreateResponseError(err.Error()))
		}
		return ctx.Status(http.StatusInternalServerError).JSON(dto.CreateResponseError(err.Error()))
	}

	return ctx.Status(http.StatusCreated).JSON(dto.CreateResponseSuccess(""))
}

func (ja journalApi) Update(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	id := ctx.Params("id")

	userToken, ok := ctx.Locals("user").(*jwt.Token)
	if !ok || userToken == nil {
		return ctx.Status(http.StatusUnauthorized).JSON(dto.CreateResponseError("autentikasi diperlukan"))
	}

	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok {
		return ctx.Status(http.StatusUnauthorized).JSON(dto.CreateResponseError("autentikasi diperlukan"))
	}

	userIDVal, ok := claims["id"]
	if !ok {
		return ctx.Status(http.StatusUnauthorized).JSON(dto.CreateResponseError("autentikasi diperlukan"))
	}

	userID, ok := userIDVal.(string)
	if !ok {
		return ctx.Status(http.StatusUnauthorized).JSON(dto.CreateResponseError("autentikasi diperlukan"))
	}

	if err := ja.journalService.Return(c, dto.ReturnJournalRequest{
		JournalID: id,
		UserId:    userID,
	}); err != nil {
		if errors.Is(err, domain.JournalNotFound) {
			return ctx.Status(http.StatusNotFound).JSON(dto.CreateResponseError(err.Error()))
		}
		if errors.Is(err, domain.JournalAlreadyCompleted) {
			return ctx.Status(http.StatusBadRequest).JSON(dto.CreateResponseError(err.Error()))
		}
		return ctx.Status(http.StatusInternalServerError).JSON(dto.CreateResponseError(err.Error()))
	}

	return ctx.Status(http.StatusOK).JSON(dto.CreateResponseSuccess(""))
}
