// api.go
package ollamago

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Generate creates a completion using the specified model
func (c *Client) Generate(ctx context.Context, req GenerateRequest) (*GenerateResponse, error) {
	if req.Model == "" {
		return nil, &RequestError{Message: "model is required"}
	}
	req.Stream = false

	var resp GenerateResponse
	if err := c.request(ctx, http.MethodPost, "/api/generate", req, &resp, false); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GenerateStream creates a streaming completion for the provided prompt
func (c *Client) GenerateStream(ctx context.Context, req GenerateRequest) (<-chan GenerateResponse, <-chan error) {
	responseChan := make(chan GenerateResponse)
	errChan := make(chan error, 1)

	go func() {
		defer close(responseChan)
		defer close(errChan)

		if req.Model == "" {
			errChan <- &RequestError{Message: "model is required"}
			return
		}

		req.Stream = true
		resp, err := c.requestStream(ctx, http.MethodPost, "/api/generate", req)
		if err != nil {
			errChan <- err
			return
		}
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				errChan <- ctx.Err()
				return
			default:
				line := scanner.Bytes()
				if len(line) == 0 {
					continue
				}

				var genResp GenerateResponse
				if err := json.Unmarshal(line, &genResp); err != nil {
					errChan <- fmt.Errorf("failed to decode response: %w", err)
					return
				}

				select {
				case responseChan <- genResp:
				case <-ctx.Done():
					errChan <- ctx.Err()
					return
				}

				if genResp.Done {
					return
				}
			}
		}

		if err := scanner.Err(); err != nil {
			errChan <- fmt.Errorf("error reading response: %w", err)
		}
	}()

	return responseChan, errChan
}

// Chat creates a chat completion using the specified model and messages
func (c *Client) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	if req.Model == "" {
		return nil, &RequestError{Message: "model is required"}
	}
	req.Stream = false
	var resp ChatResponse
	if err := c.request(ctx, http.MethodPost, "/api/chat", req, &resp, false); err != nil {
		return nil, err
	}

	return &resp, nil
}

// ChatStream creates a streaming chat completion
func (c *Client) ChatStream(ctx context.Context, req ChatRequest) (<-chan ChatResponse, <-chan error) {
	respChan := make(chan ChatResponse)
	errChan := make(chan error, 1)

	go func() {
		defer close(respChan)
		defer close(errChan)

		if req.Model == "" {
			errChan <- &RequestError{Message: "model is required"}
			return
		}

		req.Stream = true
		resp, err := c.requestStream(ctx, http.MethodPost, "/api/chat", req)
		if err != nil {
			errChan <- err
			return
		}
		defer resp.Body.Close()

		decoder := json.NewDecoder(resp.Body)
		for {
			var chatResp ChatResponse
			if err := decoder.Decode(&chatResp); err != nil {
				if err == io.EOF {
					return
				}
				errChan <- fmt.Errorf("decode error: %w", err)
				return
			}

			select {
			case respChan <- chatResp:
			case <-ctx.Done():
				errChan <- ctx.Err()
				return
			}

			if chatResp.Done {
				return
			}
		}
	}()

	return respChan, errChan
}

// Embeddings generates embeddings for the provided input
func (c *Client) Embeddings(ctx context.Context, req EmbeddingsRequest) (*EmbeddingsResponse, error) {
	if req.Model == "" {
		return nil, &RequestError{Message: "model is required"}
	}

	var resp EmbeddingsResponse
	if err := c.request(ctx, http.MethodPost, "/api/embeddings", req, &resp, false); err != nil {
		return nil, err
	}

	return &resp, nil
}

// CreateModel creates a model from a Modelfile
func (c *Client) CreateModel(ctx context.Context, req CreateModelRequest) (*ProgressResponse, error) {
	if req.Name == "" {
		return nil, &RequestError{Message: "model name is required"}
	}

	var resp ProgressResponse
	if err := c.request(ctx, http.MethodPost, "/api/create", req, &resp, req.Stream); err != nil {
		return nil, err
	}

	return &resp, nil
}

// ListModels returns a list of local models
func (c *Client) ListModels(ctx context.Context) (*ListModelsResponse, error) {
	var resp ListModelsResponse
	if err := c.request(ctx, http.MethodGet, "/api/tags", nil, &resp, false); err != nil {
		return nil, err
	}

	return &resp, nil
}

// ShowModel shows details about the specified model
func (c *Client) ShowModel(ctx context.Context, req ShowModelRequest) (*ShowModelResponse, error) {
	if req.Name == "" {
		return nil, &RequestError{Message: "model name is required"}
	}

	var resp ShowModelResponse
	if err := c.request(ctx, http.MethodPost, "/api/show", req, &resp, false); err != nil {
		return nil, err
	}

	return &resp, nil
}

// CopyModel creates a copy of a model
func (c *Client) CopyModel(ctx context.Context, req CopyModelRequest) (*StatusResponse, error) {
	if req.Source == "" || req.Destination == "" {
		return nil, &RequestError{Message: "source and destination are required"}
	}

	var resp StatusResponse
	if err := c.request(ctx, http.MethodPost, "/api/copy", req, &resp, false); err != nil {
		return nil, err
	}

	return &resp, nil
}

// DeleteModel removes a model
func (c *Client) DeleteModel(ctx context.Context, req DeleteModelRequest) (*StatusResponse, error) {
	if req.Name == "" {
		return nil, &RequestError{Message: "model name is required"}
	}

	var resp StatusResponse
	if err := c.request(ctx, http.MethodDelete, "/api/delete", req, &resp, false); err != nil {
		return nil, err
	}

	return &resp, nil
}

// PullModel downloads a model from a registry
func (c *Client) PullModel(ctx context.Context, req PullModelRequest) (*ProgressResponse, error) {
	if req.Name == "" {
		return nil, &RequestError{Message: "model name is required"}
	}

	var resp ProgressResponse
	if err := c.request(ctx, http.MethodPost, "/api/pull", req, &resp, req.Stream); err != nil {
		return nil, err
	}

	return &resp, nil
}

// PullModelStream downloads a model with progress updates
func (c *Client) PullModelStream(ctx context.Context, req PullModelRequest) (<-chan ProgressResponse, <-chan error) {
	respChan := make(chan ProgressResponse)
	errChan := make(chan error, 1)

	go func() {
		defer close(respChan)
		defer close(errChan)

		if req.Name == "" {
			errChan <- &RequestError{Message: "model name is required"}
			return
		}

		req.Stream = true
		resp, err := c.requestStream(ctx, http.MethodPost, "/api/pull", req)
		if err != nil {
			errChan <- err
			return
		}
		defer resp.Body.Close()

		decoder := json.NewDecoder(resp.Body)
		for {
			var progressResp ProgressResponse
			if err := decoder.Decode(&progressResp); err != nil {
				if err == io.EOF {
					return
				}
				errChan <- fmt.Errorf("decode error: %w", err)
				return
			}

			select {
			case respChan <- progressResp:
			case <-ctx.Done():
				errChan <- ctx.Err()
				return
			}
		}
	}()

	return respChan, errChan
}

// PushModel uploads a model to a registry
func (c *Client) PushModel(ctx context.Context, req PushModelRequest) (*ProgressResponse, error) {
	if req.Name == "" {
		return nil, &RequestError{Message: "model name is required"}
	}

	var resp ProgressResponse
	if err := c.request(ctx, http.MethodPost, "/api/push", req, &resp, req.Stream); err != nil {
		return nil, err
	}

	return &resp, nil
}

// PushModelStream uploads a model with progress updates
func (c *Client) PushModelStream(ctx context.Context, req PushModelRequest) (<-chan ProgressResponse, <-chan error) {
	respChan := make(chan ProgressResponse)
	errChan := make(chan error, 1)

	go func() {
		defer close(respChan)
		defer close(errChan)

		if req.Name == "" {
			errChan <- &RequestError{Message: "model name is required"}
			return
		}

		req.Stream = true
		resp, err := c.requestStream(ctx, http.MethodPost, "/api/push", req)
		if err != nil {
			errChan <- err
			return
		}
		defer resp.Body.Close()

		decoder := json.NewDecoder(resp.Body)
		for {
			var progressResp ProgressResponse
			if err := decoder.Decode(&progressResp); err != nil {
				if err == io.EOF {
					return
				}
				errChan <- fmt.Errorf("decode error: %w", err)
				return
			}

			select {
			case respChan <- progressResp:
			case <-ctx.Done():
				errChan <- ctx.Err()
				return
			}
		}
	}()

	return respChan, errChan
}
