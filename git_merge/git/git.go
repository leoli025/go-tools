package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func runGitCmd(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%s: %s", err, strings.TrimSpace(stderr.String()))
	}
	return strings.TrimSpace(stdout.String()), nil
}

func GetCurrentBranch() (string, error) {
	return runGitCmd("branch", "--show-current")
}

func HasUncommittedChanges() (bool, error) {
	output, err := runGitCmd("status", "--porcelain")
	if err != nil {
		return false, err
	}
	return len(output) > 0, nil
}

func FetchBranch(branch string) error {
	_, err := runGitCmd("fetch", "origin", branch)
	return err
}

func CheckoutBranch(branch string) error {
	_, err := runGitCmd("checkout", branch)
	return err
}

func PullBranch(branch string) error {
	_, err := runGitCmd("pull", "origin", branch)
	return err
}

func PushBranch(branch string) error {
	_, err := runGitCmd("push", "origin", branch)
	return err
}

func MergeBranch(branch string) (bool, error) {
	output, err := runGitCmd("merge", "--no-ff", branch)
	if err != nil {
		if strings.Contains(output, "CONFLICT") || strings.Contains(err.Error(), "CONFLICT") {
			return true, fmt.Errorf("合并冲突: %s", output)
		}
		return false, err
	}
	return false, nil
}

func HasRemoteChanges(branch string) (bool, error) {
	if err := FetchBranch(branch); err != nil {
		return false, err
	}

	localHash, err := runGitCmd("rev-parse", branch)
	if err != nil {
		return false, err
	}

	remoteHash, err := runGitCmd("rev-parse", "origin/"+branch)
	if err != nil {
		return false, err
	}

	return localHash != remoteHash, nil
}
