package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// LoadConfig 環境変数を読み込む
func LoadConfig() error {
	// .envファイルが存在する場合のみ読み込む
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			return fmt.Errorf(".envファイルの読み込みに失敗しました: %w", err)
		}
	}
	return nil
}

// GetAPIKey APIキーを取得
func GetAPIKey(provider string) string {
	key := os.Getenv(fmt.Sprintf("%s_API_KEY", provider))
	return key
}

// GetOpenAIAPIKey OpenAI APIキーを取得
func GetOpenAIAPIKey() string {
	return GetAPIKey("OPENAI")
}

// GetAnthropicAPIKey Anthropic APIキーを取得
func GetAnthropicAPIKey() string {
	return GetAPIKey("ANTHROPIC")
}
