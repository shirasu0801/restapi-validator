package cmd

import (
	"fmt"
	"os"
	"restapi_check/pkg/client"
	"restapi_check/pkg/display"
	"restapi_check/pkg/history"

	"github.com/spf13/cobra"
)

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "リクエスト履歴を管理",
	Long:  `過去に実行したリクエストの履歴を表示・再実行できます。`,
}

var historyListCmd = &cobra.Command{
	Use:   "list",
	Short: "履歴一覧を表示",
	Long:  `保存されているリクエスト履歴の一覧を表示します。`,
	Run:   runHistoryList,
}

var historyShowCmd = &cobra.Command{
	Use:   "show",
	Short: "履歴の詳細を表示",
	Long:  `指定したIDの履歴の詳細を表示します。`,
	Run:   runHistoryShow,
}

var historyReplayCmd = &cobra.Command{
	Use:   "replay",
	Short: "履歴を再実行",
	Long:  `指定したIDの履歴を再実行します。`,
	Run:   runHistoryReplay,
}

func init() {
	rootCmd.AddCommand(historyCmd)
	historyCmd.AddCommand(historyListCmd)
	historyCmd.AddCommand(historyShowCmd)
	historyCmd.AddCommand(historyReplayCmd)

	historyShowCmd.Flags().StringP("id", "i", "", "履歴ID（必須）")
	historyShowCmd.MarkFlagRequired("id")

	historyReplayCmd.Flags().StringP("id", "i", "", "履歴ID（必須）")
	historyReplayCmd.MarkFlagRequired("id")
}

func runHistoryList(cmd *cobra.Command, args []string) {
	manager, err := history.NewHistoryManager()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	histories, err := manager.GetLatest(10)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	if len(histories) == 0 {
		fmt.Println("履歴がありません。")
		return
	}

	fmt.Println("\n=== リクエスト履歴（最新10件） ===")
	for i, h := range histories {
		fmt.Printf("\n[%d] ID: %s\n", i+1, h.ID)
		fmt.Printf("    時刻: %s\n", h.Timestamp.Format("2006-01-02 15:04:05"))
		fmt.Printf("    メソッド: %s\n", h.Method)
		fmt.Printf("    URL: %s\n", h.URL)
		if h.Response != nil {
			fmt.Printf("    ステータス: %d %s\n", h.Response.StatusCode, h.Response.Status)
		}
	}
}

func runHistoryShow(cmd *cobra.Command, args []string) {
	id, _ := cmd.Flags().GetString("id")

	manager, err := history.NewHistoryManager()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	h, err := manager.GetByID(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n=== 履歴詳細 ===")
	fmt.Printf("ID: %s\n", h.ID)
	fmt.Printf("時刻: %s\n", h.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Printf("メソッド: %s\n", h.Method)
	fmt.Printf("URL: %s\n", h.URL)
	
	if len(h.Headers) > 0 {
		fmt.Println("ヘッダー:")
		for k, v := range h.Headers {
			fmt.Printf("  %s: %s\n", k, v)
		}
	}
	
	if h.Body != "" {
		fmt.Printf("ボディ: %s\n", h.Body)
	}

	if h.Response != nil {
		display.FormatResponse(h.Response, h.URL)
	}
}

func runHistoryReplay(cmd *cobra.Command, args []string) {
	id, _ := cmd.Flags().GetString("id")

	manager, err := history.NewHistoryManager()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	h, err := manager.GetByID(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("履歴を再実行中: %s %s\n", h.Method, h.URL)

	var body []byte
	if h.Body != "" {
		body = []byte(h.Body)
	}

	httpClient := client.NewHTTPClient()
	result, err := httpClient.ExecuteRequest(h.Method, h.URL, h.Headers, body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	display.FormatResponse(result, h.URL)

	// 再実行結果を履歴に保存
	newHistory := &history.RequestHistory{
		Method:   h.Method,
		URL:      h.URL,
		Headers:  h.Headers,
		Body:     h.Body,
		Response: result,
	}
	if err := manager.Save(newHistory); err != nil {
		fmt.Fprintf(os.Stderr, "警告: 履歴の保存に失敗しました: %v\n", err)
	}
}
