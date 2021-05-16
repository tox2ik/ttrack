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

	interpArgs(rest, &args)
	return args
}

func interpArgs(argv []string, arg *Arguments) {
	date :=  ""
	tail := argv
	for len(tail) > 0 {
		head := strings.ToLower(tail[0])
		tail = tail[1:]

		if arg.IsStamp(head) {
			date = head
			continue
		}

		switch head {
		case "help":
			dieHelp()
		default:
			if ! arg.SetFlagFromAlias(head) {
				arg.OutPath = head
			}
		}

	}

	if len(date) > 0 {
		argStamp, _ := ParseDate(date)
		arg.Stamp = argStamp
	} else if len(arg.InStamp) > 0 {
		dStamp, _ := ParseDate(arg.InStamp)
		arg.Stamp = dStamp
	} else {
		// the default action is to stamp in or out at the current time.
		// zero := a.Stamp
		// a.Stamp = time.Now()
		// if a.DoCount {
		// 	a.Stamp = zero
		// }
		if ! arg.DoCount {
			arg.Stamp = time.Now()
		}
	}

	if arg.DoCountCurrentMonth() {
		arg.Stamp = time.Now() // count from current-month-file
	}
}

func ParseDate(in string) (time.Time, error) {
	if strings.Contains(in, "today") {
		println(in)
	}
	if t, e := parseGo(in); e == nil {
		return t, nil
	}
	if t, e := parseGnuDate(in); e == nil {
		return t, nil
	}
	return time.Time{}, fmt.Errorf("failed to parse time: %s", in)
}

func parseGo(inputDate string) (time.Time, error) {
	// ref-time: Mon Jan 2 15:04:05 -0700 MST 2006
	layoFull := []string{
		"2006-01-02 15:04",
		"2006-01-02 15:04:05",
	}
	layoShort := []string{
		"15:04:05",
		"15:04",
		"today 15:04",
		"today 15:04:05",
		"yest 15:04",
		"yest 15:04:05",
		"yesterday 15:04",
		"yesterday 15:04:05",
	}
	for _, s := range layoShort {
		t, err := time.Parse(s, inputDate)
		if nil == err {
			if strings.Contains(inputDate, "today") {
				return todayTime(t, 0), nil
			} else
			if strings.Contains(inputDate, "yest") {
				return todayTime(t, -1), nil
			} else {
				return todayTime(t, 0), nil
			}
		}
	}
	for _, s := range layoFull {
		t, err := time.Parse(s, inputDate)
		if nil == err {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unable to parse date")
}

func todayTime(hms time.Time, dayOffset int) time.Time {
	___ymd := time.Now()
	return time.Date(
		___ymd.Year(),
		___ymd.Month(),
		___ymd.Day() + dayOffset,
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
	t, err = time.Parse(time.RFC1123Z, strings.Trim(string(out), "\r\n"))
	if nil == err {
		return t, nil
	}
	return time.Time{}, err
}


func dieHelp() {
	_, _ = flags.ParseArgs(&Arguments{}, []string{"--help"})
	additionalHelp()
	os.Exit(0)
}

func additionalHelp() {
	fmt.Printf(`
The flags also have shortcuts as indicated by the aliases in paranthesis.
This means you can use 'tt log' or 'tt sum' instaed of 'tt -l' and 'tt -s'.

cover:
- today
- yest|yesterday
- gnu-date
`)
}
