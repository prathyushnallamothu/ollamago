# OllamaGo

A powerful and feature-rich Go client for [Ollama](https://ollama.com/), enabling seamless integration of Ollama's capabilities into your Go applications.

## Features

- 🚀 Complete API coverage for Ollama
- 💬 Text generation and chat completions
- 🔄 Streaming responses support
- 🛠 Model management (create, pull, push, copy, delete)
- 📊 Embeddings generation
- 🔧 Customizable options and parameters
- 🎯 Type-safe requests and responses

## Installation

```bash
go get github.com/prathyushnallamothu/ollamago
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    ollama "github.com/prathyushnallamothu/ollamago"
)

func main() {
    // Create a new client with custom timeout
    client := ollama.NewClient(
        ollama.WithTimeout(time.Minute*5),
    )

    // Generate text
    resp, err := client.Generate(context.Background(), &ollama.GenerateRequest{
        Model:  "llama2",
        Prompt: "What is the capital of France?",
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(resp.Response)
}
```

## Features in Detail

### Chat Completions

```go
messages := []ollama.Message{
    {
        Role:    "system",
        Content: "You are a helpful assistant.",
    },
    {
        Role:    "user",
        Content: "What's the weather like today?",
    },
}

resp, err := client.Chat(context.Background(), &ollama.ChatRequest{
    Model:    "llama2",
    Messages: messages,
})
```

### Streaming Responses

```go
req := &ollama.GenerateRequest{
    Model:   "llama2",
    Prompt:  "Write a story about a space adventure",
    Stream:  true,
}

err := client.GenerateStream(context.Background(), req, func(response *ollama.GenerateResponse) error {
    fmt.Print(response.Response)
    return nil
})
```

### Model Management

```go
// List available models
models, err := client.ListModels(context.Background())

// Pull a model
err = client.PullModel(context.Background(), &ollama.PullModelRequest{
    Name: "llama2",
})

// Delete a model
err = client.DeleteModel(context.Background(), &ollama.DeleteModelRequest{
    Name: "llama2",
})
```

### Embeddings

```go
resp, err := client.CreateEmbedding(context.Background(), &ollama.EmbedRequest{
    Model:  "llama2",
    Prompt: "Hello, world!",
})
```

## Configuration Options

The client can be configured with various options:

```go
client := ollama.NewClient(
    ollama.WithTimeout(time.Minute*5),
    ollama.WithBaseURL("http://localhost:11434"),
    ollama.WithHeader("Custom-Header", "value"),
)
```

## Model Parameters

Fine-tune model behavior with various parameters:

```go
options := &ollama.Options{
    Temperature:     ollama.Float64Ptr(0.7),
    TopK:           ollama.IntPtr(40),
    TopP:           ollama.Float64Ptr(0.9),
    PresencePenalty: ollama.Float64Ptr(1.0),
    FrequencyPenalty: ollama.Float64Ptr(1.0),
}

req := &ollama.GenerateRequest{
    Model:   "llama2",
    Prompt:  "Write a poem",
    Options: options,
}
```

## Error Handling

The package provides structured error types for better error handling:

- `RequestError`: Client-side request errors
- `ResponseError`: Server-side API response errors

```go
if err != nil {
    switch e := err.(type) {
    case *ollama.RequestError:
        log.Printf("Request error: %v", e.Message)
    case *ollama.ResponseError:
        log.Printf("API error %d: %v", e.StatusCode, e.Message)
    default:
        log.Printf("Unknown error: %v", err)
    }
}
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
