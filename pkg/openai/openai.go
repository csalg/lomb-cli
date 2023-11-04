package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/avast/retry-go/v4"
)

const (
	timeout           = 20 * time.Second
	CostPer1000Tokens = 0.002
)

type Request struct {
	Messages         []Message `json:"messages"`
	Model            string    `json:"model"`
	Temperature      int       `json:"temperature"`
	MaxTokens        int       `json:"max_tokens"`
	N                int       `json:"n"`
	TopP             int       `json:"top_p"`
	PresencePenalty  int       `json:"presence_penalty"`
	FrequencyPenalty int       `json:"frequency_penalty"`
}

type Message struct {
	Role    MessageRole `json:"role"`
	Content string      `json:"content"`
}

type MessageRole string

const (
	MessageRoleUser      MessageRole = "user"
	MessageRoleAssistant MessageRole = "assistant"
	MessageRoleSystem    MessageRole = "system"
)

func NewRequest() *Request {
	return &Request{
		Messages:         []Message{},
		Model:            "gpt-3.5-turbo",
		Temperature:      0,
		MaxTokens:        2049,
		N:                1,
		TopP:             1,
		PresencePenalty:  0,
		FrequencyPenalty: 0,
	}
}

func (r *Request) AddMessage(m Message) {
	r.Messages = append(r.Messages, m)
}

type Response struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
	Index        int     `json:"index"`
}

func Post(messages []Message, token string) (Response, error) {
	var res Response
	err := retry.Do(
		func() error {
			var err error
			res, err = post(messages, token)
			if err != nil {
				return err
			}
			if res.Choices[0].FinishReason == "length" {
				return fmt.Errorf("response was truncated")
			}
			return nil
		})
	if err != nil {
		return Response{}, fmt.Errorf("error sending request: %w", err)
	}
	return res, nil
}

func post(messages []Message, token string) (Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	url := "https://api.openai.com/v1/chat/completions"
	method := "POST"
	reqBody := NewRequest()
	reqBody.Messages = messages
	payload, err := json.Marshal(reqBody)
	if err != nil {
		return Response{}, fmt.Errorf("error marshalling request body: %w", err)
	}

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, method, url, strings.NewReader(string(payload)))

	if err != nil {
		return Response{}, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return Response{}, fmt.Errorf("error sending request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return Response{}, fmt.Errorf("error reading response body: %w", err)
		}
		return Response{}, fmt.Errorf("error response status code: %d, body: %s", res.StatusCode, body)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return Response{}, fmt.Errorf("error reading response body: %w", err)
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return Response{}, fmt.Errorf("error unmarshalling response body: %w", err)
	}
	return response, nil
}
