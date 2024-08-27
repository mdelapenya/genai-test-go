package server

import (
	"github.com/gofiber/fiber/v2"
)

func (s *FiberServer) OpenAIRagHandler(c *fiber.Ctx) error {
	return ragHandler(c, s.openAIChat.ConversationalRetrieval)
}

func (s *FiberServer) OpenAILLHandler(c *fiber.Ctx) error {
	return llmHandler(c, s.openAIChat.Model)
}
