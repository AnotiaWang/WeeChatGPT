package model

import (
	"ChatGPT_WeChat_Bot/logicerr"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type ChatCompletionResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Param   string `json:"param"`
		Code    string `json:"code"`
	} `json:"error"`
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int32  `json:"created"`
	Model   string `json:"model"`

	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`

	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
		Index        int    `json:"index"`
	}
}

func MakeMessage(content string, sysMsg string) []Message {
	messages := make([]Message, 0)

	if sysMsg != "" {
		messages = append(messages, Message{
			Role:    "system",
			Content: sysMsg,
		})
	}
	messages = append(messages, Message{
		Role:    "user",
		Content: content,
	})

	return messages
}

func ChatCompletion(ctx context.Context, messages []Message) (string, error) {
	config := ctx.Value(ConfigKey).(*Config)

	data, _ := json.Marshal(ChatCompletionRequest{
		Model:    config.OpenAI.Model,
		Messages: messages,
	})

	req, _ := http.NewRequest("POST", config.OpenAI.Endpoint, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.OpenAI.SecretKey)

	client := http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	body := ChatCompletionResponse{}
	err := json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return "", errors.New(logicerr.DecodeJSONFailedError.Error() + ": " + err.Error())
	}

	if body.Error.Message != "" {
		return "", errors.New(logicerr.ChatCompletionFailedError.Error() + ": " + body.Error.Message)
	}

	content := strings.TrimFunc(body.Choices[0].Message.Content, func(r rune) bool {
		return r == ' ' || r == '\t' || r == '\n'
	})

	return content, nil
}
