package sheet

import (
	"testing"
)

func TestFetchSheetData(t *testing.T) {
	sheetId := "xxx"
	rangeName := "Sheet1!A:D"
	apiKey := "xxx"
	data, err := FetchSheetData(sheetId, rangeName, apiKey)
	if err != nil {
		t.Errorf("FetchSheetData failed: %v", err)
	}
	if len(data) == 0 {
		t.Errorf("FetchSheetData returned empty data")
	}
	t.Logf("data: %v\n", data)
}
