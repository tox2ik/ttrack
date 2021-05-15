package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/jessevdk/go-flags"

	"genja.org/ttrack/glue"
	. "genja.org/ttrack/model"
)


func guessArgs(argv []string, a *Arguments) []string {
	var dateString string
	var head string

	tail := argv

	for len(tail) > 0 {
		head = strings.ToLower(tail[0])
		tail = tail[1:]

		isStamp := !strings.Contains(head, "/") && (
			strings.Contains(head, ":") || // regex [0-9]:[0-9]
				strings.Contains(head, "now") ||
				strings.Contains(head, "hour") ||
				strings.Contains(head, "hours") ||
				strings.Contains(head, "second") ||
				strings.Contains(head, "seconds") ||
				strings.Contains(head, "tomorrow") ||
				strings.Contains(head, "yesterday") ||
				strings.Contains(head, "minute") ||
				strings.Contains(head, "years") ||
				strings.Contains(head, "year") ||
				strings.Contains(head, "days") ||
				strings.Contains(head, "day"))

		if isStamp {
			dateString = head
			continue
		}

		if "count" == head || "c" == head {
			a.DoCount = true
			continue
		}

		if "per-day" == head || "sum" == head || "count-per-pay" == head || "cc" == head {
			a.DoCount = true
			a.SumPerDay = true
			continue
		}

		if head == "mark" {
			a.DoMark = true
			continue
		}

		if head == "log" {
			a.DoLog = true
			continue
		}

		if head == "in" || head == "out" {
			a.Mark = head
			continue
		}

		a.OutPath = head
	}

	if len(dateString) > 0 {
		stamp, _ := parseGnuDate(dateString)
		a.Stamp = stamp
	} else {
		// the default action is to stamp in or out at the current time.
		zero := a.Stamp
		a.Stamp = time.Now()
		if a.DoCount {
			a.Stamp = zero
		}
	}

	if a.DoCount && len(a.OutPath) == 0 {
		a.Stamp = time.Now() // count from current-month-file

	}

	return tail
}

func parseArgs(argv []string) Arguments {

	aa := Arguments{}
	rest, err := flags.ParseArgs(&aa, argv)

	rest = guessArgs(rest, &aa)

	if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
		os.Exit(0)
	}

	if err != nil {
		fmt.Printf("Nope: %s\n", err)
		os.Exit(1)
	}

	return aa
}


func parseGnuDate(inputDate string) (time.Time, error) {
	var t time.Time
	var err error
	var out []byte
	formats := [...] string{
		"2018-01-20 04:35:11",
		"12:59:59",
	}
	for _, e := range formats {
		t, err = time.Parse(e, inputDate)
		if err == nil {
			break
		}
	}
	if err != nil {
		// maybe-todo: handle schmuck-os date and winders.
		// The semantics of GNU `date -d` are very useful.
		// For more info read `info date`; section 29.7 Relative Items in date strings
		// https://www.gnu.org/software/coreutils/manual/html_node/Relative-items-in-date-strings.html#Relative-items-in-date-strings
		// The intro-quote of section 29 Date input formats is also worth a read.
		out, err = exec.Command("date", "--rfc-email", "-d", inputDate).Output()
		t, err = time.Parse(time.RFC1123Z, strings.Trim(string(out), "\n"))
	}
	if nil == err {
		return t, nil
	}
	glue.Debug("Failed to parse time: %s", inputDate)
	return time.Time{}, err
}

