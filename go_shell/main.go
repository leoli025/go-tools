package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

//go:embed scripts
var scriptFS embed.FS

func main() {
	args := os.Args[1:]
	fmt.Println("[--- 开始遍历嵌入的 .sh 文件 ---]")
	err := walkEmbedFS(scriptFS, "scripts", func(path string, d fs.DirEntry) error {
		// 只处理 .sh 结尾的文件
		if !d.IsDir() && strings.HasSuffix(path, ".sh") {
			fmt.Println("=================================")
			if err := executeEmbeddedScript(path, args...); err != nil {
				fmt.Printf("❌ 执行失败 [%s]: %v\n", path, err)
				return err
			}
			fmt.Printf("✅ 执行成功: %s\n", path)
		}
		return nil
	})

	if err != nil {
		panic(err)
	}
	fmt.Println("=================================")
}

func walkEmbedFS(fSys fs.FS, root string, fn func(path string, d fs.DirEntry) error) error {
	entries, err := fs.ReadDir(fSys, root)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		tmp := path.Join(root, entry.Name())
		if err := fn(tmp, entry); err != nil {
			return err
		}
		// 如果是目录，递归遍历
		if entry.IsDir() {
			if err := walkEmbedFS(fSys, tmp, fn); err != nil {
				return err
			}
		}
	}
	return nil
}

func executeEmbeddedScript(embedPath string, args ...string) error {
	content, err := scriptFS.ReadFile(embedPath)
	if err != nil {
		return fmt.Errorf("读取嵌入文件失败: %w", err)
	}

	// 创建临时文件来存储脚本内容，因为 exec.Command 需要文件路径或明确的解释器
	tmpDir, err := os.MkdirTemp("", "go-embed-sh-*")
	if err != nil {
		return fmt.Errorf("创建临时目录失败: %w", err)
	}
	defer func() {
		_ = os.RemoveAll(tmpDir) // 清理临时文件
	}()

	scriptPath := filepath.Join(tmpDir, "run.sh")
	if err := os.WriteFile(scriptPath, content, 0755); err != nil {
		return fmt.Errorf("写入临时脚本文件失败: %w", err)
	}

	// 根据操作系统选择执行方式
	var cmd *exec.Cmd
	// 合并脚本路径和传入的参数
	cmdArgs := append([]string{scriptPath}, args...)
	if runtime.GOOS == "windows" {
		// Windows 下通常没有原生的 sh 需要提前安装并添加到环境变量中
		cmd = exec.Command("sh", cmdArgs...)
	} else {
		// Linux/macOS 直接使用 sh
		cmd = exec.Command("sh", cmdArgs...)
	}

	// 捕获输出
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("命令执行错误: %w, 输出: %s", err)
	}
	return nil
}
