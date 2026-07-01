package main

import (
	"encoding/json"
	"fmt"
	"go-tools/google_sheet/sheet"
	"log"
	"os"
)

// GenerateTranslationFiles 生成翻译文件
func GenerateTranslationFiles(translations map[string]map[string]string, outputDir string) error {
	// 创建输出目录
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	// 获取所有语言
	languages := make(map[string]bool)
	for _, langMap := range translations {
		for lang := range langMap {
			if lang != "key" {
				languages[lang] = true
			}
		}
	}

	// 为每种语言生成 JSON 文件
	for lang := range languages {
		langData := make(map[string]string)
		for key, langMap := range translations {
			if text, ok := langMap[lang]; ok {
				langData[key] = text
			}
		}

		filePath := fmt.Sprintf("%s/%s.json", outputDir, lang)
		file, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("无法创建文件 %s: %v", filePath, err)
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(langData); err != nil {
			return fmt.Errorf("无法编码 JSON: %v", err)
		}

		log.Printf("已生成翻译文件: %s", filePath)
	}

	return nil
}

func main() {
	sheetId := "xxx"
	rangeName := "Sheet1!A:D"
	apiKey := "xxx"
	// 获取表格数据
	data, err := sheet.FetchSheetData(sheetId, rangeName, apiKey)
	if err != nil {
		log.Fatalf("获取表格数据失败: %v", err)
	}

	// 解析翻译数据
	translations := sheet.ParseTranslationData(data)
	if len(translations) == 0 {
		log.Println("没有找到翻译数据")
		return
	}

	// 生成翻译文件
	if err := GenerateTranslationFiles(translations, "./output"); err != nil {
		log.Fatalf("生成翻译文件失败: %v", err)
	}

	log.Println("翻译文件生成完成!")
}
