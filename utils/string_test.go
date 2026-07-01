package utils

import "testing"

func TestRandAlphaNumber(t *testing.T) {
	t.Logf("%s", RandAlphaNumber(12))
}
