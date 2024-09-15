package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/embeddings"
	"github.com/parakeet-nest/parakeet/llm"
	"github.com/parakeet-nest/parakeet/enums/option"

)

func main() {
	ollamaUrl := "http://localhost:11434"

	//smallChatModel := "qwen2:0.5b"
	smallChatModel := "gemma2:2b"
	embeddingsModel := "all-minilm:33m"


	systemContent := `instruction: 
	translate the user question in docker command using the given context.
	Stay brief.`

	store := embeddings.BboltVectorStore{}
	store.Initialize("../embeddings.db")

	options := llm.SetOptions(map[string]interface{}{
		option.Temperature: 0.0,
		option.RepeatLastN: 2,
		option.RepeatPenalty: 3.0,
		option.TopK: 10,
		option.TopP: 0.5,
	})

	for {
		question := input(smallChatModel)
		if question == "bye" {
			break
		}

		// Create an embedding from the question
		embeddingFromQuestion, err := embeddings.CreateEmbedding(
			ollamaUrl,
			llm.Query4Embedding{
				Model:  embeddingsModel,
				Prompt: question,
			},
			"question",
		)
		if err != nil {
			log.Fatalln("😡:", err)
		}
		fmt.Println("🔎 searching for similarity...")
		similarities, _ := store.SearchTopNSimilarities(embeddingFromQuestion, 0.4, 3)

		contextContent := embeddings.GenerateContextFromSimilarities(similarities)
		//fmt.Println(documentsContent)
		fmt.Println("🎉 similarities:", len(similarities))

		// Prepare the query
		query := llm.Query{
			Model: smallChatModel,
			Messages: []llm.Message{
				{Role: "system", Content: systemContent},
				{Role: "system", Content: contextContent},
				{Role: "user", Content: question},
			},
			Options: options,
		}

		// Answer the question
		_, err = completion.ChatStream(ollamaUrl, query,
			func(answer llm.Answer) error {
				fmt.Print(answer.Message.Content)
				return nil
			})

		if err != nil {
			log.Fatal("😡:", err)
		}

		fmt.Println()

	}
}

func input(smallChatModel string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("🐳 [%s] ask me something> ", smallChatModel)
	question, _ := reader.ReadString('\n')
	return strings.TrimSpace(question)
}
