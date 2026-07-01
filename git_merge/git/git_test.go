package git

import "testing"

func TestGetCurrentBranch(t *testing.T) {
	currentBranch, err := GetCurrentBranch()
	if err != nil {
		t.Errorf("❌ 获取当前分支失败: %v", err)
	}
	t.Logf("✅ 当前分支:%s", currentBranch)
}

func TestHasRemoteChanges(t *testing.T) {
	hasChanges, err := HasRemoteChanges("dev")
	if err != nil {
		t.Errorf("❌ 检查远程变更失败: %v", err)
	}
	t.Logf("✅ 远程变更:%v", hasChanges)
}
