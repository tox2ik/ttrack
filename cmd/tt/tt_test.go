package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"genja.org/ttrack/glue"
	. "genja.org/ttrack/model"
)

var ow io.Writer = bytes.NewBuffer(nil)

func TestMain(m *testing.M) {
	glue.WipeTestFiles()
	code := m.Run()
	glue.WipeTestFiles()
	os.Exit(code)
}


func TestCount(t *testing.T) {
	_ = os.Remove("/tmp/tt-should-be-empty")
	args := Arguments{
		DoCount: true,
		OutPath: "/tmp/tt-should-be-empty", // resolved automatically if not specified
	}

	_ = mainAct(args, ow)

	s, err := os.Stat(args.OutPath)

	if err != nil {
		fmt.Println(err)
	}
	assert.Equal(t, 0, int(s.Size()))
}
