/*
Topic: Parakeet
Generate a chat completion using function calling
no streaming
This example:
- uses Mistral model
- make a list of tools
- use the list of tools to generate content for the prompt
- retrieve the function from the list of tools
*/

package main

import (
	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/llm"

	"fmt"
	"log"
)

func main() {
	ollamaUrl := "http://localhost:11434"
	// if working from a container
	//ollamaUrl := "http://host.docker.internal:11434"
	model := "sam4096/qwen2tools:0.5b"

	toolsList := []llm.Tool{
		{
			Type: "function",
			Function: llm.Function{
				Name:        "hello",
				Description: "Say hello to a given person with his name",
				Parameters: llm.Parameters{
					Type: "object",
					Properties: map[string]llm.Property{
						"name": {
							Type:        "string",
							Description: "The name of the person",
						},
					},
					Required: []string{"name"},
				},
			},
		},
		{
			Type: "function",
			Function: llm.Function{
				Name:        "addNumbers",
				Description: "Make an addition of the two given numbers",
				Parameters: llm.Parameters{
					Type: "object",
					Properties: map[string]llm.Property{
						"a": {
							Type:        "number",
							Description: "first operand",
						},
						"b": {
							Type:        "number",
							Description: "second operand",
						},
					},
					Required: []string{"a", "b"},
				},
			},
		},
	}

	messages := []llm.Message{
		{Role: "user", Content: `say "hello" to Bob`},
	}

	options := llm.Options{
		Temperature:   0.0,
		RepeatLastN:   2,
		RepeatPenalty: 2.0,
	}

	query := llm.Query{
		Model:    model,
		Messages: messages,
		Tools:    toolsList,
		Options:  options,
		Format:   "json",
	}

	answer, err := completion.Chat(ollamaUrl, query)
	if err != nil {
		log.Fatal("😡:", err)
	}


	result, err := answer.Message.ToolCalls[0].Function.ToJSONString()
	if err != nil {
		log.Fatal("😡:", err)
	}
	fmt.Println(result)

	messages = []llm.Message{
		{Role: "user", Content: `add 2 and 40`},
	}

	query = llm.Query{
		Model:    model,
		Messages: messages,
		Tools:    toolsList,
		Options:  options,
		Format:   "json",
	}

	answer, err = completion.Chat(ollamaUrl, query)
	if err != nil {
		log.Fatal("😡:", err)
	}

	result, err = answer.Message.ToolCalls[0].Function.ToJSONString()
	if err != nil {
		log.Fatal("😡:", err)
	}
	fmt.Println(result)

}
