package sheet

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type SheetResponse struct {
	Values [][]interface{} `json:"values"`
}

// FetchSheetData 获取表格数据
func FetchSheetData(sheetId string, rangeName string, apiKey string) ([][]interface{}, error) {
	// 构建 API 请求 URL
	apiUrl := fmt.Sprintf(
		"https://sheets.googleapis.com/v4/spreadsheets/%s/values/%s?key=%s",
		sheetId,
		rangeName,
		apiKey,
	)

	// 解析代理地址
	proxyURL, err := url.Parse("http://127.0.0.1:7897")
	if err != nil {
		panic(err)
	}

	// 创建 Transport，使用指定代理
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   time.Second * 10,
	}

	req, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/149.0.0.0 Safari/537.36")
	req.Header.Set("Content-Type", "application/json")

	// 发送 HTTP GET 请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP 请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("请求失败，状态码: %d，响应: %s", resp.StatusCode, string(body))
	}

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 解析 JSON 响应
	var sheetResp SheetResponse
	if err := json.Unmarshal(body, &sheetResp); err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %v", err)
	}

	if len(sheetResp.Values) == 0 {
		return nil, fmt.Errorf("表格中没有数据")
	}

	return sheetResp.Values, nil
}

// ParseTranslationData 解析翻译数据
func ParseTranslationData(data [][]interface{}) map[string]map[string]string {
	translations := make(map[string]map[string]string)

	if len(data) < 2 {
		return translations
	}

	// 获取表头（第一行）
	headers := data[0]
	langColumns := make(map[int]string)

	// 解析表头，找出语言列
	for i, header := range headers {
		if headerStr, ok := header.(string); ok {
			langColumns[i] = headerStr
		}
	}

	// 解析数据行
	for i := 1; i < len(data); i++ {
		row := data[i]
		if len(row) == 0 {
			continue
		}

		key, ok := row[0].(string)
		if !ok {
			continue
		}

		translations[key] = make(map[string]string)

		for colIdx, lang := range langColumns {
			if colIdx < len(row) {
				if text, ok := row[colIdx].(string); ok {
					translations[key][lang] = text
				}
			}
		}
	}

	return translations
}
