package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"

	. "genja.org/ttrack/model"
	"genja.org/ttrack/ttio"
)

var stdErr = log.New(os.Stderr, "", 0)


func act(args Arguments) (err error) {
	var tuples Tuples

	if args.SumPerDay {
		args.DoCount = true
		_, tuples, err = ttio.ParseRecords(ttio.Open(args))
		if err != nil {
			return
		}
		err = tuples.ReportHoursPerDay(os.Stdout)
		if err != nil {
			return
		}
	} else
	if args.DoCount {
		_, tuples, err = ttio.ParseRecords(ttio.Open(args))
		if err != nil {
			return
		}
		err = tuples.ReportHours(os.Stdout)
		if err != nil {
			return
		}
	}

	if args.DoLog {
		ttio.AppendLog(args)
	}

	if args.DoMark {
		ttio.MarkSession(args)
	}

	ls := ""

	addStampByDefault := ! (args.DoCount || args.DoLog)
if addStampByDefault {
	if ! args.Stamp.IsZero() {
			ls = ttio.AddStamp(args)

			showLastTuple(ls, args)
		}
	}
	return err

}

func showLastTuple(stampLine string, args Arguments) {
	if strings.Contains(stampLine, "out:") {
		stampsFile := ttio.Open(args)
		_, tuples, _ := ttio.ParseRecords(stampsFile)
		stampsFile.Close()
		t := tuples.Last()
		fmt.Printf("%s  %5.2f\n", t.Day, t.Seconds/3600)
	}

}

func main() {

	arguments := parseArgs(os.Args[1:])
	err := act(arguments)

	if nil != err {
		stdErr.Print(err)
	}
}

func dieHelp() {
	_, _ = flags.ParseArgs(&Arguments{}, []string{"--help"})
	additionalHelp()
	os.Exit(0)
}
