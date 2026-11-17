package ai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/nagi-hikaru/gofuck/internal/config"
	"github.com/nagi-hikaru/gofuck/internal/status"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

const (
	maxAIResponse = 64 * 1024
)

// Suggest 向指定模型发送原命令和 stderr，并返回解析后的候选。
func Suggest(ctx context.Context, cfg config.ProviderConfig, modelName string, input Input) ([]Candidate, error) {
	client := &http.Client{Timeout: 20 * time.Second}
	completions := openai.NewChatCompletionService(
		option.WithBaseURL(cfg.BaseURL),
		option.WithAPIKey(cfg.APIKey),
		option.WithHTTPClient(client),
		option.WithMaxRetries(0),
		option.WithMiddleware(limitModelResponse),
	)

	body, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	reply, err := completions.New(ctx, openai.ChatCompletionNewParams{
		Model: modelName,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt),
			openai.UserMessage(string(body)),
		},
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONObject: &openai.ResponseFormatJSONObjectParam{},
		},
		MaxTokens: openai.Int(512),
	}, option.WithJSONSet("thinking.type", "disabled"))

	if err != nil {
		var syntaxError *json.SyntaxError
		var typeError *json.UnmarshalTypeError
		if errors.Is(err, io.EOF) || errors.As(err, &syntaxError) || errors.As(err, &typeError) {
			return nil, status.ErrInvalidModelResponse
		}
		return nil, fmt.Errorf("request model: %w", err)
	}
	if reply == nil || len(reply.Choices) == 0 {
		return nil, status.ErrInvalidModelResponse
	}
	choice := reply.Choices[0]
	if choice.FinishReason != "stop" {
		return nil, status.ErrIncompleteModelResponse
	}

	var result Result
	if err := json.Unmarshal([]byte(choice.Message.Content), &result); err != nil || result.Candidates == nil {
		return nil, status.ErrInvalidModelResponse
	}
	return result.Candidates, nil
}

// limitModelResponse 在 SDK 解码前限制响应体大小
func limitModelResponse(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
	resp, err := next(req)
	if err != nil {
		return resp, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("model API returned HTTP %d", resp.StatusCode)
	}
	if resp.ContentLength > maxAIResponse {
		resp.Body.Close()
		return nil, status.ErrModelResponseTooLarge
	}
	return resp, nil
}
