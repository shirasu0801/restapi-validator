package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "restapi_check",
	Short: "APIレスポンスチェックツール",
	Long: `APIレスポンスチェックツールは、指定したエンドポイントにリクエストを送り、
レスポンス（ステータスコード、ボディ、実行時間）を表示するCLIツールです。
AIがレスポンス内容を解析し、エラー解決策の提示やドキュメント生成を自動で行います。`,
}

// Execute コマンドを実行
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
