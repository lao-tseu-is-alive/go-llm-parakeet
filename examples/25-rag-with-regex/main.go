package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/content"
	"github.com/parakeet-nest/parakeet/embeddings"
	"github.com/parakeet-nest/parakeet/llm"

	"github.com/parakeet-nest/parakeet/enums/option"

)

func main() {
	ollamaUrl := "http://localhost:11434"
	var embeddingsModel = "all-minilm:33m" // This model is for the embeddings of the documents
	var smallChatModel = "qwen2:0.5b" // This model is for the chat completion

	store := createEmbeddingsFromDocument("./chronicles.md", ollamaUrl, embeddingsModel)

	getCompletion(`Who are the monsters of Chronicles of Aethelgard?`,
		ollamaUrl, embeddingsModel, smallChatModel, store, 0.3)

	/*
	getCompletion(`Tell me more about Keegorg`,
		ollamaUrl, embeddingsModel, smallChatModel, store, 0.3)
	*/
}

func createEmbeddingsFromDocument(pathDocument, ollamaUrl, embeddingsModel string) embeddings.MemoryVectorStore {

	store := embeddings.MemoryVectorStore{
		Records: make(map[string]llm.VectorRecord),
	}

	rulesContent, err := content.ReadTextFile(pathDocument)
	if err != nil {
		log.Fatalln("😡:", err)
	}
	//chunks := content.SplitTextWithDelimiter(rulesContent, "<!-- SPLIT -->")
	//chunks := content.SplitText(rulesContent, "##")
	chunks := content.SplitTextWithRegex(rulesContent, `## *`)

	// Create embeddings from documents and save them in the store
	for idx, doc := range chunks {
		fmt.Println("Creating embedding from document ", idx)
		embedding, err := embeddings.CreateEmbedding(
			ollamaUrl,
			llm.Query4Embedding{
				Model:  embeddingsModel,
				Prompt: doc,
			},
			strconv.Itoa(idx),
		)
		if err != nil {
			fmt.Println("😡:", err)
		} else {
			embedding.SimpleMetaData = "📝 chunk num: " + strconv.Itoa(idx)
			store.Save(embedding)
		}
	}

	return store
}

func getContentFromSimilarities(userContent, ollamaUrl, embeddingsModel string, store embeddings.MemoryVectorStore, limit float64) string {

	// Create an embedding from the question
	embeddingFromQuestion, err := embeddings.CreateEmbedding(
		ollamaUrl,
		llm.Query4Embedding{
			Model:  embeddingsModel,
			Prompt: userContent,
		},
		"question",
	)
	if err != nil {
		log.Fatalln("😡:", err)
	}
	fmt.Println("🔎 searching for similarity...")

	similarities, _ := store.SearchSimilarities(embeddingFromQuestion, limit)

	for _, similarity := range similarities {
		// Do something with the similarity
		fmt.Println("Similarity:", similarity.SimpleMetaData)
	}

	fmt.Println("🎉 number of similarities:", len(similarities))

	documentsContent := embeddings.GenerateContentFromSimilarities(similarities)

	return documentsContent

}

func getCompletion(userContent, ollamaUrl, embeddingsModel, smallChatModel string, store embeddings.MemoryVectorStore, limit float64) {

	systemContent := `You are the dungeon master,
	expert at interpreting and answering questions based on provided sources.
	Using only the provided context, answer the user's question 
	to the best of your ability using only the resources provided. 
	Be verbose!`

	documentsContent := getContentFromSimilarities(userContent, ollamaUrl, embeddingsModel, store, 0.3)

	options := llm.SetOptions(map[string]interface{}{
		option.Temperature: 0.0,
		option.RepeatLastN: 2,
		option.RepeatPenalty: 3.0,
		option.TopK: 10,
		option.TopP: 0.5,
	})

	query := llm.Query{
		Model: smallChatModel,
		Messages: []llm.Message{
			{Role: "system", Content: systemContent},
			{Role: "system", Content: documentsContent},
			{Role: "user", Content: userContent},
		},
		Options: options,
	}

	fmt.Println()
	fmt.Println("🤖 answer:")

	// Answer the question
	_, err := completion.ChatStream(ollamaUrl, query,
		func(answer llm.Answer) error {
			fmt.Print(answer.Message.Content)
			return nil
		})

	if err != nil {
		log.Fatal("😡:", err)
	}

	fmt.Println()

}
