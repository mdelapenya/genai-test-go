//go:build local_dev
// +build local_dev

package server

import (
	"context"
	"genai-test-go/internal/database"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/memory"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/pgvector"
)

func init() {
	godotenv.Load("../../.env")

	log.Println("Initializing Local Development Environment")
	// Create an embeddings client using the OpenAI API. Requires environment variable OPENAI_API_KEY to be set.
	llm, err := openai.New()
	if err != nil {
		log.Fatal(err)
	}

	server.llm = llm

	e, err := embeddings.NewEmbedder(llm)
	if err != nil {
		log.Fatal(err)
	}

	optionsVector := []vectorstores.Option{
		vectorstores.WithScoreThreshold(0.80), // use for precision, when you want to get only the most relevant documents
		//vectorstores.WithNameSpace(""),            // use for set a namespace in the storage
		//vectorstores.WithFilters(map[string]interface{}{"language": "en"}), // use for filter the documents
		//vectorstores.WithEmbedder(embedder), // use when you want add documents or doing similarity search
		//vectorstores.WithDeduplicater(vectorstores.NewSimpleDeduplicater()), //  This is useful to prevent wasting time on creating an embedding
	}

	connString := database.MustVectorDatabase(context.Background())

	store, err := pgvector.New(
		context.Background(),
		pgvector.WithConnectionURL(connString),
		pgvector.WithEmbedder(e),
		pgvector.WithCollectionName("documentscollection"),
		pgvector.WithCollectionTableName("documentstable"),
		pgvector.WithEmbeddingTableName("documentsembeddings"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Read the Grafana LGTM document (it was very recently created)
	// so the LLM does not know about it. Therefore, only using RAG
	// will found the answer.

	// get current dir

	bs, err := os.ReadFile(filepath.Join("testdata", "grafana-lgtm.txt"))
	if err != nil {
		log.Fatal(err)
	}

	if database.CheckInitialEmbeddings(connString) {
		_, err = store.AddDocuments(context.Background(), []schema.Document{
			{
				PageContent: string(bs),
			},
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	server.conversationalRetrieval = chains.NewConversationalRetrievalQAFromLLM(
		llm, vectorstores.ToRetriever(store, 10, optionsVector...), memory.NewConversationBuffer())
}
