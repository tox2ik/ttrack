package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"genja.org/ttrack/glue"
	. "genja.org/ttrack/model"
	"genja.org/ttrack/ttio"
	"genja.org/ttrack/tui"
)

func main() {
	args := ParseArgs(os.Args[1:])
	glue.Die(mainAct(args, os.Stdout))
}

var ew = os.Stderr

func mainAct(args Arguments, ow io.Writer) (err error) {

	{
		var tuples Tuples
		if args.DoSumPerDay || args.DoCount {
			_, tuples, err = ttio.ParseRecords(args, ew)
			if err != nil {
				return
			}
		}

		if args.DoSumPerDay {
			args.DoCount = true
			err = tuples.ReportHoursPerDay(ow)
		} else
		if args.DoCount {
			err = tuples.ReportHours(ow)
		}

		if err != nil {
			return
		}

	}
	{
		if args.DoMark {
			ttio.MarkSession(args, ew)
		}

		if args.DoAddStamp() {
			if strings.Contains(ttio.AddStamp(args, os.Stderr), "out:") {
				_, tuples, _ := ttio.ParseRecords(args, ew)
				fmt.Println(tuples.Last().FormatDur())
			}
		}

	}

	if args.DoInteractive {
		err = tui.Run(args)
	}


	if args.DoLog {
		ttio.AppendLog(args, ew)
	}
	return err

}

