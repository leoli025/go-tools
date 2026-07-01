package config

import "testing"

func TestNewConfig(t *testing.T) {
	cfg := NewConfig()
	t.Log(cfg)
}
