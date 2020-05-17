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
		}

		if head == "log" {
			a.DoLog = true
		}

		if head == "in" || head == "out" {
			a.Mark = head
		}

		a.OutPath = head
	}

	if len(dateString) > 0 {
		stamp, _ := parseGnuDate(dateString)
		a.Stamp = stamp
	} else {
		// the default action is to stamp in or out at the current time.
		a.Stamp = time.Now()
	}
	return tail
}

func parseArgs(argv []string) Arguments {

	aa := Arguments{}
	// aa.Stamp = time.Now()

	rest, err := flags.ParseArgs(&aa, os.Args[1:])

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
	runPar := parseArgs(os.Args[1:])
	fmt.Printf("%#v\n\n\n", runPar)

	if runPar.SumPerDay {
		runPar.DoCount = true
		err = CountPerDay(Open(runPar))
	} else if runPar.DoCount {
		err = Count(Open(runPar))
	}

	if runPar.DoLog {
		AppendLog(runPar)
	}

	if runPar.DoMark {
		MarkSession(runPar)
	}

	ls := ""

	if ! runPar.Stamp.IsZero() {
		ls = AddStamp(runPar)
	}
	showLastTuple(ls, runPar)

	if nil != err {
		stdErr.Print(err)
	}
}
