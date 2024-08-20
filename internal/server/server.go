package server

import (
	"github.com/gofiber/fiber/v2"

	"sampleapp/internal/database"
)

type FiberServer struct {
	*fiber.App

	db database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "sampleapp",
			AppName:      "sampleapp",
		}),

		db: database.New(),
	}

	return server
}
