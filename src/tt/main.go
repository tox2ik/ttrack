package main

import (
	"../model"
	"../report"
	"./io"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

var stdErr = log.New(os.Stderr, "", 0)

func die(e error) {
	if e != nil {
		panic(e)
	}
}

func parseInputDate(inputDate string) (time.Time, error){
	var t time.Time
	var err error
	formats := [...] string {
		"2018-01-20 04:35:11",
		"12:59:59",
	}
	for _, e := range formats {
		t, err = time.Parse(e, inputDate)
		if err == nil { break }
	}
	if err != nil {
		// maybe-todo: handle schmuck-os date and winders.
		var out []byte
		// The semantics of GNU `date -d` are moset useful and you should consider installing coreutils.
		// For more info read `info date`; section 29.7 Relative items in date strings
		// https://www.gnu.org/software/coreutils/manual/html_node/Relative-items-in-date-strings.html#Relative-items-in-date-strings
		// The intro-quote of section 29 Date input formats is also worth a read.
		out, err = exec.Command("date", "--rfc-email", "-d", inputDate).Output()
		osDate := strings.Trim(fmt.Sprintf("%s", out), "\n")
		t, err = time.Parse(time.RFC1123Z, osDate)
	}
	die(err)
	if nil == err {
		return t, nil;
	}
	return time.Time{}, err
}

func parseArgs(argv []string) model.Arguments {
	var dateString string
	var s string
	var a = model.Arguments{DoCount: false}
	for len(argv) > 0 {
		s = argv[0]
		argv = argv[1:]
		isStamp := (strings.Contains(s, ":") || strings.Contains(s, "-")) && !strings.Contains(s, "/")

		if "count" == s {
			a.DoCount = true
		} else if "per-day" == s || "day" == s {
			a.DoCount = true
			a.SumPerDay = true
		} else if "count-per-day" == s {
			a.DoCount = true
			a.SumPerDay = true
		} else if s == "mark" {
		} else if s == "log" {
			a.DoLog = true
		} else if s ==  "in" || s == "out" {
			a.Mark = s
		} else if isStamp {
			dateString = s
		} else {
			a.OutPath = s
		}
	}
	if len(dateString) > 0 {
		stamp, _ := parseInputDate(dateString)
		a.Stamp = stamp
	} else {
		a.Stamp = time.Now()
	}
	return a
}

func showLastTuple(stampLine string, args model.Arguments) {
	if strings.Contains(stampLine, "out:") {
		stampsFile := io.Open(args)
		_, tuples, _ := report.ParseRecords(stampsFile)
		stampsFile.Close()
		fmt.Print(report.FormatTuple(tuples[len(tuples)-1]))
	}

}

func main() {
	var err error
	args := parseArgs(os.Args[1:])

	if args.DoCount && args.SumPerDay {
		err = report.CountPerDay(io.Open(args))
	} else if args.DoCount {
		err = report.Count(io.Open(args))
	} else if args.DoLog {
		io.AppendLog(args)
	} else if args.DoMark {
		io.MarkSession(args)
	} else  {
		showLastTuple(io.AddStamp(args), args)
	}

	if nil != err {
		stdErr.Print(err)
	}
}

