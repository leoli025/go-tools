package sheet

import (
	"go-tools/google_sheet/config"
	"testing"
)

func TestFetchSheetData(t *testing.T) {
	cfg := config.NewConfig()
	data, err := FetchSheetData(cfg.SheetId, cfg.RangeName, cfg.ApiKey)
	if err != nil {
		t.Errorf("FetchSheetData failed: %v", err)
	}
	if len(data) == 0 {
		t.Errorf("FetchSheetData returned empty data")
	}
	t.Logf("data: %v\n", data)
}
