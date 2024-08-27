package ai

import (
	"context"
	"fmt"

	"github.com/Jeffail/gabs/v2"
	"github.com/tmc/langchaingo/llms"
)

type Evaluator interface {
	Evaluate(ctx context.Context, question string, answer string, reference string) (Response, error)
}

type ModelEvaluator struct {
	model llms.Model
}

type Response struct {
	Answer string "json:answer"
	Reason string "json:reason"
}

// NewEvaluator creates a new Evaluator with the provided model.
func NewEvaluator(m llms.Model) *ModelEvaluator {
	return &ModelEvaluator{
		model: m,
	}
}

// systemPrompt is a string that will be used to provide instructions to the LLM. The most
// specific the instructions, the better the LLM will understand the context of the question,
// for that reason we are providing instructions on how the LLM should respond to the question,
// the format of the response, and several examples for both correct and incorrect answers.
// Each example will include a question, an answer, and a reference to validate the answer.
const sysmtemPrompt string = `
### Instructions
You are a strict validator.
You will be provided with a question, an answer, and a reference.
Your task is to validate whether the answer is correct for the given question, based on the reference.

Follow these instructions:
- Respond only 'yes', 'no' or 'unsure' and always include the reason for your response
- Respond with 'yes' if the answer is correct
- Respond with 'no' if the answer is incorrect
- If you are unsure, simply respond with 'unsure'
- Respond with 'no' if the answer is not clear or concise
- Respond with 'no' if the answer is not based on the reference
- Do not include the examples in your response

Your response must be a json object with the following structure:
{
	"response": "yes",
	"reason": "The answer is correct because it is based on the reference provided."
}

### Example 1
Question: Is Madrid the capital of Spain?
Answer: No, it's Barcelona.
Reference: The capital of Spain is Madrid
###
Response: {
	"response": "no",
	"reason": "The answer is incorrect because the reference states that the capital of Spain is Madrid."
}

### Example 2
Question: What is the planet closest to the sun?
Answer: Mercury
Reference: The planet closest to the sun is Mercury
###
Response: {
	"response": "yes",
	"reason": "The answer is correct because it is based on the reference provided."
}
`

// userPrompt is the prompt used by Users talking to the LLM.
// It's a template that is used to provide the question, answer, and reference to the LLM.
// This structured format helps the LLM to understand the context of the question, answer and reference.
const userPrompt string = `
###
Question: %s
###
Answer: %s
###
Reference: %s
###
`

func (e *ModelEvaluator) Evaluate(ctx context.Context, question string, answer string, reference string) (Response, error) {
	res, err := e.model.GenerateContent(ctx, []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, sysmtemPrompt),
		llms.TextParts(llms.ChatMessageTypeHuman, fmt.Sprintf(userPrompt, question, answer, reference)),
	},
		llms.WithModel("gpt-4"),
		llms.WithTemperature(0), // deterministic responses
		llms.WithSeed(42),
		llms.WithTopP(0), // top-p zero means that the model will always provide the most likely answer
	)
	if err != nil {
		return Response{}, err
	}

	var r Response
	c, err := gabs.ParseJSON([]byte(res.Choices[0].Content))
	if err != nil {
		return r, err
	}

	r = Response{
		Answer: c.Path("response").Data().(string),
		Reason: c.Path("reason").Data().(string),
	}

	return r, nil
}
