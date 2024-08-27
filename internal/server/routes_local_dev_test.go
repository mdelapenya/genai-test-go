//go:build local_dev
// +build local_dev

package server

import (
	"genai-test-go/internal/ai"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Jeffail/gabs/v2"
)

func TestLLMs(t *testing.T) {
	server.RegisterFiberRoutes()

	testCases := []struct {
		name     string
		basepath string
	}{
		{
			name: "(OpenAI) Direct LLM should not know when the Grafana LGTM module is available",
			// talking to the LLM directly will not provide a good answer, because the model does not have
			// the information about the Grafana module.
			basepath: "/openai/llm",
		},
		{
			name: "(OpenAI) Using RAG must know when the Grafana LGTM module is available",
			// using RAG will provide a better answer because we already added the embeddings for Grafana
			// to the database. See the local_development.go file.
			basepath: "/openai/rag",
		},
		{
			name: "(Ollama) Direct LLM should not know when the Grafana LGTM module is available",
			// talking to the LLM directly will not provide a good answer, because the model does not have
			// the information about the Grafana module.
			basepath: "/ollama/llm",
		},
		{
			name: "(Ollama) Using RAG must know when the Grafana LGTM module is available",
			// using RAG will provide a better answer because we already added the embeddings for Grafana
			// to the database. See the local_development.go file.
			basepath: "/ollama/rag",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			temperature := "0.0"

			pathWithQS := tc.basepath + "?t=" + temperature

			req, err := http.NewRequest("GET", pathWithQS, nil)
			if err != nil {
				t.Fatalf("error creating request. Err: %v", err)
			}

			// use 10 seconds timeout for the request
			resp, err := server.App.Test(req, 10_000)
			if err != nil {
				t.Fatalf("error making request to server. Err: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Fatalf("expected status OK; got %v", resp.Status)
			}

			// read body
			bs, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("error reading response body. Err: %v", err)
			}

			c, err := gabs.ParseJSON(bs)
			if err != nil {
				t.Fatalf("error parsing response body. Err: %v", err)
			}

			type payload struct {
				response  ai.Response
				reference string
				question  string
			}

			answer := c.Path("message").String()
			evaluator := c.Path("evaluator").Data().(map[string]interface{})
			p := payload{
				response: ai.Response{
					Evaluation: evaluator["Evaluation"].(string),
					Reason:     evaluator["Reason"].(string),
				},
				question:  c.Path("question").String(),
				reference: c.Path("reference").String(),
			}

			t.Logf("AI response: %v, answer was %s", evaluator, answer)

			if strings.Contains(tc.basepath, "/llm") {
				if p.response.Evaluation != "no" {
					t.Fatalf("expected the LLM to not know about the Grafana LGTM module: %v", p)
				}
			} else if strings.Contains(tc.basepath, "/rag") {
				if p.response.Evaluation != "yes" {
					t.Fatalf("expected the RAG to know about the Grafana LGTM module: %v", p)
				}
			}
		})
	}
}
