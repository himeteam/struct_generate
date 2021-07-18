package struct_generate

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseFile(t *testing.T) {
	d, _ := os.Getwd()
	testDataPkg := filepath.Join(d, "test_data")

	ParseFile(testDataPkg)
}
