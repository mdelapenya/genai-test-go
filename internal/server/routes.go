package server

import (
	"genai-test-go/internal/ai"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
)

const (
	question string = "Since which Testcontainers for Go version is the Grafana LGTM module available?"

	// Using must/should is important
	reference = `- Answer must not mention any other module
- Answer must mention the version of Testcontainers for Go, which is v0.33.0
- Answer must be less than 5 sentences`
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
	answer, err := chains.Run(c.Context(), s.conversationalRetrieval, question, chains.WithTemperature(parseTemperature(c)))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	answer = strings.ReplaceAll(answer, "\"", "'")

	evaluator := ai.NewEvaluator(server.llm)
	aiResp, err := evaluator.Evaluate(c.Context(), question, answer, reference)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	resp := fiber.Map{
		"evaluator": aiResp,
		"answer":    answer,
		"question":  question,
		"reference": reference,
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

	answer := strings.ReplaceAll(response.Choices[0].Content, "\"", "'")

	evaluator := ai.NewEvaluator(server.llm)
	aiResp, err := evaluator.Evaluate(c.Context(), question, answer, reference)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	resp := fiber.Map{
		"evaluator": aiResp,
		"answer":    answer,
		"question":  question,
		"reference": reference,
	}

	return c.JSON(resp)
}
