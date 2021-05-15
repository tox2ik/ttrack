package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"genja.org/ttrack/glue"
	. "genja.org/ttrack/model"
)


func TestMain(m *testing.M) {
	glue.WipeTestFiles()
	code := m.Run()
	glue.WipeTestFiles()
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
