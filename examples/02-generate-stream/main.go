/*
Topic: Parakeet
Generate a simple completion with Ollama and parakeet
The output is streamed
*/

package main

import (
	"fmt"
	"log"

	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/llm"
	"github.com/parakeet-nest/parakeet/enums/option"

)

func main() {
	ollamaUrl := "http://localhost:11434"
	// if working from a container
	//ollamaUrl := "http://host.docker.internal:11434"
	model := "tinydolphin"

	// Define the options
	//options := llm.DefaultOptions()
	//options.Temperature = 0.5
	// or:

	options := llm.SetOptions(map[string]interface{}{
		option.Temperature: 0.5,
	})

	firstQuestion := llm.GenQuery{
		Model:   model,
		Prompt:  "Who is James T Kirk?",
		Options: options,
	}

	answer, err := completion.GenerateStream(ollamaUrl, firstQuestion,
		func(answer llm.GenAnswer) error {
			fmt.Print(answer.Response)
			return nil
		})

	if err != nil {
		log.Fatal("😡:", err)
	}

	secondQuestion := llm.GenQuery{
		Model:   model,
		Prompt:  "Who is his best friend?",
		Context: answer.Context,
		Options: options,
	}

	fmt.Println()

	_, err = completion.GenerateStream(ollamaUrl, secondQuestion,
		func(answer llm.GenAnswer) error {
			fmt.Print(answer.Response)
			return nil
		})

	if err != nil {
		log.Fatal("😡:", err)
	}
}
