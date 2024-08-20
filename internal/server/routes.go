package server

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/tmc/langchaingo/chains"
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
	response, err := chains.Run(c.Context(), s.conversationalRetrieval, "¿Qué es un TTV?")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	resp := fiber.Map{
		"message-from-rag": response,
	}

	return c.JSON(resp)
}

func (s *FiberServer) LLHandler(c *fiber.Ctx) error {
	completion, err := s.llm.Call(c.Context(), "¿Qué es un TTV?")
	if err != nil {
		log.Fatal(err)
	}

	resp := fiber.Map{
		"message-from-llm": completion,
	}

	return c.JSON(resp)
}
