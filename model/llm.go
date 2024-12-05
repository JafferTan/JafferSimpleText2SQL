package model

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type RequestBody struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Response struct {
	ID      string `json:"id"`
	Choices []struct {
		FinishReason string      `json:"finish_reason"`
		Index        int         `json:"index"`
		Logprobs     interface{} `json:"logprobs"`
		Message      struct {
			Content      string      `json:"content"`
			Role         string      `json:"role"`
			FunctionCall interface{} `json:"function_call"`
			ToolCalls    interface{} `json:"tool_calls"`
		} `json:"message"`
	} `json:"choices"`
	Created           int64       `json:"created"`
	Model             string      `json:"model"`
	Object            string      `json:"object"`
	ServiceTier       interface{} `json:"service_tier"`
	SystemFingerprint interface{} `json:"system_fingerprint"`
	Usage             struct {
		CompletionTokens int `json:"completion_tokens"`
		PromptTokens     int `json:"prompt_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

func CallLLM(prompt string) (Response, error) {
	// 创建 HTTP 客户端
	client := &http.Client{}
	// 构建请求体
	requestBody := RequestBody{
		// 模型列表：https://help.aliyun.com/zh/model-studio/getting-started/models
		Model: "qwen-max",
		Messages: []Message{
			{
				Role:    "system",
				Content: "You are a helpful assistant.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Fatal(err)
	}
	// 创建 POST 请求
	req, err := http.NewRequest("POST", "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	// 设置请求头
	// 若没有配置环境变量，请用百炼API Key将下行替换为：apiKey := "sk-xxx"
	os.Setenv("DASHSCOPE_API_KEY", "")
	apiKey := os.Getenv("DASHSCOPE_API_KEY")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	// 读取响应体
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var res Response
	err = json.Unmarshal(bodyText, &res)
	if err != nil {
		return res, err
	}
	return res, nil
}
