package server

import (
	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/rs/zerolog"
	"time"
)

type Config struct {
	AppName      string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	Logger       zerolog.Logger
}

func New(cfg Config) *fiber.App {
	app := fiber.New(fiber.Config{
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
		AppName:      cfg.AppName,
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
	})

	app.Use(cors.New())
	app.Use(requestid.New(requestid.Config{
		Next:       nil,
		Header:     fiber.HeaderXRequestID,
		Generator:  utils.UUID,
		ContextKey: "requestID",
	}))
	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: &cfg.Logger,
		Fields: []string{fiberzerolog.FieldLatency, fiberzerolog.FieldStatus, fiberzerolog.FieldMethod, fiberzerolog.FieldURL, fiberzerolog.FieldError, fiberzerolog.FieldRequestID},
	}))

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("pong")
	})

	return app
}
