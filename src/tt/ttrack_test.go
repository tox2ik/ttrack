package main

import (
	"os"
	io "./io"
	"strings"
	"testing"
	"time"
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


func TestFirstStamp(t *testing.T) {
	last := identifyLastStamp("/dev/null")
	if determineNextStamp(last) != "in" {
		t.Errorf("first state must be \"in\"")
	}
}

func TestSecondStamp(t *testing.T) {
	out, _ = io.OpenOutputFile("/tmp/ttrack")
	writeStamp(out, time.Now(), "")
	writeStamp(out, time.Now(), "")
	if identifyLastStamp("/tmp/ttrack") != "out" {
		t.Errorf("out follows in")
	}
}

func TestOutOfSequenceA(t *testing.T) {
	defer expectPanic(t)
	out, _ = io.OpenOutputFile("/tmp/ttrack")
	writeStamp(out, time.Now(), "in")
	writeStamp(out, time.Now(), "in")
}

func TestOutOfSequenceB(t *testing.T) {
	defer expectPanic(t)
	out, _ = io.OpenOutputFile("/tmp/ttrack")
	writeStamp(out, time.Now(), "out")
	writeStamp(out, time.Now(), "out")
}

func TestInvalidMark(t *testing.T) {
	defer expectPanic(t)
	out, _ = io.OpenOutputFile("/tmp/ttrack")
	writeStamp(out, time.Now(), "typo")
}


func TestSupportUtasOut(t *testing.T) {
	//defer expectPanic(t)
	out, _ = io.OpenOutputFile("/tmp/ttrack")
	s := writeStamp(out, time.Now(), "ut")
	if ! strings.Contains(s, "out") {
		t.Errorf("should have converted ut to out")
	}
}
