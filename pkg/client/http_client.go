package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ResponseResult レスポンス結果を保持する構造体
type ResponseResult struct {
	StatusCode   int
	Status       string
	Body         string
	Duration     time.Duration
	Headers      http.Header
}

// HTTPClient HTTPクライアントの構造体
type HTTPClient struct {
	client *http.Client
}

// NewHTTPClient 新しいHTTPクライアントを作成
func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ExecuteRequest HTTPリクエストを実行
func (c *HTTPClient) ExecuteRequest(method, url string, headers map[string]string, body []byte) (*ResponseResult, error) {
	var reqBody io.Reader
	if body != nil {
		reqBody = bytes.NewBuffer(body)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("リクエスト作成エラー: %w", err)
	}

	// ヘッダーを設定
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// デフォルトでContent-Typeを設定（ボディがある場合）
	if body != nil && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	startTime := time.Now()
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("リクエスト実行エラー: %w", err)
	}
	defer resp.Body.Close()

	duration := time.Since(startTime)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("レスポンス読み込みエラー: %w", err)
	}

	return &ResponseResult{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Body:       string(respBody),
		Duration:   duration,
		Headers:    resp.Header,
	}, nil
}
