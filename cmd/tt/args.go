package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jessevdk/go-flags"

	. "genja.org/ttrack/model"
	"genja.org/ttrack/ttio"
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
	date := ""
	tail := argv

	for len(tail) > 0 {
		head := strings.ToLower(tail[0])
		tail = tail[1:]

		if arg.IsStamp(head) {
			date = head
			continue
		}

		switch head {
		case "h", "help":
			dieHelp()
		default:
			if ! arg.SetFlagFromAlias(head) {
				arg.OutPath = head
			}
		}

	}

	if len(arg.InStamp) > 0 && len(date) == 0 {
		date = arg.InStamp
	}

	if len(date) > 0 {
		argStamp, _ := ttio.ParseDate(date)
		arg.Stamp = argStamp
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
