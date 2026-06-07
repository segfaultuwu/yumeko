package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/segfaultuwu/yumeko/internal/tools"
)

const (
	MistralBaseURL = "https://api.mistral.ai/v1"
)

type MistralResp struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

type MistralReq struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type MistralService struct {
	apiKey       string
	systemPrompt string
}

func (m *MistralService) GetModel() string {
	return "mistral-small-latest"
}

type ToolCall struct {
	Tool string         `json:"tool"`
	Args map[string]any `json:"args"`
}

func NewMistralService(apiKey string) *MistralService {
	return &MistralService{
		apiKey: apiKey,
		systemPrompt: `You are Yumeko, a friendly Discord assistant.

Personality:
- You are a catgirl, but keep it subtle.
- You are warm, helpful, and playful.
- Do not overuse "nya", emojis, or roleplay.
- Always answer in the user's language.

Response style:
- Be concise by default.
- For programming questions, give practical code and direct fixes.
- If the user shows an error, explain the cause and provide the corrected code.
- Avoid long introductions.

Tools:
If you need to use a tool, respond ONLY with valid JSON in this format:
{"tool":"tool_name","args":{}}

Do not wrap tool JSON in markdown.
Do not add text before or after tool JSON.

Available tools:
- ping: checks if the tools system works
- time: returns the current server time

If no tool is needed, answer normally.`,
	}
}

func (m *MistralService) GetBaseURL() string {
	return MistralBaseURL
}

func (m *MistralService) GetAPIKey() string {
	return m.apiKey
}

func (m *MistralService) GetHeaders() map[string]string {
	return map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + m.apiKey,
	}
}

func (m MistralService) GetModelURL() string {
	return "/chat/completions"
}

func (m *MistralService) GetRequestBody(prompt string) interface{} {
	return MistralReq{
		Model: m.GetModel(),
		Messages: []Message{
			{
				Role:    "system",
				Content: m.systemPrompt,
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}
}

func (m *MistralService) GetResponseBody() interface{} {
	return &MistralResp{}
}

func (m *MistralService) GetResponse(resp interface{}) string {
	return resp.(*MistralResp).Choices[0].Message.Content
}

func (m *MistralService) Ask(prompt string) (string, error) {
	body := m.GetRequestBody(prompt)

	data, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	url := m.GetBaseURL() + m.GetModelURL()

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return "", err
	}

	for key, value := range m.GetHeaders() {
		req.Header.Set(key, value)
	}

	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return "", fmt.Errorf("mistral error: status %d: %s", res.StatusCode, string(resBody))
	}

	var mistralResp MistralResp
	if err := json.Unmarshal(resBody, &mistralResp); err != nil {
		return "", err
	}

	if len(mistralResp.Choices) == 0 {
		return "", fmt.Errorf("mistral returned no choices")
	}

	return mistralResp.Choices[0].Message.Content, nil
}

func ParseToolCall(text string) (*ToolCall, bool) {
	text = strings.TrimSpace(text)

	if !strings.HasPrefix(text, "{") {
		return nil, false
	}

	var call ToolCall
	if err := json.Unmarshal([]byte(text), &call); err != nil {
		return nil, false
	}

	if call.Tool == "" {
		return nil, false
	}

	if call.Args == nil {
		call.Args = map[string]any{}
	}

	return &call, true
}

func (m *MistralService) AskWithTools(
	ctx context.Context,
	prompt string,
	registry *tools.Registry,
) (string, error) {
	first, err := m.Ask(prompt)
	if err != nil {
		return "", err
	}

	call, ok := ParseToolCall(first)
	if !ok {
		return first, nil
	}

	result, err := registry.Execute(ctx, call.Tool, call.Args)
	if err != nil {
		return "", err
	}

	finalPrompt := fmt.Sprintf(`User asked:
%s

You used tool:
%s

Tool result:
%s

Now answer the user normally in their language.
Do not output JSON.
Do not mention internal tool mechanics unless useful.`, prompt, call.Tool, result)

	return m.Ask(finalPrompt)
}
