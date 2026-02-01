package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"restapi_check/pkg/config"
)

// OpenAIClient OpenAI APIクライアント
type OpenAIClient struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// Message OpenAI APIのメッセージ構造
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest OpenAI APIのリクエスト構造
type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// ChatResponse OpenAI APIのレスポンス構造
type ChatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

// NewOpenAIClient 新しいOpenAIクライアントを作成
func NewOpenAIClient() (*OpenAIClient, error) {
	apiKey := config.GetOpenAIAPIKey()
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEYが設定されていません。.envファイルに設定してください")
	}

	return &OpenAIClient{
		apiKey:  apiKey,
		baseURL: "https://api.openai.com/v1/chat/completions",
		client:  &http.Client{},
	}, nil
}

// AnalyzeError エラーレスポンスを分析して解決策を提案
func (c *OpenAIClient) AnalyzeError(statusCode int, responseBody string) (string, error) {
	prompt := fmt.Sprintf(`以下のHTTPエラーレスポンスを分析し、原因と解決策を日本語で提案してください。

ステータスコード: %d
レスポンスボディ:
%s

以下の形式で回答してください:
1. エラーの原因
2. 修正すべき箇所
3. 推奨される解決策`, statusCode, responseBody)

	return c.Chat(prompt)
}

// GenerateGoStruct JSONレスポンスからGo構造体を生成
func (c *OpenAIClient) GenerateGoStruct(responseBody string) (string, error) {
	prompt := fmt.Sprintf(`以下のJSONレスポンスから、Go言語の構造体（Struct）を生成してください。
フィールド名は適切なGoの命名規則に従ってください（大文字で始まる、キャメルケース）。

JSON:
%s

構造体のみを出力してください。`, responseBody)

	return c.Chat(prompt)
}

// GenerateOpenAPISchema JSONレスポンスからOpenAPIスキーマを生成
func (c *OpenAIClient) GenerateOpenAPISchema(responseBody string) (string, error) {
	prompt := fmt.Sprintf(`以下のJSONレスポンスから、OpenAPI 3.0形式のスキーマ定義を生成してください。

JSON:
%s

OpenAPIスキーマのみを出力してください。`, responseBody)

	return c.Chat(prompt)
}

// SuggestTestData テストデータの推論を提案
func (c *OpenAIClient) SuggestTestData(responseBody string) (string, error) {
	prompt := fmt.Sprintf(`以下のAPIレスポンスに基づいて、次にテストすべき境界値（Boundary Value）やテストケースを提案してください。

レスポンス:
%s

以下の形式で回答してください:
1. 境界値テストの提案
2. 異常系テストの提案
3. パフォーマンステストの提案`, responseBody)

	return c.Chat(prompt)
}

// Chat OpenAI APIにチャットリクエストを送信
func (c *OpenAIClient) Chat(prompt string) (string, error) {
	reqBody := ChatRequest{
		Model: "gpt-4o-mini",
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("リクエストのマーシャルに失敗: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("リクエスト作成エラー: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("リクエスト実行エラー: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("APIエラー (ステータス: %d): %s", resp.StatusCode, string(body))
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", fmt.Errorf("レスポンスのデコードに失敗: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("レスポンスに選択肢がありません")
	}

	return chatResp.Choices[0].Message.Content, nil
}
