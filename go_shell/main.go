package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

//go:embed init.sh
var scriptContent string

func main() {
	cmd := exec.Command("bash")
	cmd.Stdin = strings.NewReader(scriptContent)
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Println("==> 脚本执行失败:", err)
	}
}
