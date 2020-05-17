package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jessevdk/go-flags"
)

var stdErr = log.New(os.Stderr, "", 0)

func die(e error) {
	if e != nil {
		panic(e)
	}
}

type Arguments struct {
	DoCount   bool `short:"c" long:"count" description:"Count stamps"`
	DoLog     bool `short:"l" long:"log" description:"Describe last time stamp"`
	DoMark    bool `short:"m" long:"mark" description:"Sign out and back in"`
	DoDry     bool `short:"n" long:"dry" description:"Dry run"`
	SumPerDay bool `short:"s" long:"sum" description:"Count average per day"`

	Stamp   time.Time `short:"d" long:"date" description:"Time to record"`
	OutPath string
	Mark    string `short:"r" long:"record" description:"Sign in or out" choice:"in" choice:"out"`
}

func parseArgs(argv []string) Arguments {
	var dateString string
	var s string
	var a = Arguments{DoCount: false}
	for len(argv) > 0 {
		s = argv[0]
		argv = argv[1:]

		isStamp := !strings.Contains(s, "/") && (
			strings.Contains(s, ":") ||
			strings.Contains(s, "-") ||
			strings.Contains(s, "+") ||
			strings.Contains(s, "now") ||
			strings.Contains(s, "hour") ||
			strings.Contains(s, "minute") ||
			strings.Contains(s, "days") ||
			strings.Contains(s, "Day"))

		if "count" == s || "c" == s {
			a.DoCount = true
		} else if "per-Day" == s || "Day" == s {
			a.DoCount = true
			a.SumPerDay = true
		} else if "count-per-Day" == s {
			a.DoCount = true
			a.SumPerDay = true
		} else if s == "Mark" {
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

func showLastTuple(stampLine string, args Arguments) {
	if strings.Contains(stampLine, "out:") {
		stampsFile := Open(args)
		_, tuples, _ := ParseRecords(stampsFile)
		stampsFile.Close()
		t := lastTuple(tuples)
		fmt.Printf("%s  %5.2f\n", t.Day, t.Seconds/3600)
	}

}

func main() {
	var err error
	args := parseArgs(os.Args[1:])

	if args.DoCount && args.SumPerDay {
		err = CountPerDay(Open(args))
	} else if args.DoCount {
		err = Count(Open(args))
	} else if args.DoLog {
		AppendLog(args)
	} else if args.DoMark {
		MarkSession(args)
	} else  {
		showLastTuple(AddStamp(args), args)
	}

	if nil != err {
		stdErr.Print(err)
	}
}

