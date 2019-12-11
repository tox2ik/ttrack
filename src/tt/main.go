package main

import (
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

func createStorageFolder() (string, error) {
	var home = os.Getenv("HOME")
	var direnv = os.Getenv("TIMETRACK_DIR")
	var ttdir string

	if len(ttdir) >= 1 {
		ttdir = direnv
	} else if len(home) > 0 {
		ttdir = sprintf("%s/ttrack", home)
	} else if len(ttdir) + len(home) == 0 {
		ttdir = "/tmp/ttrack"
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


func openOutputFile(t time.Time) (*os.File, error) {
	storageFolder, err := createStorageFolder()
	die(err)
	if len(storageFolder) > 0 {
		monthPath := sprintf("%s/%s", storageFolder, strings.ToLower(t.Format("Jan")))
		var out, err =  os.OpenFile(monthPath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
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
	for lastLn = len(o)-1; lastLn >-1; {
		if o[lastLn] == 10 && lastLn != len(o)-1 {
			break
		}
		lastLn--
	}
	for colon = lastLn; colon < (len(o)-1); {
		if o[colon] == 58 {
			break
		}
		colon++
	}
	if colon > lastLn {
		return strings.Trim(string(o[lastLn:colon]), "\n")
	}
	return "in"
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

func oppositeMark(mark string) string {
	var r string;
	r = "out"

	switch mark {
	case "in": r = "out"
	case "out": r = "in"
	}

	ln("last:" + mark, " opposite:" + r)

	return r
}

func enforceSequence(requestedMark string, out *os.File) error {
	last := identifyLastStamp(out.Name())
	if requestedMark == last {
		msg := sprintf(
			"invalid in/out sequence %s -> %s;\n    - end the last session or start a new one first", last, requestedMark)
		log.New(os.Stderr, "", 0).Println(msg)
		return errors.New(msg)
	}
	return nil
}


func main() {
	var err error
	var out* os.File
	var stamp time.Time
	var inDate string
	var inputMark string

	// binFile := os.Args[1]
	argv := os.Args[1:]

	if len(argv) == 2 {
		inputMark = os.Args[2]
		inDate = os.Args[1]
	} else if len(argv) == 1 {
		inDate = os.Args[1]
		stamp = time.Now()
	} else if len(argv) == 0 {
		stamp = time.Now()
	}


	if len(inDate) > 0 {
		stamp, err = parseInputDate(inDate)
		die(err)
	}

	out, err = openOutputFile(stamp)


	if len(inputMark) > 0 {
		if nil != enforceSequence(inputMark, out) {
			os.Exit(1)
		}
	} else {
		inputMark = oppositeMark(identifyLastStamp(out.Name()))
	}

	stampLine := formatTime(stamp, noramalizeMark(inputMark))
	_, err = out.WriteString(stampLine+"\n")
	die(err)

	out.Sync()
	fmt.Printf("%s -> %s\n", stampLine, out.Name())
}

