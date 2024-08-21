package main

import (
	"fmt"
	"genai-test-go/internal/server"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	if err := server.Run(); err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
