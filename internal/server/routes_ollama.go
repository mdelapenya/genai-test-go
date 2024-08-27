package server

import (
	"github.com/gofiber/fiber/v2"
)

func (s *FiberServer) OllamaLLHandler(c *fiber.Ctx) error {
	return llmHandler(c, s.ollamaModel)
}

func (s *FiberServer) OllamaRagHandler(c *fiber.Ctx) error {
	return ragHandler(c, s.ollamaConversationalRetrieval)
}
