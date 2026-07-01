package main

import (
	"fmt"
	"go-tools/git_merge/git"
	"os"
	"strings"

	"github.com/fatih/color"
)

var validBranches = []string{"dev", "test", "master"}

func isValidBranch(branch string) bool {
	for _, b := range validBranches {
		if b == branch {
			return true
		}
	}
	return false
}

func printValidBranches() {
	fmt.Print("✅ 支持的分支:")
	color.Green(strings.Join(validBranches, ", "))
}

func main() {
	targetBranch := ""

	if len(os.Args) > 1 {
		targetBranch = os.Args[1]
	}

	if !isValidBranch(targetBranch) {
		color.Red("❌ 无效的目标分支'%s'", targetBranch)
		printValidBranches()
		os.Exit(1)
	}

	currentBranch, err := git.GetCurrentBranch()
	if err != nil {
		color.Red("❌ 获取当前分支失败:%v", err)
		os.Exit(1)
	}

	if currentBranch == targetBranch {
		color.Red("❌ 当前分支已是目标分支'%s'", targetBranch)
		os.Exit(1)
	}

	color.Cyan("✅ 当前分支:%s", currentBranch)
	color.Cyan("✅ 目标分支:%s", targetBranch)

	hasChanges, err := git.HasUncommittedChanges()
	if err != nil {
		color.Red("❌ 检查未提交文件失败:%v", err)
		os.Exit(1)
	}

	if hasChanges {
		color.Red("❌ 取消合并:存在未提交的文件或记录")
		os.Exit(1)
	}

	color.Yellow("- 切换到目标分支")
	if err := git.CheckoutBranch(targetBranch); err != nil {
		color.Red("❌ 切换分支失败:%v", err)
		os.Exit(1)
	}

	color.Yellow("- 检查是否有更新")
	hasRemoteChanges, err := git.HasRemoteChanges(targetBranch)
	if err != nil {
		color.Red("❌ 检查远程更新失败:%v", err)
		_ = git.CheckoutBranch(currentBranch)
		os.Exit(1)
	}

	if hasRemoteChanges {
		color.Yellow("- 远程分支有更新,拉取最新代码...")
		if err := git.PullBranch(targetBranch); err != nil {
			color.Red("❌ 拉取代码失败:%v", err)
			_ = git.CheckoutBranch(currentBranch)
			os.Exit(1)
		} else {
			color.Green("✅ 拉取成功")
		}
	}

	color.Yellow("- 合并分支:%s=>%s", currentBranch, targetBranch)
	hasConflict, err := git.MergeBranch(currentBranch)
	if err != nil {
		color.Red("❌ 合并失败:%v", err)
		if hasConflict {
			color.Red("❌ 存在合并冲突,请手动解决后再继续")
		} else {
			_ = git.CheckoutBranch(currentBranch)
		}
		os.Exit(1)
	} else {
		color.Green("✅ 合并成功")
	}
	color.Yellow("- 推送至远程仓库...")
	if err := git.PushBranch(targetBranch); err != nil {
		color.Red("❌ 推送至远程仓库失败:%v", err)
		_ = git.CheckoutBranch(currentBranch)
		os.Exit(1)
	}
	color.Green("✅ 推送成功")
	color.Yellow("- 切回原分支:%s", currentBranch)
	if err := git.CheckoutBranch(currentBranch); err != nil {
		color.Red("❌ 切回原分支失败:%v", err)
		os.Exit(1)
	}

	color.Green("✅ 操作完成!!!")
}
