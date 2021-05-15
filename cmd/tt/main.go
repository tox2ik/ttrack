package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	. "genja.org/ttrack/model"
	"genja.org/ttrack/ttio"
)

var stdErr = log.New(os.Stderr, "", 0)


func parseAndRun(args Arguments) (err error) {
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


	if ! (args.DoCount || args.DoLog) {
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
	// fmt.Printf("%#v\n\n\n", runPar)
	//fmt.Println("tt.main()")

	args := parseArgs(os.Args[1:])
	if err := parseAndRun(args); nil != err {
		stdErr.Print(err)
	}
}
