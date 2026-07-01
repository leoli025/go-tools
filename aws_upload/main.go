package main

import (
	"bufio"
	"fmt"
	"go-tools/aws_upload/config"
	"go-tools/utils"
	"log"
	"mime"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Config 保存程序的所有配置
type Config struct {
	LocalDir    string   // 要遍历的本地目录
	Bucket      string   // S3 存储桶名称
	Region      string   // AWS 区域
	S3Prefix    string   // S3 上的路径前缀（如 "pro/props"）
	AllowedExt  []string // 允许上传的文件后缀（如 []string{".jpg", ".png"}）
	OutputFile  string   // 输出映射关系的文件路径
	Concurrency int      // 并发上传数
}

// UploadResult 保存单个文件的上传结果
type UploadResult struct {
	LocalPath string // 本地文件绝对路径
	URL       string // 上传后的云端 URL
	Error     error  // 上传过程中的错误
}

// 默认配置
var defaultConfig = Config{
	Bucket:      "",                                                // 存储桶
	Region:      "",                                                // 区域
	S3Prefix:    "pro/props",                                       // S3 路径前缀
	LocalDir:    "./uploads",                                       // 默认当前目录下的 uploads 文件夹
	AllowedExt:  []string{".jpg", ".jpeg", ".png", ".gif", ".svg"}, // 默认图片后缀
	OutputFile:  "upload_mapping.txt",                              // 输出映射文件
	Concurrency: 5,                                                 // 同时上传 5 个文件
}

func main() {
	cfg := defaultConfig
	awsConf := config.NewConfig()
	cfg.Bucket = awsConf.Bucket
	cfg.Region = awsConf.Region
	cfg.S3Prefix = awsConf.S3Prefix
	cfg.LocalDir = awsConf.LocalDir

	log.Printf("开始扫描目录: %s", cfg.LocalDir)

	// 2. 遍历目录，获取所有符合条件的文件
	files, err := walkDir(cfg.LocalDir, cfg.AllowedExt)
	if err != nil {
		log.Fatalf("扫描目录失败: %v", err)
	}
	if len(files) == 0 {
		log.Println("未找到任何需要上传的文件，程序退出。")
		return
	}
	log.Printf("共找到 %d 个符合条件的文件。", len(files))

	// 3. 初始化 AWS Session
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(cfg.Region),
		Credentials: credentials.NewStaticCredentials(awsConf.AccessId, awsConf.AccessSecret, ""),
	})
	if err != nil {
		log.Fatalf("创建 AWS Session 失败: %v", err)
	}
	s3Client := s3.New(sess)

	// 4. 并发上传所有文件
	results := uploadFilesConcurrently(s3Client, files, cfg)

	// 5. 将结果写入文档
	err = writeResultsToFile(results, cfg.OutputFile)
	if err != nil {
		log.Fatalf("写入结果文件失败: %v", err)
	}

	log.Printf("处理完成！共处理 %d 个文件，成功 %d 个，失败 %d 个。映射关系已保存至: %s",
		len(results), countSuccess(results), countFailure(results), cfg.OutputFile)
}

// walkDir 遍历目录，返回所有符合后缀过滤条件的文件绝对路径列表
func walkDir(root string, allowedExts []string) ([]string, error) {
	var files []string
	extMap := make(map[string]bool)
	for _, ext := range allowedExts {
		extMap[strings.ToLower(ext)] = true
	}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("警告: 访问路径 %s 出错: %v", path, err)
			return nil // 跳过有问题的文件/目录
		}
		if info.IsDir() {
			return nil // 跳过目录
		}
		ext := strings.ToLower(filepath.Ext(path))
		if extMap[ext] {
			absPath, err := filepath.Abs(path)
			if err != nil {
				log.Printf("警告: 获取文件 %s 绝对路径失败: %v", path, err)
				return nil
			}
			files = append(files, absPath)
		}
		return nil
	})

	return files, err
}

// generateS3Key 生成带时间戳和随机字符串的 S3 对象键
func generateS3Key(originalFilename string, prefix string) string {
	ext := filepath.Ext(originalFilename)
	timeStr := time.Now().Format("20060102150405")
	newFilename := fmt.Sprintf("aws_driver_%s_%s%s", timeStr, utils.RandAlphaNumber(12), ext)
	return path.Join(prefix, newFilename)
}

// uploadFileToS3 上传单个文件到 S3，返回上传后的 URL 或错误
func uploadFileToS3(client *s3.S3, localPath string, bucket string, s3Prefix string) (string, error) {
	file, err := os.Open(localPath)
	if err != nil {
		return "", fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	// 生成 S3 对象键
	baseName := filepath.Base(localPath)
	s3Key := generateS3Key(baseName, s3Prefix)

	// 尝试自动检测 Content-Type
	contentType := mime.TypeByExtension(filepath.Ext(localPath))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// 准备上传输入
	input := &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(s3Key),
		Body:        file,
		ContentType: aws.String(contentType),
	}

	// 执行上传
	_, err = client.PutObject(input)
	if err != nil {
		return "", fmt.Errorf("上传到 S3 失败: %w", err)
	}

	// 拼接云端 URL (注意: 如果存储桶是私有访问，此URL需要签名才能访问)
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucket, *client.Config.Region, s3Key)
	return url, nil
}

// uploadFilesConcurrently 并发上传多个文件
func uploadFilesConcurrently(client *s3.S3, files []string, cfg Config) []UploadResult {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, cfg.Concurrency)
	results := make([]UploadResult, len(files))

	for i, v := range files {
		wg.Add(1)
		go func(idx int, localPath string) {
			defer wg.Done()
			semaphore <- struct{}{} // 获取信号量

			log.Printf("正在上传: %s", localPath)
			url, err := uploadFileToS3(client, localPath, cfg.Bucket, cfg.S3Prefix)
			results[idx] = UploadResult{
				LocalPath: localPath,
				URL:       url,
				Error:     err,
			}
			if err != nil {
				log.Printf("上传失败: %s, 错误: %v", localPath, err)
			} else {
				log.Printf("上传成功: %s -> %s", localPath, url)
			}

			<-semaphore // 释放信号量
		}(i, v)
	}

	wg.Wait()
	return results
}

// writeResultsToFile 将上传结果写入文档，格式: 本地路径\t云端URL (失败则记录错误)
func writeResultsToFile(results []UploadResult, outputFile string) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString("# 本地文件路径\t\t云端URL / 错误信息\n")
	if err != nil {
		return err
	}

	for _, res := range results {
		var line string
		if res.Error != nil {
			line = fmt.Sprintf("%s\t\t[上传失败] %s\n", res.LocalPath, res.Error.Error())
		} else {
			line = fmt.Sprintf("%s\t\t%s\n", res.LocalPath, res.URL)
		}
		_, err = writer.WriteString(line)
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}

// countSuccess 统计成功上传的数量
func countSuccess(results []UploadResult) int {
	count := 0
	for _, res := range results {
		if res.Error == nil {
			count++
		}
	}
	return count
}

// countFailure 统计失败的数量
func countFailure(results []UploadResult) int {
	return len(results) - countSuccess(results)
}
