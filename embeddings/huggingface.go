package embeddings

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type HuggingFace struct {
	ModelName string
	ModelURL  string
}

const ENDPOINT_HUGGINGFACE_EMBEDDING = "/embeddings"

func (h *HuggingFace) Embedding(ctx context.Context, text string) ([]float32, error) {
	// 构成请求内容
	reqBody := struct {
		Input string `json:"input"`
		Model string `json:"model"`
	}{
		Input: text,
		Model: h.ModelName,
	}
	reqBodyBuf, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	// 发送请求
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(reqBodyBuf).
		Post(h.ModelURL + ENDPOINT_HUGGINGFACE_EMBEDDING)
	if err != nil {
		return nil, err
	}

	// 解析响应
	var respBody struct {
		Code  int       `json:"code,omitempty"`
		Model string    `json:"model"`
		Data  []float32 `json:"data"`
		Error string    `json:"error,omitempty"`
	}
	err = json.Unmarshal(resp.Body(), &respBody)
	if err != nil {
		return nil, err
	}
	if respBody.Code != http.StatusOK && respBody.Code != 0 {
		return nil, fmt.Errorf("huggingface api error: '%s'", resp.Body())
	}
	if len(respBody.Data) == 0 {
		return nil, fmt.Errorf("huggingface api error: embedding is empty")
	}

	return respBody.Data, nil
}

func (h *HuggingFace) Embeddings(ctx context.Context, texts []string) ([][]float32, error) {
	// 构成请求内容
	reqBody := struct {
		Input []string `json:"input"`
		Model string   `json:"model"`
	}{
		Input: texts,
		Model: h.ModelName,
	}
	reqBodyBuf, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	// 发送请求
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(reqBodyBuf).
		Post(h.ModelURL + ENDPOINT_HUGGINGFACE_EMBEDDING)
	if err != nil {
		return nil, err
	}

	// 解析响应
	var respBody struct {
		Code  int         `json:"code,omitempty"`
		Model string      `json:"model"`
		Data  [][]float32 `json:"data"`
		Error string      `json:"error,omitempty"`
	}
	err = json.Unmarshal(resp.Body(), &respBody)
	if err != nil {
		return nil, err
	}
	if respBody.Code != http.StatusOK && respBody.Code != 0 {
		return nil, fmt.Errorf("huggingface api error: '%s'", resp.Body())
	}
	if len(respBody.Data) == 0 {
		return nil, fmt.Errorf("huggingface api error: embedding is empty")
	}

	return respBody.Data, nil
}
