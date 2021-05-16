package main

import (
	"testing"
	"time"

	"genja.org/ttrack/model"
)

type dateTestInput struct {
	y int
	m time.Month
	d int

	h int
	i int
	s int

	T string
}

func TestParseArgs_Date(t *testing.T) {
	n := time.Now()
	year := n.Year()
	m := n.Month()
	d := n.Day()

	h := n.Hour()
	i := n.Minute()
	s := n.Second()

	samplesSimple := []dateTestInput{
		{y: year, m: m, d: d, h: 12, i: 12, s: 12, T: "12:12:12"},
		{y: year, m: m, d: d, h: 13, i: 37, s: 00, T: "13:37"},
		{y: 1990, m: 2, d: 3, h: 04, i: 05, s: 06, T: "1990-02-03 04:05:06"},
		{y: 2260, m: 6, d: 6, h: 06, i: 06, s: 06, T: "2260-06-06 06:06:06"},
	}

	samplesNatural := []dateTestInput{
		{y: year, m: m, d: d, h: 13, i: 30, s: 00, T: "today 13:30"},
		{y: year, m: m, d: d - 1, h: 18, i: 50, s: 00, T: "yesterday 18:50"},
		{y: year, m: m, d: d - 1, h: 18, i: 55, s: 33, T: "yest 18:55:33"},
	}

	samplesRelative := []dateTestInput{
		{y: year, m: m, d: d, h: h, i: i - 3, s: s, T: "3 min ago"},
		// this is fed to gnu date, no need to test further examples.
	}

	samples := []dateTestInput{}
	samples = append(samples, samplesSimple...)
	samples = append(samples, samplesNatural...)
	samples = append(samples, samplesRelative...)

	for _, is := range samples {
		shouldBe := time.Date(is.y, is.m, is.d, is.h, is.i, is.s, 0, n.Location())
		args := ParseArgs([]string{"-d", is.T})
		failed := gatherFailed(args, shouldBe)
		if len(failed) > 0 {
			t.Errorf("parse -d {%s} failed for fields: %s\n\t=> %q", is.T, failed, args.Stamp)
		}
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
	if len(failed) > 0 {
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
