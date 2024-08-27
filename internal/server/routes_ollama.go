package server

import (
	"github.com/gofiber/fiber/v2"
)

func (s *FiberServer) OllamaLLHandler(c *fiber.Ctx) error {
	return llmHandler(c, s.ollamaChat.Model)
}

func (s *FiberServer) OllamaRagHandler(c *fiber.Ctx) error {
	return ragHandler(c, s.ollamaChat.ConversationalRetrieval)
}
