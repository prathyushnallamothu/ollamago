package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	ollama "github.com/prathyushnallamothu/ollamago"
)

func main() {
	// Create a new client with custom timeout
	client := ollama.NewClient(
		ollama.WithTimeout(time.Minute*5),
		ollama.WithHeader("Custom-Header", "value"),
	)

	// Example 1: Basic Generation
	fmt.Println("\n=== Basic Generation ===")
	if err := basicGeneration(client); err != nil {
		log.Fatal("Basic generation failed:", err)
	}

	// Example 2: Streaming Generation
	fmt.Println("\n=== Streaming Generation ===")
	if err := streamingGeneration(client); err != nil {
		log.Fatal("Streaming generation failed:", err)
	}

	// Example 3: Chat
	fmt.Println("\n=== Chat Example ===")
	if err := chatExample(client); err != nil {
		log.Fatal("Chat example failed:", err)
	}

	// Example 4: Model Management
	fmt.Println("\n=== Model Management ===")
	if err := modelManagement(client); err != nil {
		log.Fatal("Model management failed:", err)
	}
}

func basicGeneration(client *ollama.Client) error {
	ctx := context.Background()
	//topK := int(40)
	resp, err := client.Generate(ctx, ollama.GenerateRequest{
		Model:  "llama3.2:latest",
		Prompt: "What is the capital of France?",
		Stream: false,
	})
	if err != nil {
		return fmt.Errorf("generate failed: %w", err)
	}

	fmt.Printf("Response: %s\n", resp.Response)
	return nil
}

func streamingGeneration(client *ollama.Client) error {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	numPredict := int(500)
	respChan, errChan := client.GenerateStream(ctx, ollama.GenerateRequest{
		Model:  "llama3.2:latest",
		Prompt: "Write a short story about a brave knight.",
		Options: &ollama.Options{
			Temperature: float64Ptr(0.9),
			NumPredict:  &numPredict,      // Set maximum token length
			Stop:        []string{"\n\n"}, // Stop on double newline
		},
	})

	var fullResponse strings.Builder

	// Handle streaming responses
	for {
		select {
		case resp, ok := <-respChan:
			if !ok {
				if fullResponse.Len() > 0 {
					return nil // Normal completion
				}
				return fmt.Errorf("response channel closed without data")
			}
			fullResponse.WriteString(resp.Response)
			fmt.Print(resp.Response)
			if resp.Done {
				fmt.Println("\n\nGeneration complete!")
				return nil
			}
		case err := <-errChan:
			if err != nil {
				return fmt.Errorf("streaming error: %w", err)
			}
		case <-ctx.Done():
			return fmt.Errorf("generation timed out: %w", ctx.Err())
		}
	}
}

func chatExample(client *ollama.Client) error {
	ctx := context.Background()

	messages := []ollama.Message{
		{
			Role:    "system",
			Content: "You are a helpful AI assistant.",
		},
		{
			Role:    "user",
			Content: "Hello! Can you help me with some Go programming?",
		},
	}

	resp, err := client.Chat(ctx, ollama.ChatRequest{
		Model:    "llama3.2:latest",
		Messages: messages,
		Options: &ollama.Options{
			Temperature: float64Ptr(0.7),
		},
	})
	if err != nil {
		return fmt.Errorf("chat failed: %w", err)
	}

	fmt.Printf("Assistant: %s\n", resp.Message.Content)
	return nil
}

func modelManagement(client *ollama.Client) error {
	ctx := context.Background()

	// List available models
	models, err := client.ListModels(ctx)
	if err != nil {
		return fmt.Errorf("listing models failed: %w", err)
	}

	fmt.Println("Available models:")
	for _, model := range models.Models {
		fmt.Printf("- %s (modified: %s)\n", model.Name, model.ModifiedAt.Format(time.RFC3339))
	}

	// Show details of a specific model
	modelInfo, err := client.ShowModel(ctx, ollama.ShowModelRequest{
		Name: "llama3.2:latest",
	})
	if err != nil {
		if _, ok := err.(*ollama.ResponseError); ok {
			fmt.Println("Model not found or not available")
			return nil
		}
		return fmt.Errorf("showing model failed: %w", err)
	}

	fmt.Printf("\nModel details for llama3.2:latest:\n")
	fmt.Printf("License: %s\n", modelInfo.License)
	fmt.Printf("Parameters: %s\n", modelInfo.Parameters)
	if modelInfo.Details.Family != "" {
		fmt.Printf("Family: %s\n", modelInfo.Details.Family)
	}

	return nil
}

// Helper functions for Options
func float64Ptr(v float64) *float64 {
	return &v
}

func int64Ptr(v int64) *int64 {
	return &v
}
