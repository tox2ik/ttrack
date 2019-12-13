package main

import (
	R "../report"
	A "../arguments"
	"errors"
	"fmt"
	_ "github.com/araddon/dateparse"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

var printf = fmt.Printf
var sprintf = fmt.Sprintf
var echo = fmt.Println
var ln = fmt.Println

func die(e error) {
	if e != nil {
		panic(e)
	}
}

func dump(e interface{}) {
	fmt.Printf("%v -- %T\n", e, e)
}

func createStorageFolder(ttdir string) (string, error) {
	var home = os.Getenv("HOME")
	var direnv = os.Getenv("TIMETRACK_DIR")
	//var ttdir string

	if len(ttdir) == 0 {
		if len(direnv) >= 1 {
			ttdir = direnv
		} else if len(home) > 0 {
			ttdir = sprintf("%s/ttrack", home)
		} else if len(ttdir) + len(home) == 0 {
			ttdir = "/tmp/ttrack"
		}
	}

	_, err := os.Stat(ttdir)
	if err != nil && os.IsNotExist(err) {
		return "", err
	} else {
		if nil != os.MkdirAll(ttdir, 0775) {
			fmt.Println("Failed to create dir: " + ttdir)
			os.Exit(1)
		}
	}
	return ttdir, nil
}

func openOutputFile(path string) (*os.File, error) {
	var out, err =  os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	return out, err
}

func openCurrentMonthOutputFile(t time.Time) (*os.File, error) {
	storageFolder, err := createStorageFolder("")
	die(err)
	if len(storageFolder) > 0 {
		monthPath := sprintf("%s/%s", storageFolder, strings.ToLower(t.Format("Jan")))
		var out, err =  os.OpenFile(monthPath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
		die(err)
		return out, err
	}
	return & os.File{}, err
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
		osDate := strings.Trim(sprintf("%s", out), "\n")
		t, err = time.Parse(time.RFC1123Z, osDate)
	}
	die(err)
	if nil == err {
		return t, nil;
	}
	return time.Time{}, err
}

// if last marker==out: ... => out
// if last marker==in : ... => in
// if no records        ... => in
func identifyLastStamp(path string) string {
	var out* os.File
	var err error
	var colon int
	var lastLn int
	if len(path) > 0 {
		out, err = os.Open(path)
		die(err)
	}
	stat, _ := os.Stat(path)
	o := make([]byte, stat.Size())
	_, err = out.Read(o)
	lastIndex := int(stat.Size()-1)


	if lastIndex <= 0 {
		return "" //  first record is always in? ... maybe
	}

	//var isEmptyFile bool
	//isEmptyFile = lastIndex <= 1
	//
	//if isEmptyFile {
	//	return "out" //  first record is always in? ... maybe
	//}

	if lastIndex >=2 {
		for lastLn = lastIndex; lastLn >-1; {
			if o[lastLn] == 10 && lastLn != lastIndex {
				break
			}
			lastLn--
		}
		if lastLn == -1 {
			lastLn = 0
			colon = 1
		} else {
			colon = lastLn
		}
		for ; colon < lastIndex && colon >= 0; {
			if o[colon] == 58 {
				break
			}
			colon++
		}
		if colon > lastLn {
			return strings.Trim(string(o[lastLn:colon]), "\n")
		}
	}
	return "in"
}



// in : 2007-03-04 12:20:00 1173010800
// out: 2007-03-04 12:20:00 1173010800
func formatTime(t time.Time, mark string) string {
	if mark == "" { mark = "in" }
	return fmt.Sprintf("%-4s %d-%02d-%02d %02d:%02d:%02d %d",
		mark + ":",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second(),
		t.Unix())
}

func noramalizeMark(mark string) string {
	r := "in"
	switch mark {
	case "inn": r = "in"
	case "ut": r = "out"
	case "out": r = "out"
	case "in": r = "in"
	}
	return r
}

func determineNextStamp(mark string) string {
	var r string;
	switch mark {
	case "in": r = "out"
	case "out": r = "in"
	case "": r = "in"
	}
	return r
}

func enforceSequence(requestedMark string, out *os.File) error {
	last := identifyLastStamp(out.Name())
	if requestedMark != "in" && requestedMark != "out" {
		return errors.New(sprintf("Unrecognized marker: %s", requestedMark))
	}
	if requestedMark == last {
		msg := sprintf(
			"invalid in/out sequence %s -> %s;\n    - end the last session or start a new one first",
			last,
			requestedMark)
		//log.New(os.Stderr, "", 0).Println(msg)
		return errors.New(msg)
	}
	return nil
}



func writeStamp(out *os.File, stamp time.Time, mark string) string {
	var err error
	if len(mark) > 0 {
		err := enforceSequence(noramalizeMark(mark), out)
		if nil != err {
			panic(err)
		}
	} else {
		lastMark := identifyLastStamp(out.Name())
		mark = determineNextStamp(lastMark)
	}
	stampLine := formatTime(stamp, noramalizeMark(mark))
	_, err = out.WriteString(stampLine+"\n")
	die(err)
	out.Sync()
	return stampLine
}



func parseArgs(argv []string) A.Arguments {
	var dateString string
	var s string
	var a = A.Arguments{DoCount: false}
	for len(argv) > 0 {
		s = argv[0]
		argv = argv[1:]
		isStamp := (strings.Contains(s, ":") || strings.Contains(s, "-")) && !strings.Contains(s, "/")

		if "count" == s {
			a.DoCount = true
		} else if "per-day" == s {
			a.SumPerDay = true
		} else if len(s) <= 3 {
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

func open(args A.Arguments) *os.File {
	var out* os.File
	if len(args.OutPath) > 0 {
		out, _ = openOutputFile(args.OutPath)
	} else {
		out, _ = openCurrentMonthOutputFile(args.Stamp)
	}
	return out
}



func main() {
	var stampsFile * os.File
	var args A.Arguments = parseArgs(os.Args[1:])

	if args.DoCount {
		R.Count(open(args), args)

	} else  {
		stampsFile = open(args)
		stampLine := writeStamp(stampsFile, args.Stamp, args.Mark)
		_ = stampsFile.Close()
		log.New(os.Stderr, "", 0).Print(fmt.Sprintf("%s -> %s\n", stampLine, stampsFile.Name()))

	}
}




