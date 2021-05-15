package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"

	"genja.org/ttrack/glue"
	. "genja.org/ttrack/model"
	"genja.org/ttrack/ttio"
)

func main() {
	glue.Die(act(ParseArgs(os.Args[1:])))
}

func act(args Arguments) (err error) {
	var tuples Tuples

	if args.SumPerDay || args.DoCount {
		_, tuples, err = ttio.ParseRecords(ttio.Open(args))
		if err != nil {
			return
		}
	}

	if args.SumPerDay {
		args.DoCount = true
		err = tuples.ReportHoursPerDay(os.Stdout)
	} else
	if args.DoCount {
		err = tuples.ReportHours(os.Stdout)
	}

	if err != nil {
		return
	}

	if args.DoLog {
		ttio.AppendLog(args)
	}

	if args.DoMark {
		ttio.MarkSession(args)
	}

	if args.DoAddStamp() {
		if strings.Contains(ttio.AddStamp(args), "out:") {
			_, tuples, _ := ttio.ParseRecords(ttio.Open(args))
			fmt.Println(tuples.Last().Format())
		}
	}
	return err

}

func dieHelp() {
	_, _ = flags.ParseArgs(&Arguments{}, []string{"--help"})
	additionalHelp()
	os.Exit(0)
}
