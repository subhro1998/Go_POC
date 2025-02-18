package excelreader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadExcelSuccess(t *testing.T) {
	excelAllData := ReadExcel("user_data.xlsx")
	if len(excelAllData) != 3 {
		t.Errorf("Expected Excel Data length of 3, but got %d", len(excelAllData))
	}
}

func TestReadExcelFileNotFound(t *testing.T) {
	assert.Panics(t, func() { ReadExcel("XYZ.xlsx") }, "The code did not panic")
}

func TestExtractPrivilegesSuccess(t *testing.T) {
	row := []string{"Name", "Email", "Role", "Privilege1, Privilege2    ", "password"}
	output := extractPrivileges(row, 1)

	assert.Equal(t, 2, len(output))
	assert.ElementsMatch(t, output, []string{"Privilege1", "Privilege2"})
}

func TestExtractPrivilegesBlankPrivileges(t *testing.T) {
	row := []string{"Name", "Email", "Role", "", "password"}
	output := extractPrivileges(row, 1)

	assert.Equal(t, 0, len(output))
}
