package cmd

import (
	"fmt"
	"os"
	"restapi_check/pkg/ai"
	"restapi_check/pkg/config"

	"github.com/spf13/cobra"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "AIでレスポンスを分析",
	Long: `AIを使用してレスポンスを分析します。
エラー解決アドバイス、スキーマ生成、テストデータ推論などが可能です。`,
}

var analyzeErrorCmd = &cobra.Command{
	Use:   "error",
	Short: "エラーレスポンスを分析",
	Long:  `4xx/5xxエラー時に、AIが原因を推測し「修正すべき箇所」を提案します。`,
	Run:   runAnalyzeError,
}

var generateStructCmd = &cobra.Command{
	Use:   "struct",
	Short: "Go構造体を生成",
	Long:  `レスポンスのJSONから、Goの構造体（Struct）を自動生成します。`,
	Run:   runGenerateStruct,
}

var generateSchemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "OpenAPIスキーマを生成",
	Long:  `レスポンスのJSONから、OpenAPI（Swagger）定義を自動生成します。`,
	Run:   runGenerateSchema,
}

var suggestTestCmd = &cobra.Command{
	Use:   "test",
	Short: "テストデータを推論",
	Long:  `レスポンスに基づき、次にテストすべき境界値（Boundary Value）をAIが提案します。`,
	Run:   runSuggestTest,
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
	analyzeCmd.AddCommand(analyzeErrorCmd)
	analyzeCmd.AddCommand(generateStructCmd)
	analyzeCmd.AddCommand(generateSchemaCmd)
	analyzeCmd.AddCommand(suggestTestCmd)

	analyzeErrorCmd.Flags().IntP("status", "s", 0, "ステータスコード（必須）")
	analyzeErrorCmd.Flags().StringP("body", "b", "", "レスポンスボディ（必須）")
	analyzeErrorCmd.MarkFlagRequired("status")
	analyzeErrorCmd.MarkFlagRequired("body")

	generateStructCmd.Flags().StringP("body", "b", "", "レスポンスボディ（必須）")
	generateStructCmd.MarkFlagRequired("body")

	generateSchemaCmd.Flags().StringP("body", "b", "", "レスポンスボディ（必須）")
	generateSchemaCmd.MarkFlagRequired("body")

	suggestTestCmd.Flags().StringP("body", "b", "", "レスポンスボディ（必須）")
	suggestTestCmd.MarkFlagRequired("body")
}

func runAnalyzeError(cmd *cobra.Command, args []string) {
	if err := config.LoadConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "警告: %v\n", err)
	}

	statusCode, _ := cmd.Flags().GetInt("status")
	body, _ := cmd.Flags().GetString("body")

	client, err := ai.NewOpenAIClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	result, err := client.AnalyzeError(statusCode, body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n=== AIエラー分析結果 ===")
	fmt.Println(result)
}

func runGenerateStruct(cmd *cobra.Command, args []string) {
	if err := config.LoadConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "警告: %v\n", err)
	}

	body, _ := cmd.Flags().GetString("body")

	client, err := ai.NewOpenAIClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	result, err := client.GenerateGoStruct(body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n=== 生成されたGo構造体 ===")
	fmt.Println(result)
}

func runGenerateSchema(cmd *cobra.Command, args []string) {
	if err := config.LoadConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "警告: %v\n", err)
	}

	body, _ := cmd.Flags().GetString("body")

	client, err := ai.NewOpenAIClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	result, err := client.GenerateOpenAPISchema(body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n=== 生成されたOpenAPIスキーマ ===")
	fmt.Println(result)
}

func runSuggestTest(cmd *cobra.Command, args []string) {
	if err := config.LoadConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "警告: %v\n", err)
	}

	body, _ := cmd.Flags().GetString("body")

	client, err := ai.NewOpenAIClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	result, err := client.SuggestTestData(body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n=== AIテストデータ推論 ===")
	fmt.Println(result)
}
