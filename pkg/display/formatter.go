package display

import (
	"encoding/json"
	"fmt"
	"restapi_check/pkg/client"
)

// FormatResponse レスポンス結果を整形して表示
func FormatResponse(result *client.ResponseResult, url string) {
	fmt.Println("\n=== APIレスポンス結果 ===")
	fmt.Printf("URL: %s\n", url)
	fmt.Printf("ステータスコード: %d %s\n", result.StatusCode, result.Status)
	fmt.Printf("実行時間: %d ms\n", result.Duration.Milliseconds())
	
	// JSONを整形して表示
	if result.Body != "" {
		var jsonData interface{}
		if err := json.Unmarshal([]byte(result.Body), &jsonData); err == nil {
			// JSONとして有効な場合、整形して表示
			formattedJSON, err := json.MarshalIndent(jsonData, "", "  ")
			if err == nil {
				fmt.Println("レスポンスボディ:")
				fmt.Println(string(formattedJSON))
			} else {
				fmt.Printf("レスポンスボディ:\n%s\n", result.Body)
			}
		} else {
			// JSONでない場合、そのまま表示
			fmt.Printf("レスポンスボディ:\n%s\n", result.Body)
		}
	}
}
