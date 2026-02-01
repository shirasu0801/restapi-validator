package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"restapi_check/pkg/ai"
	"restapi_check/pkg/client"
	"restapi_check/pkg/config"
	"restapi_check/pkg/display"
	"restapi_check/pkg/history"
	"strings"

	"github.com/spf13/cobra"
)

var (
	urlFlag      string
	methodFlag   string
	headerFlags  []string
	bodyFlag     string
	bodyFileFlag string
	analyzeFlag  bool
)

var requestCmd = &cobra.Command{
	Use:   "request",
	Short: "HTTPリクエストを実行",
	Long: `指定したURLに対してHTTPリクエストを実行し、レスポンスを表示します。
メソッド、ヘッダー、ボディを指定できます。`,
	Run: runRequest,
}

func init() {
	rootCmd.AddCommand(requestCmd)
	
	requestCmd.Flags().StringVarP(&urlFlag, "url", "u", "", "リクエスト先のURL（必須）")
	requestCmd.Flags().StringVarP(&methodFlag, "method", "m", "GET", "HTTPメソッド (GET, POST, PUT, DELETE, PATCH)")
	requestCmd.Flags().StringArrayVarP(&headerFlags, "header", "H", []string{}, "HTTPヘッダー (例: -H 'Authorization: Bearer token')")
	requestCmd.Flags().StringVarP(&bodyFlag, "body", "b", "", "リクエストボディ（JSON文字列）")
	requestCmd.Flags().StringVarP(&bodyFileFlag, "body-file", "f", "", "リクエストボディファイルのパス")
	requestCmd.Flags().BoolVarP(&analyzeFlag, "analyze", "a", false, "エラー時にAI分析を自動実行")
	
	requestCmd.MarkFlagRequired("url")
}

func runRequest(cmd *cobra.Command, args []string) {
	if urlFlag == "" {
		fmt.Fprintln(os.Stderr, "エラー: URLを指定してください (--url または -u)")
		os.Exit(1)
	}

	// ヘッダーをパース
	headers := make(map[string]string)
	for _, header := range headerFlags {
		parts := strings.SplitN(header, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			headers[key] = value
		}
	}

	// ボディを取得
	var body []byte
	if bodyFileFlag != "" {
		var err error
		body, err = os.ReadFile(bodyFileFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: ボディファイルの読み込みに失敗しました: %v\n", err)
			os.Exit(1)
		}
	} else if bodyFlag != "" {
		// JSONとして有効かチェック
		var jsonData interface{}
		if err := json.Unmarshal([]byte(bodyFlag), &jsonData); err != nil {
			fmt.Fprintf(os.Stderr, "エラー: ボディが有効なJSONではありません: %v\n", err)
			os.Exit(1)
		}
		body = []byte(bodyFlag)
	}

	// HTTPリクエストを実行
	httpClient := client.NewHTTPClient()
	result, err := httpClient.ExecuteRequest(methodFlag, urlFlag, headers, body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	// 結果を表示
	display.FormatResponse(result, urlFlag)

	// 履歴に保存
	historyManager, err := history.NewHistoryManager()
	if err == nil {
		historyEntry := &history.RequestHistory{
			Method:   methodFlag,
			URL:      urlFlag,
			Headers:  headers,
			Body:     string(body),
			Response: result,
		}
		if err := historyManager.Save(historyEntry); err != nil {
			fmt.Fprintf(os.Stderr, "警告: 履歴の保存に失敗しました: %v\n", err)
		}
	}

	// エラー時にAI分析を実行（オプション）
	if analyzeFlag && (result.StatusCode >= 400) {
		fmt.Println("\n=== AIエラー分析を実行中... ===")
		if err := config.LoadConfig(); err != nil {
			fmt.Fprintf(os.Stderr, "警告: %v\n", err)
		}

		aiClient, err := ai.NewOpenAIClient()
		if err != nil {
			fmt.Fprintf(os.Stderr, "AI分析スキップ: %v\n", err)
			return
		}

		analysis, err := aiClient.AnalyzeError(result.StatusCode, result.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "AI分析エラー: %v\n", err)
			return
		}

		fmt.Println(analysis)
	}
}
