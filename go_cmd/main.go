package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	//go:embed dev.yaml
	configDev string

	//go:embed pro.yaml
	configPro string
)

type Command struct {
	Name string   `yaml:"name"`
	Desc string   `yaml:"desc"`
	Exec []string `yaml:"exec"`
}

type CommandWrapper struct {
	Command Command `yaml:"command"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("用法: go_cmd <环境> [<命令名>]")
		fmt.Println("环境: dev, pro")
		return
	}

	config := ""
	env := os.Args[1]
	switch env {
	case "dev":
		config = configDev
	case "pro":
		config = configPro
	}
	if config == "" {
		fmt.Printf("环境[%s]的配置文件不存在\n", env)
		fmt.Println("可用环境: dev, pro")
		return
	}

	var commands []CommandWrapper
	err := yaml.Unmarshal([]byte(config), &commands)
	if err != nil {
		fmt.Printf("配置文件解析失败: %v\n", err)
		return
	}

	if len(os.Args) == 2 {
		fmt.Printf("环境[%s]可用命令:\n", env)
		for _, wrapper := range commands {
			fmt.Printf(" - %-15s %s\n", wrapper.Command.Name, wrapper.Command.Desc)
		}
		return
	}

	cmdName := os.Args[2]
	for _, wrapper := range commands {
		if wrapper.Command.Name == cmdName {
			fmt.Printf("执行命令: %s (%s)\n", wrapper.Command.Name, wrapper.Command.Desc)
			for _, execCmd := range wrapper.Command.Exec {
				fmt.Printf("==> %s\n", execCmd)
				err := executeCommand(execCmd)
				if err != nil {
					return
				}
			}
			return
		}
	}

	fmt.Printf("未知命令: %s\n", cmdName)
	fmt.Printf("环境[%s]可用命令:\n", env)
	for _, wrapper := range commands {
		fmt.Printf("  %-15s %s\n", wrapper.Command.Name, wrapper.Command.Desc)
	}
}

func executeCommand(cmdStr string) error {
	parts := strings.Fields(cmdStr)
	if len(parts) == 0 {
		return nil
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		fmt.Printf("--> 命令执行失败: %v\n", err)
	}
	return err
}
