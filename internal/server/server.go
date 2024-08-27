package server

import (
	"fmt"
	"genai-test-go/internal/ai"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/tmc/langchaingo/llms"
)

var server = &FiberServer{
	App: fiber.New(fiber.Config{
		ServerHeader: "genai-test-go",
		AppName:      "genai-test-go",
	}),
}

type FiberServer struct {
	*fiber.App

	evaluatorModel llms.Model

	// to talk to OpenAI
	openAIChat *ai.Chat

	// to talk to Ollama as a local model
	ollamaChat *ai.Chat
}

func Run() error {
	server.RegisterFiberRoutes()
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	return server.Listen(fmt.Sprintf(":%d", port))
}
