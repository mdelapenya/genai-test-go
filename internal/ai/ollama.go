package ai

import (
	"context"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/ollama"
)

const ollamaImg string = "ilopezluna/ollama-llama3.1:0.3.6-8b"
const Model string = "llama3.1:8b"

func MustGetOllamaConnectionString(ctx context.Context) string {
	ollamaContainer, err := ollama.Run(ctx,
		ollamaImg,
		testcontainers.CustomizeRequest(testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				Name: "ollama-ctr",
			},
			Reuse: true,
		}),
	)
	if err != nil {
		panic(err)
	}

	connectionStr, err := ollamaContainer.ConnectionString(ctx)

	if err != nil {
		panic(err)
	}

	return connectionStr
}
