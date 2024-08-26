package server

import (
	"log"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
)

func (s *FiberServer) RegisterFiberRoutes() {
	s.App.Get("/", s.HelloWorldHandler)

	chatApis := s.App.Group("/chat")

	chatApis.Add(http.MethodGet, "/rag", s.RagHandler)
	chatApis.Add(http.MethodGet, "/llm", s.LLHandler)
}

func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "Hello World",
	}

	return c.JSON(resp)
}

func (s *FiberServer) RagHandler(c *fiber.Ctx) error {
	response, err := chains.Run(c.Context(), s.conversationalRetrieval, "¿Qué es un TTV?", chains.WithTemperature(0.5))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response = strings.ReplaceAll(response, "\"", "'")

	resp := fiber.Map{
		"message": response,
	}

	return c.JSON(resp)
}

func (s *FiberServer) LLHandler(c *fiber.Ctx) error {
	response, err := s.llm.GenerateContent(c.Context(), []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeHuman, "¿Qué es un TTV?"),
	}, llms.WithTemperature(0.5))
	if err != nil {
		log.Fatal(err)
	}

	text := strings.ReplaceAll(response.Choices[0].Content, "\"", "'")

	resp := fiber.Map{
		"message": text,
	}

	return c.JSON(resp)
}
