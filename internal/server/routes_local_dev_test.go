//go:build local_dev
// +build local_dev

package server

import (
	"context"
	"genai-test-go/internal/ai"
	"io"
	"net/http"
	"testing"

	"github.com/Jeffail/gabs/v2"
)

// Using must/should is important
const reference = `- Answer must not mention any other module
- Answer must mention the version of Testcontainers for Go, which is v0.33.0
- Answer must be less than 5 sentences`

func TestLLMs(t *testing.T) {
	server.RegisterFiberRoutes()

	testCases := []struct {
		name     string
		basepath string
	}{
		{
			name: "Direct LLM should not know when the Grafana LGTM module is available",
			// talking to the LLM directly will not provide a good answer, because the model does not have
			// the information about the Grafana module.
			basepath: "/chat/llm",
		},
		{
			name: "Using RAG must know when the Grafana LGTM module is available",
			// using RAG will provide a better answer because we already added the embeddings for Grafana
			// to the database. See the local_development.go file.
			basepath: "/chat/rag",
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

			answer := c.Path("message").String()

			evaluator := ai.NewEvaluator(server.llm)

			aiResp, err := evaluator.Evaluate(context.Background(), question, answer, reference)
			if err != nil {
				t.Fatalf("error evaluating response. Err: %v", err)
			}

			t.Logf("AI response: %v, answer was %s", aiResp, answer)

			if tc.basepath == "/chat/llm" {
				if aiResp.Answer != "no" {
					t.Fatalf("expected the LLM to not know what TTV means: %v", aiResp)
				}
			} else if tc.basepath == "/chat/rag" {
				if aiResp.Answer != "yes" {
					t.Fatalf("expected the RAG to know what TTV means: %v", aiResp)
				}
			}
		})
	}
}
