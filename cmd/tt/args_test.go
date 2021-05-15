package main

import (
	"testing"
	"time"

	"genja.org/ttrack/model"
)

func TestParseArgs_hms(t *testing.T) {
	n := time.Now()
	t12 := time.Date(n.Year(), n.Month(), n.Day(), 12, 12, 12, 0, n.Location())
	failed := gatherFailed(ParseArgs([]string{"-d", "12:12:12"}), t12)
	if len(failed) > 0 {
		t.Errorf("parse -d hh:mm:ss => some filelds failed: %s", failed)
	}
}

func TestParseArgs_hm(t *testing.T) {
	n := time.Now()
	t1337 := time.Date(n.Year(), n.Month(), n.Day(), 13, 37, 0, 0, n.Location())
	failed := gatherFailed(ParseArgs([]string{"-d", "13:37"}), t1337)
	if len(failed) > 0 {
		t.Errorf("parse -d hh:mm => some filelds failed: %s", failed)
	}
}

func gatherFailed(args model.Arguments, expected time.Time) string {
	failed := ""
	if args.Stamp.Year() != expected.Year() {
		failed += "Y"
	}
	if args.Stamp.Month() != expected.Month() {
		failed += "M"
	}
	if args.Stamp.Day() != expected.Day() {
		failed += "D"
	}
	if len(failed) >0 {
		failed += " "
	}
	if args.Stamp.Hour() != expected.Hour() {
		failed += "h"
	}
	if args.Stamp.Minute() != expected.Minute() {
		failed += "m"
	}
	if args.Stamp.Second() != expected.Second() {
		failed += "s"
	}
	return failed
}
