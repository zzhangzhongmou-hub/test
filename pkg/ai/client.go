package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"test/configs"
	"time"
)

type Client struct {
	httpClient *http.Client
	config     configs.AIConfig
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AIRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
}

type ChatRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
	Stream    bool      `json:"stream"`
}

type ChatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func NewClient() *Client {
	cfg := configs.Cfg.AI

	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 10
	}

	return &Client{
		httpClient: &http.Client{Timeout: time.Duration(cfg.Timeout) * time.Second},
		config:     cfg,
	}
}

func (c *Client) EvaluateCode(ctx context.Context, codeContent string) (comment string, score int, err error) {
	if !c.config.Enable {
		return "", 0, fmt.Errorf("AI功能未启用")
	}

	prompt := fmt.Sprintf(`你是一位资深编程导师，请对以下学生代码/作业进行评价。
要求：
	1. 给出简要评语（优点+改进建议，100字以内）
	2. 给出建议分数（0-100的整数）
	3. 返回格式必须严格如下，不要有多余内容：
	评语：{你的评语}
	分数：{数字}

学生提交内容：

%s`, codeContent)

	reqBody := ChatRequest{
		Model: c.config.Model,
		Messages: []Message{
			{Role: "system", Content: "你是一位严格的编程导师，善于发现代码问题并给出建设性意见。只返回指定格式，不要多余解释。"},
			{Role: "user", Content: prompt},
		},
		MaxTokens: c.config.MaxTokens,
		Stream:    false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", 0, fmt.Errorf("请求序列化失败: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.config.BaseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", 0, err
	}

	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", 0, fmt.Errorf("AI请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", 0, fmt.Errorf("AI接口错误，状态码: %d", resp.StatusCode)
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", 0, fmt.Errorf("解析响应失败: %w", err)
	}

	if chatResp.Error != nil {
		return "", 0, fmt.Errorf("AI接口错误: %s", chatResp.Error.Message)
	}

	if len(chatResp.Choices) == 0 {
		return "", 0, fmt.Errorf("AI无响应内容")
	}
	
	content := chatResp.Choices[0].Message.Content
	return parseAIResponse(content)
}
