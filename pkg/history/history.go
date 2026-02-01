package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"restapi_check/pkg/client"
	"time"
)

// RequestHistory リクエスト履歴の構造体
type RequestHistory struct {
	ID          string                 `json:"id"`
	Timestamp   time.Time              `json:"timestamp"`
	Method      string                 `json:"method"`
	URL         string                 `json:"url"`
	Headers     map[string]string      `json:"headers"`
	Body        string                 `json:"body,omitempty"`
	Response    *client.ResponseResult `json:"response,omitempty"`
}

// HistoryManager 履歴管理の構造体
type HistoryManager struct {
	historyDir string
	historyFile string
}

// NewHistoryManager 新しい履歴マネージャーを作成
func NewHistoryManager() (*HistoryManager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("ホームディレクトリの取得に失敗: %w", err)
	}

	historyDir := filepath.Join(homeDir, ".restapi_check")
	historyFile := filepath.Join(historyDir, "history.json")

	// ディレクトリが存在しない場合は作成
	if err := os.MkdirAll(historyDir, 0755); err != nil {
		return nil, fmt.Errorf("履歴ディレクトリの作成に失敗: %w", err)
	}

	return &HistoryManager{
		historyDir:  historyDir,
		historyFile: historyFile,
	}, nil
}

// Save リクエスト履歴を保存
func (h *HistoryManager) Save(history *RequestHistory) error {
	histories, err := h.LoadAll()
	if err != nil {
		histories = []*RequestHistory{}
	}

	// IDを設定（タイムスタンプベース）
	if history.ID == "" {
		history.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	}
	history.Timestamp = time.Now()

	histories = append(histories, history)

	// 最新100件のみ保持
	if len(histories) > 100 {
		histories = histories[len(histories)-100:]
	}

	data, err := json.MarshalIndent(histories, "", "  ")
	if err != nil {
		return fmt.Errorf("JSONマーシャルエラー: %w", err)
	}

	if err := os.WriteFile(h.historyFile, data, 0644); err != nil {
		return fmt.Errorf("ファイル書き込みエラー: %w", err)
	}

	return nil
}

// LoadAll すべての履歴を読み込む
func (h *HistoryManager) LoadAll() ([]*RequestHistory, error) {
	if _, err := os.Stat(h.historyFile); os.IsNotExist(err) {
		return []*RequestHistory{}, nil
	}

	data, err := os.ReadFile(h.historyFile)
	if err != nil {
		return nil, fmt.Errorf("ファイル読み込みエラー: %w", err)
	}

	var histories []*RequestHistory
	if err := json.Unmarshal(data, &histories); err != nil {
		return nil, fmt.Errorf("JSONアンマーシャルエラー: %w", err)
	}

	return histories, nil
}

// GetByID IDで履歴を取得
func (h *HistoryManager) GetByID(id string) (*RequestHistory, error) {
	histories, err := h.LoadAll()
	if err != nil {
		return nil, err
	}

	for _, history := range histories {
		if history.ID == id {
			return history, nil
		}
	}

	return nil, fmt.Errorf("履歴が見つかりません: %s", id)
}

// GetLatest 最新の履歴を取得
func (h *HistoryManager) GetLatest(n int) ([]*RequestHistory, error) {
	histories, err := h.LoadAll()
	if err != nil {
		return nil, err
	}

	if len(histories) == 0 {
		return []*RequestHistory{}, nil
	}

	start := len(histories) - n
	if start < 0 {
		start = 0
	}

	return histories[start:], nil
}
