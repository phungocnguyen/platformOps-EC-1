package services

import (
	"github.com/tealeg/xlsx"
	"testing"
)

func createTestSheet(cellValues []string) *xlsx.File {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row

	file = xlsx.NewFile()

	sheet, _ = file.AddSheet("Sheet1")
	row = sheet.AddRow()

	for _, c := range cellValues {
		cell := row.AddCell()
		cell.Value = c
		cell.NumFmt = "general"
	}

	return file

}

func TestExcelHeaderParsing(t *testing.T) {

	file := createTestSheet([]string{"Req #", "Category", "Requirements", "Discussion", "Check Text", "Fix Text"})
	file2 := createTestSheet([]string{"Req #", "Category", "BAR", "Requirements", "Discussion", "Check Text", "FOO", "Fix Text"})

	sheet := file.Sheets[0]
	sheet2 := file2.Sheets[0]

	var startrow, reqIdInd, CategoryInd, RequirementsInd, DiscussionInd, CheckTextInd, FixTextInd int
	startrow, reqIdInd, CategoryInd, RequirementsInd, DiscussionInd, CheckTextInd, FixTextInd = parseHeaders(sheet)

	if startrow != 0 {
		t.Fatalf("Failed to identify the starting row")
	}
	if reqIdInd != 0 {
		t.Fatalf("Expected %d, but got %d", 0, reqIdInd)
	}
	if CategoryInd != 1 {
		t.Fatalf("Expected %d, but got %d", 1, CategoryInd)
	}
	if RequirementsInd != 2 {
		t.Fatalf("Expected %d, but got %d", 2, RequirementsInd)
	}
	if DiscussionInd != 3 {
		t.Fatalf("Expected %d, but got %d", 3, DiscussionInd)
	}
	if CheckTextInd != 4 {
		t.Fatalf("Expected %d, but got %d", 4, CheckTextInd)
	}
	if FixTextInd != 5 {
		t.Fatalf("Expected %d, but got %d", 5, FixTextInd)
	}

	startrow, reqIdInd, CategoryInd, RequirementsInd, DiscussionInd, CheckTextInd, FixTextInd = parseHeaders(sheet2)

	if startrow != 0 {
		t.Fatalf("Failed to identify the starting row")
	}
	if reqIdInd != 0 {
		t.Fatalf("Expected %d, but got %d", 0, reqIdInd)
	}
	if CategoryInd != 1 {
		t.Fatalf("Expected %d, but got %d", 1, CategoryInd)
	}
	if RequirementsInd != 3 {
		t.Fatalf("Expected %d, but got %d", 3, RequirementsInd)
	}
	if DiscussionInd != 4 {
		t.Fatalf("Expected %d, but got %d", 4, DiscussionInd)
	}
	if CheckTextInd != 5 {
		t.Fatalf("Expected %d, but got %d", 5, CheckTextInd)
	}
	if FixTextInd != 7 {
		t.Fatalf("Expected %d, but got %d", 7, FixTextInd)
	}

}

func TestLoadFromExcel(t *testing.T) {
	t.Skip("skipping test.")
	b, controls := LoadFromExcel("../path-to-some-file.xlsx")

	if b.Name == "" {
		t.Fatalf("Load failed")
	}
	if controls == nil {
		t.Fatalf("Failed to load the controls")
	}

}
