package main

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var out *os.File

func TestMain(m *testing.M) {
	_ = os.Remove("/tmp/ttrack")
	code := m.Run()

	_ = os.Remove("/tmp/ttrack")
	os.Exit(code)
}

func expectPanic(t *testing.T) {
	if r := recover(); r == nil {
		t.Errorf("Should have panicked. %v", r)
	}
}


// todo: use a memory buffer https://stackoverflow.com/questions/40316052/in-memory-file-for-testing

func TestFirstStamp(t *testing.T) {
	last := identifyLastStamp("/dev/null")
	if determineNextStamp(last) != "in" {
		t.Errorf("first state must be \"in\"")
	}
}

func TestSecondStamp(t *testing.T) {
	out, _ = OpenOutputFile("/tmp/ttrack")
	writeStamp(out, time.Now(), "")
	writeStamp(out, time.Now(), "")
	if identifyLastStamp("/tmp/ttrack") != "out" {
		t.Errorf("out follows in")
	}
}

func TestOutOfSequenceA(t *testing.T) {
	defer expectPanic(t)
	out, _ = OpenOutputFile("/tmp/ttrack")
	writeStamp(out, time.Now(), "in")
	writeStamp(out, time.Now(), "in")
}

func TestOutOfSequenceB(t *testing.T) {
	defer expectPanic(t)
	out, _ = OpenOutputFile("/tmp/ttrack")
	writeStamp(out, time.Now(), "out")
	writeStamp(out, time.Now(), "out")
}

func TestInvalidMark(t *testing.T) {
	defer expectPanic(t)
	out, _ = OpenOutputFile("/tmp/ttrack")
	writeStamp(out, time.Now(), "typo")
}


func TestSupportUtasOut(t *testing.T) {
	out, _ = OpenOutputFile("/tmp/ttrack")
	writeStamp(out, time.Now(), "inn")
	s := writeStamp(out, time.Now(), "ut")
	if ! strings.Contains(s, "out") {
		t.Errorf("should have converted ut to out, inn to in")
	}
}

func TestHalfHour(t *testing.T) {
	out, _ = OpenOutputFile("/tmp/ttrack")
	writeStamp(out, time.Now(), "in")
	writeStamp(out, time.Now().Add(time.Duration(time.Minute * 30)), "out")
	_ = out.Sync()
	_, tuples, err := ParseRecords(out)
	assert.Empty(t, err)
	assert.Equal(t, float32(1800), tuples.Seconds(), "30 min should be 1800 Seconds")
}

func TestUnseekable(t *testing.T) {
	var file *os.File
	_, _, _ = ParseRecords(file)
	expectPanic(t)
}

/*
func TestOpenMonth(t *testing.T) {
	month := Open(Arguments{})
	fi, _ := month.Stat()
	fm := fi.Mode().Perm()
	fmt.Printf("%b\n", fm)
	fmt.Printf("%o\n", fm)
	fmt.Printf("%v\n", (fm & 0o200) == 0o200)
	fmt.Printf("%v\n", (fm & 0o020) == 0o020)
	fmt.Printf("%v\n", (fm & 0o002) == 0o002)
	fmt.Printf("%o\n", fm)

	fmt.Printf("---\n")
	fmt.Printf("%4d  %5o  %010b\n", 493, 493, 493)
	fmt.Printf("%4d  %5o  %010b\n", 40,  40,  40)
	fmt.Printf("%4d  %5o  %010b\n", 5,   5,   5)
	fmt.Printf("---------------\n")
	fmt.Printf("%4d  %5o  %010b\n", 0o755, 0o755, 0o755 )
	fmt.Printf("%4d  %5o  %010b\n", 0o50,  0o50,  0o50 )
	fmt.Printf("%4d  %5o  %010b\n", 0o5,   0o5,   0o5 )
	fmt.Printf("---------------\n")

	fmt.Printf("%4d  %5o  %010b\n", 1<<uint(0), 1<<uint(0), 1<<uint(0))
	fmt.Printf("%4d  %5o  %010b\n", 1<<uint(1), 1<<uint(1), 1<<uint(1))
	fmt.Printf("%4d  %5o  %010b\n", 1<<uint(2), 1<<uint(2), 1<<uint(2))

}
*/
