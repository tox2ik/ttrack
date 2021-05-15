package glue

import (
	"os"
	"strings"
	"testing"
)

func TestTestStampFile(t *testing.T) {

	eff := TestStampFile()
	fi, err := os.Stat(eff)

	if ! strings.HasSuffix(eff, ".1") {
		t.Error("First file should end with 1")
	}

	if fi != nil {
		t.Error("TestStampFile should not create a file")
	}

	if err == nil {
		t.Error("Stat should err on non-file")
	}
}

