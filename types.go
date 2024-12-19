// types.go
package ollamago

import (
	"encoding/json"
	"fmt"
	"time"
)

// Version represents the current version of the client
const Version = "0.1.0"

// Options represents model parameters and inference options
type Options struct {
	NumKeep          *int     `json:"num_keep,omitempty"`
	Seed            *int     `json:"seed,omitempty"`
	NumPredict      *int     `json:"num_predict,omitempty"`
	TopK            *int     `json:"top_k,omitempty"`
	TopP            *float64 `json:"top_p,omitempty"`
	TFSZ            *float64 `json:"tfs_z,omitempty"`
	TypicalP        *float64 `json:"typical_p,omitempty"`
	RepeatLastN     *int     `json:"repeat_last_n,omitempty"`
	Temperature     *float64 `json:"temperature,omitempty"`
	RepeatPenalty   *float64 `json:"repeat_penalty,omitempty"`
	PresencePenalty *float64 `json:"presence_penalty,omitempty"`
	FrequencyPenalty *float64 `json:"frequency_penalty,omitempty"`
	Mirostat        *int     `json:"mirostat,omitempty"`
	MirostatTau     *float64 `json:"mirostat_tau,omitempty"`
	MirostatEta     *float64 `json:"mirostat_eta,omitempty"`
	PenalizeNewline *bool    `json:"penalize_newline,omitempty"`
	Stop            []string `json:"stop,omitempty"`
	NumGPU          *int     `json:"num_gpu,omitempty"`
	NumThread       *int     `json:"num_thread,omitempty"`
	NumCtx          *int     `json:"num_ctx,omitempty"`
	LogitsAll       *bool    `json:"logits_all,omitempty"`
	EmbeddingOnly   *bool    `json:"embedding_only,omitempty"`
	F16KV           *bool    `json:"f16_kv,omitempty"`
}

// Message represents a chat message
type Message struct {
	Role      string     `json:"role"`
	Content   string     `json:"content,omitempty"`
	Images    []Image    `json:"images,omitempty"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
	Name      string     `json:"name,omitempty"`
}

// Image represents an image for multimodal models
type Image struct {
	Data string `json:"data"`
}

// Function represents a function definition
type Function struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Parameters  json.RawMessage `json:"parameters,omitempty"`
}

// ToolCall represents a function call from the model
type ToolCall struct {
	ID       string          `json:"id"`
	Type     string          `json:"type"`
	Function FunctionCall    `json:"function"`
}

// FunctionCall represents the details of a function call
type FunctionCall struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

// Tool represents a tool available to the model
type Tool struct {
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

// GenerateRequest represents a completion request
type GenerateRequest struct {
	Model     string   `json:"model"`
	Prompt    string   `json:"prompt,omitempty"`
	System    string   `json:"system,omitempty"`
	Template  string   `json:"template,omitempty"`
	Context   []int    `json:"context,omitempty"`
	Stream    bool     `json:"stream"`
	Raw       bool     `json:"raw,omitempty"`
	Format    string   `json:"format,omitempty"`
	Images    []Image  `json:"images,omitempty"`
	Options   *Options `json:"options,omitempty"`
	KeepAlive string   `json:"keep_alive,omitempty"`
}

// GenerateResponse represents a completion response
type GenerateResponse struct {
	Model             string  `json:"model,omitempty"`
	CreatedAt        string  `json:"created_at,omitempty"`
	Response         string  `json:"response"`
	Done             bool    `json:"done,omitempty"`
	Context          []int   `json:"context,omitempty"`
	TotalDuration    int64   `json:"total_duration,omitempty"`
	LoadDuration     int64   `json:"load_duration,omitempty"`
	PromptEvalCount  int     `json:"prompt_eval_count,omitempty"`
	EvalCount        int     `json:"eval_count,omitempty"`
	EvalDuration     int64   `json:"eval_duration,omitempty"`
}

// ChatRequest represents a chat completion request
type ChatRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	Format    string    `json:"format,omitempty"`
	Stream    bool      `json:"stream"`
	Tools     []Tool    `json:"tools,omitempty"`
	Options   *Options  `json:"options,omitempty"`
	KeepAlive string    `json:"keep_alive,omitempty"`
}

// ChatResponse represents a chat completion response
type ChatResponse struct {
	Model            string   `json:"model,omitempty"`
	CreatedAt        string   `json:"created_at,omitempty"`
	Message          Message  `json:"message"`
	Done             bool     `json:"done,omitempty"`
	TotalDuration    int64    `json:"total_duration,omitempty"`
	LoadDuration     int64    `json:"load_duration,omitempty"`
	PromptEvalCount  int      `json:"prompt_eval_count,omitempty"`
	EvalCount        int      `json:"eval_count,omitempty"`
	EvalDuration     int64    `json:"eval_duration,omitempty"`
}

// EmbedRequest represents an embedding request
type EmbedRequest struct {
	Model     string   `json:"model"`
	Prompt    string   `json:"prompt,omitempty"`
	Options   *Options `json:"options,omitempty"`
	KeepAlive string   `json:"keep_alive,omitempty"`
}

// EmbedResponse represents an embedding response
type EmbedResponse struct {
	Embeddings []float64 `json:"embedding"`
}

// CreateRequest represents a model creation request
type CreateRequest struct {
	Name      string `json:"name"`
	Path      string `json:"-"` // Local file path, not sent to API
	Modelfile string `json:"modelfile"`
	Stream    bool   `json:"stream,omitempty"`
}

// PullRequest represents a model download request
type PullRequest struct {
	Name     string `json:"name"`
	Insecure bool   `json:"insecure,omitempty"`
	Stream   bool   `json:"stream,omitempty"`
}

// PushRequest represents a model upload request
type PushRequest struct {
	Name     string `json:"name"`
	Insecure bool   `json:"insecure,omitempty"`
	Stream   bool   `json:"stream,omitempty"`
}

// CopyRequest represents a model copy request
type CopyRequest struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

// DeleteRequest represents a model deletion request
type DeleteRequest struct {
	Name string `json:"name"`
}

// ShowModelRequest represents a request to show model details
type ShowModelRequest struct {
    Name string `json:"model"`
}

// ShowModelResponse represents detailed information about a model
type ShowModelResponse struct {
    ModelFile  string                 `json:"modelfile,omitempty"`
    Template   string                 `json:"template,omitempty"`
    Parameters string                 `json:"parameters,omitempty"`
    License    string                 `json:"license,omitempty"`
    Details    ModelDetails           `json:"details,omitempty"`
    ModelInfo  map[string]interface{} `json:"model_info,omitempty"`
    ModifiedAt time.Time              `json:"modified_at,omitempty"`
}

// CopyModelRequest represents a request to copy a model
type CopyModelRequest struct {
    Source      string `json:"source"`
    Destination string `json:"destination"`
}

// DeleteModelRequest represents a request to delete a model
type DeleteModelRequest struct {
    Name string `json:"model"`
}

// PullModelRequest represents a request to pull a model from a registry
type PullModelRequest struct {
    Name     string `json:"model"`
    Insecure bool   `json:"insecure,omitempty"`
    Stream   bool   `json:"stream,omitempty"`
}

// PushModelRequest represents a request to push a model to a registry
type PushModelRequest struct {
    Name     string `json:"model"`
    Insecure bool   `json:"insecure,omitempty"`
    Stream   bool   `json:"stream,omitempty"`
}

// EmbeddingsRequest represents a request to generate embeddings
type EmbeddingsRequest struct {
    Model     string    `json:"model"`
    Prompt    string    `json:"prompt"`
    Options   *Options  `json:"options,omitempty"`
    KeepAlive string    `json:"keep_alive,omitempty"`
}

// EmbeddingsResponse represents the response containing embeddings
type EmbeddingsResponse struct {
    Embedding []float64 `json:"embedding"`
}

// CreateModelRequest represents a request to create a new model
type CreateModelRequest struct {
    Model     string `json:"model"`
    Path      string `json:"-"` // used locally, not sent to API
    Modelfile string `json:"modelfile"`
    Stream    bool   `json:"stream,omitempty"`
	Name      string `json:"name"`
}

// ListModelsResponse represents the response containing available models
type ListModelsResponse struct {
    Models []ModelInfo `json:"models"`
}

// ModelInfo represents information about a model
type ModelInfo struct {
    Name       string       `json:"name"`
    ModifiedAt time.Time    `json:"modified_at"`
    Digest     string       `json:"digest,omitempty"`
    Size       int64        `json:"size"`
    Details    ModelDetails `json:"details,omitempty"`
}

// ListResponse represents a model list response
type ListResponse struct {
	Models []Model `json:"models"`
}

// Model represents model information
type Model struct {
	Name       string       `json:"name"`
	Size       int64        `json:"size"`
	ModifiedAt time.Time    `json:"modified_at"`
	Digest     string       `json:"digest,omitempty"`
	Details    ModelDetails `json:"details,omitempty"`
}

// ModelDetails represents detailed model information
type ModelDetails struct {
	Format           string   `json:"format,omitempty"`
	Family           string   `json:"family,omitempty"`
	Families         []string `json:"families,omitempty"`
	ParameterSize    string   `json:"parameter_size,omitempty"`
	QuantizationLevel string   `json:"quantization_level,omitempty"`
}

// ShowResponse represents detailed model information
type ShowResponse struct {
	License    string                 `json:"license,omitempty"`
	Modelfile  string                 `json:"modelfile,omitempty"`
	Template   string                 `json:"template,omitempty"`
	System     string                 `json:"system,omitempty"`
	Parameters string                 `json:"parameters,omitempty"`
	Details    ModelDetails           `json:"details,omitempty"`
}

// StatusResponse represents a basic status response
type StatusResponse struct {
	Status string `json:"status"`
}

// ProgressResponse represents a progress status response
type ProgressResponse struct {
	Status    string `json:"status"`
	Digest    string `json:"digest,omitempty"`
	Total     int64  `json:"total,omitempty"`
	Completed int64  `json:"completed,omitempty"`
	Error     string `json:"error,omitempty"`
}

// RequestError represents a client request error
type RequestError struct {
	Message string
}

func (e *RequestError) Error() string {
	return e.Message
}

// ResponseError represents an API response error
type ResponseError struct {
	StatusCode int
	Message    string
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("status %d: %s", e.StatusCode, e.Message)
}