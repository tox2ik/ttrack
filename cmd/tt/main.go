package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"genja.org/ttrack/glue"
	. "genja.org/ttrack/model"
	"genja.org/ttrack/ttio"
	"genja.org/ttrack/tui"
)

func main() {
	args := ParseArgs(os.Args[1:])
	err := mainAct(args, os.Stdout)
	glue.Die(err)
}

var ew = os.Stderr

func mainAct(args Arguments, ow io.Writer) (err error) {

	var tuples Tuples
	if args.DoSumPerDay || args.DoCount || args.DoList || args.DoDump {
		_, tuples, err = ttio.ParseRecords(args, ew)
		if err != nil {
			return
		}
	}

	{
		if args.DoList {
			err = tuples.ReportHoursHuman(ow)
			return
		}
		if args.DoDump {
			f := ttio.Open(args)
			var bb []byte
			bb, err = os.ReadFile(f.Name())
			_,_ = ow.Write(bb)
			return
		}

		if args.DoSumPerDay {
			args.DoCount = true
			err = tuples.ReportHoursPerDay(ow)
			return
		}

		if args.DoCount {
			err = tuples.ReportHours(ow)
			return
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

	{

		if args.DoEdit {
			out := []byte{}
			stampsFile := ttio.Open(args)
			edComm := ""
			edArg := []string{}

			editor := strings.TrimSpace(os.Getenv("EDITOR"))
			if strings.Contains(editor, " ") {
				ed := strings.Split(editor, " ")
				if len(ed) >= 1 {
					edComm = ed[0]
					edArg = ed [1:]
				} else if len(ed) == 1 {
					edComm = ed[0]
				}
			} else {
				edComm = editor
			}
			edArg = append(edArg, stampsFile.Name(), "+333")
			out, err = exec.Command(edComm, edArg...).Output()
			if err != nil {
				_,_ = ow.Write([]byte(err.Error()))
			}
			_,_ = ow.Write(out)
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

