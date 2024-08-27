package server

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
)

const question string = "Since which Testcontainers for Go version is the Grafana LGTM module available?"

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

func parseTemperature(c *fiber.Ctx) float64 {
	var temperature float64
	// read it from the query string
	if temp := c.Query("t"); temp != "" {
		t, err := strconv.ParseFloat(temp, 64)
		if err != nil {
			return 0.0
		}
		temperature = t
	}

	if temperature > 1.0 {
		temperature = 1.0
	} else if temperature < 0.0 {
		temperature = 0.0
	}

	return temperature
}

func (s *FiberServer) RagHandler(c *fiber.Ctx) error {
	response, err := chains.Run(c.Context(), s.conversationalRetrieval, question, chains.WithTemperature(parseTemperature(c)))
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
		llms.TextParts(llms.ChatMessageTypeHuman, question),
	}, llms.WithTemperature(parseTemperature(c)))
	if err != nil {
		log.Fatal(err)
	}

	text := strings.ReplaceAll(response.Choices[0].Content, "\"", "'")

	resp := fiber.Map{
		"message": text,
	}

	return c.JSON(resp)
}
