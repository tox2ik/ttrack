package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/jessevdk/go-flags"

	. "genja.org/ttrack/model"
)

func additionalHelp() {
	fmt.Printf(`
The flags also have shortcuts as indicated by the aliases in paranthesis.
This means you can use «tt log» or «tt sum» instaed of «tt -l» and «tt -s».

`)

}

func ParseArgs(argv []string) Arguments {

	args := Arguments{}
	rest, err := flags.ParseArgs(&args, argv)
	ferr, ok := err.(*flags.Error)
	if ok && ferr.Type == flags.ErrHelp {
		additionalHelp()
		os.Exit(0)
	} else
	if err != nil {
		fmt.Printf("Nope: %s\n", err)
		os.Exit(1)
	}

	guessArgs(rest, &args)
	return args
}

func guessArgs(argv []string, a *Arguments) {
	var dateString string
	var head string

	tail := argv

	for len(tail) > 0 {
		head = strings.ToLower(tail[0])
		tail = tail[1:]

		isStamp := !strings.Contains(head, "/") && (
			strings.Contains(head, ":") || // should be regex [0-9]:[0-9]
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

		switch head {
		case "count":
		case "c":
			a.DoCount = true
		case "sum":
		case "cc":
			a.DoCount = true
			a.SumPerDay = true
		case "mark":
			a.DoMark = true
		case "log":
			a.DoLog = true
		case "in":
		case "out":
			a.Mark = head
		case "help":
			dieHelp()
		default:
			a.OutPath = head // check is file ?
		}

	}

	if len(dateString) > 0 {
		stamp, _ := ParseDate(dateString)
		a.Stamp = stamp
	} else if len(a.InStamp) > 0 {
		stmp, _ := ParseDate(a.InStamp)
		a.Stamp = stmp
	} else {
		// the default action is to stamp in or out at the current time.
		// zero := a.Stamp
		// a.Stamp = time.Now()
		// if a.DoCount {
		// 	a.Stamp = zero
		// }

		if ! a.DoCount {
			a.Stamp = time.Now()
		}
	}

	if a.DoCount && len(a.OutPath) == 0 {
		a.Stamp = time.Now() // count from current-month-file

	}
}

func ParseDate(in string) (time.Time, error) {
	if t, e := parseGo(in); e == nil {
		return t, nil
	}
	if t, e := parseGnuDate(in); e == nil {
		return t, nil
	}
	return time.Time{}, fmt.Errorf("failed to parse time: %s", in)
}

func parseGo(inputDate string) (time.Time, error) {
	var t time.Time
	var err error
	// ref-time: Mon Jan 2 15:04:05 -0700 MST 2006
	layFull1 := "2006-01-02 15:04"
	layFull2 := "2006-01-02 15:04:05"
	layShort1 := "15:04:05"
	layShort2 := "15:04"
	// short
	t, err = time.Parse(layShort2, inputDate)
	if err == nil {
		return todayTime(t), nil
	}
	t, err = time.Parse(layShort1, inputDate)
	if err == nil {
		return todayTime(t), nil
	}
	// full
	t, err = time.Parse(layFull2, inputDate)
	if err == nil {
		return t, nil
	}
	t, err = time.Parse(layFull1, inputDate)
	if err == nil {
		return t, nil
	}
	return t, fmt.Errorf("unable to parse date")
}

func todayTime (hms time.Time) time.Time {
	___ymd := time.Now()
	return time.Date(
		___ymd.Year(),
		___ymd.Month(),
		___ymd.Day(),
		hms.Hour(),
		hms.Minute(),
		hms.Second(),
		0,
		hms.Location())
}

func parseGnuDate(inputDate string) (time.Time, error) {
	var t time.Time
	var err error
	var out []byte

	comment := `
maybe-todo: handle schmuck-os date and winders.
The semantics of GNU 'date -d' are very useful.
For more info read 'info date'; section 29.7 Relative Items in date strings
https://www.gnu.org/software/coreutils/manual/html_node/Relative-items-in-date-strings.html#Relative-items-in-date-strings
The intro-quote of section 29 Date input formats is also worth a read.
`
	comment += ""

	out, err = exec.Command("date", "--rfc-email", "-d", inputDate).Output()
	t, err = time.Parse(time.RFC1123Z, strings.Trim(string(out), "\n"))
	if nil == err {
		return t, nil
	}
	return time.Time{}, err
}
