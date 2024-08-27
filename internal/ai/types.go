package ai

import (
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
)

type Chat struct {
	// model to be used
	Model llms.Model

	// using RAG for conversational retrieval
	ConversationalRetrieval chains.ConversationalRetrievalQA
}

func NewChat(model llms.Model, conversationalRetrieval chains.ConversationalRetrievalQA) *Chat {
	return &Chat{
		Model:                   model,
		ConversationalRetrieval: conversationalRetrieval,
	}
}
