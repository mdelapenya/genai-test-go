package server

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/pgvector"

	"genai-test-go/internal/database"
)

type FiberServer struct {
	*fiber.App

	db    database.Service
	store vectorstores.VectorStore
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "genai-test-go",
			AppName:      "genai-test-go",
		}),

		db: database.New(),
	}

	// Create an embeddings client using the OpenAI API. Requires environment variable OPENAI_API_KEY to be set.
	llm, err := openai.New()
	if err != nil {
		log.Fatal(err)
	}

	e, err := embeddings.NewEmbedder(llm)
	if err != nil {
		log.Fatal(err)
	}

	server.Hooks().OnListen(func(listenData fiber.ListenData) error {
		if fiber.IsChild() {
			return nil
		}

		connString := database.MustVectorDatabase(context.Background())

		store, err := pgvector.New(
			context.Background(),
			pgvector.WithConnectionURL(connString),
			pgvector.WithEmbedder(e),
		)
		if err != nil {
			log.Fatal(err)
		}

		server.store = store

		return nil
	})

	return server
}
