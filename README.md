# GenAI Testing in Go

In this demo project, we'd like the different models to be able to determine in which Testcontainers for Go version the Grafana module was added. The models should be able to answer questions like:

- In which version of Testcontainers for Go was the Grafana module added?

For that, we will use two approaches: using an LLM without `RAG` and using the LLM with `RAG`, creating the embeddings for the valid response, to demonstrate that using the `RAG` model can provide more accurate results. Finally, we will use the `GPT-4` model as an evaluator to validate the different responses.

## Challenges of Testing GenAI applications

The real challenge arises when trying to test the responses generated by language models. Traditionally, we could settle for verifying that the response includes certain keywords, which is insufficient and prone to errors.

This approach is not only fragile but also lacks the ability to assess the relevance or coherence of the response.

An alternative is to employ cosine similarity to compare the embeddings of a "_reference_" response and the actual response, providing a more semantic form of evaluation.

This method measures the similarity between two vectors/embeddings by calculating the cosine of the angle between them. If both vectors point to the same direction, it means the "_reference_" response is semantically the same as the actual response.

However, this method introduces the problem of selecting an appropriate threshold to determine the acceptability of the response, in addition to the opacity of the evaluation process.

## Towards a more effective method

The real problem here arises from the fact that answers provided by the LLM are in natural language and non-deterministic.
Because of this, using current testing methods to verify them is difficult, as these methods are better suited to testing predictable values. 

However, we already have a great tool for understanding non-deterministic answers in natural language: LLMs themselves.
Thus, the key may lie in using one LLM to evaluate the adequacy of responses generated by another LLM. 

This proposal involves defining detailed validation criteria and using an LLM as an **Evaluator** to determine if the responses meet the specified requirements. This approach can be applied to validate answers to specific questions, drawing on both general knowledge and specialized information.

By incorporating detailed instructions and examples, the Evaluator can provide accurate and justified evaluations, offering clarity on why a response is considered correct or incorrect. For this reason, the evaluator model should be the most accurate and reliable model available.
In this project, we are using `OpenAI's GPT-4` as the evaluator model.

Please check the bootstrap of the [local development mode file](./internal/server/local_development.go) to understand how the evaluator model is loaded and used. We will start and a `pgVector` instance for RAG, pre-loaded with information about the recently created Grafana module for Testcontainers for Go.

## Local Development

First, create a `.env` file in the root of the project with the contents of the [`.env.example`](.env.example) file. This file contains the environment variables required by the application.

The project uses [air](https://github.com/air-verse/air) for live reloading the application. To start the application with live reload, run the following command:

```bash
make watch
```

This will start the application and watch for changes in the source code. When a change is detected, the application will be recompiled and restarted.

The project uses [Testcontainers for Go](https://github.com/testcontainers/testcontainers-go) for running the local development environment for the application. It is formed by the following services:

- a local Postgres database, with the `pgVector` module enabled, as a Docker container.
- a local Ollama instance, with the `llama3.1:8b` model already loaded in it.

These containers are reused across multiple builds, so you don't end up with multiple containers running at the same time.

Please check the [`server/local_development.go`](./internal/server/local_development.go) file. This file leverages Go build tags and a Go's init function to conditionally start the Postgres container only during local development. The `make build-dev` command adds the proper build tags to the execution of the Go toolchain in order to make it possible. To understand how this works, check [the following blog post](https://www.docker.com/blog/local-development-of-go-applications-with-testcontainers/).

In a nutshell, you start the application using `air`, that watches for changes; and when a change is detected, the application is recompiled with the `make build-dev` command, so the local development environment is included in the build process. The database container will be started only once, and it will be reused across multiple builds.

## Running the tests

To run the integrations tests, which automatically starts the local development environment, execute the following command:

```bash
make itest
```
