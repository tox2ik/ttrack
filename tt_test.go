package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const TtStampsFile =  "/tmp/tt-stamps-test"

var tfn = 0

// parallel tests need individual output files
func TtStampFileX() string {
	tfn++
	f := fmt.Sprintf("%s.%d", TtStampsFile, tfn)
	_ = os.Remove(f)
	return f

}

func wipeTestFile(ff... string) {
	_ = os.Remove(TtStampsFile)
	for _, f := range ff {
		if f != "" {
			_ = os.Remove(f)
		}
	}
}

func TestMain(m *testing.M) {
	wipeTestFile()
	code := m.Run()
	wipeTestFile()
	os.Exit(code)
}


func TestCount(t *testing.T) {
	_ = os.Remove("/tmp/tt-should-be-empty")
	count := Arguments{
		DoCount: true,
		OutPath: "/tmp/tt-should-be-empty", // resolved automatically if not specified
	}
	_ = parseAndRun(count)
	_ = parseAndRun(count)
	_ = parseAndRun(count)

	s, err := os.Stat(count.OutPath)

	if err != nil {
		fmt.Println(err)
	}
	assert.Equal(t, 0, int(s.Size()))
}
