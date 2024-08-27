package server

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

var server = &FiberServer{
	App: fiber.New(fiber.Config{
		ServerHeader: "genai-test-go",
		AppName:      "genai-test-go",
	}),
}

type FiberServer struct {
	*fiber.App

	// using RAG for conversational retrieval
	conversationalRetrieval chains.ConversationalRetrievalQA

	// to talk to the LLM directly
	llm llms.Model

	// to talk to Ollama as a local model
	ollamaModel *ollama.LLM

	// using RAG for conversational retrieval with Ollama
	ollamaConversationalRetrieval chains.ConversationalRetrievalQA
}

func Run() error {
	server.RegisterFiberRoutes()
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	return server.Listen(fmt.Sprintf(":%d", port))
}
