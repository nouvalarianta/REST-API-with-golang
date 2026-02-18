package main

import (
	"golang-yt/dto"
	"golang-yt/internal/api"
	"golang-yt/internal/config"
	"golang-yt/internal/connection"
	"golang-yt/internal/repository"
	"golang-yt/internal/service"
	"net/http"

	jwtMid "github.com/gofiber/contrib/jwt"

	"github.com/gofiber/fiber/v2"
)

func main() {
	cnf := config.Get()
	dbConnection := connection.GetDatabase(cnf.DatabaseURL)

	app := fiber.New()

	jwtMidd := jwtMid.New(jwtMid.Config{
		SigningKey: jwtMid.SigningKey{Key: []byte(cnf.Jwt.Key)},
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(http.StatusUnauthorized).JSON(dto.CreateResponseError("autentikasi di perlukan"))
		},
	})

	customerRepository := repository.NewCustomer(dbConnection)
	userRepository := repository.NewUser(dbConnection)
	bookRepository := repository.NewBook(dbConnection)
	bookStockRepository := repository.NewBookStock(dbConnection)
	journalRepository := repository.NewJournal(dbConnection)
	mediaRepository := repository.Newmedia(dbConnection)
	chargeRepository := repository.NewCharge(dbConnection)

	customerService := service.NewCustomer(customerRepository)
	userService := service.NewAuth(cnf, userRepository)
	bookService := service.NewBook(cnf, bookRepository, bookStockRepository, mediaRepository)
	bookStockService := service.NewBookStock(bookRepository, bookStockRepository)
	journalService := service.NewJournal(journalRepository, bookRepository, bookStockRepository, customerRepository, chargeRepository)
	mediaService := service.Newmedia(cnf, mediaRepository)

	api.NewAuth(app, userService)
	api.NewCustomer(app, customerService, jwtMidd)
	api.NewBook(app, bookService, jwtMidd)
	api.NewBookStock(app, bookStockService, jwtMidd)
	api.NewJournal(app, journalService, jwtMidd)
	api.Newmedia(app, mediaService, jwtMidd, cnf)

	_ = app.Listen(cnf.Server.Host + ":" + cnf.Server.Port)
}

// func handler(ctx *fiber.Ctx) error{
// 	return ctx.Status(200).JSON("data")
// }
