package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

func Open(args Arguments) *os.File {
	var out* os.File
	if len(args.OutPath) > 0 {
		out, _ = OpenOutputFile(args.OutPath)
	} else {
		out, _ = OpenCurrentMonthOutputFile(args.Stamp)
	}
	return out
}

func createStorageFolder(ttdir string) (string, error) {
	var home = os.Getenv("HOME")
	var direnv = os.Getenv("TIMETRACK_DIR")
	//var ttdir string

	if len(ttdir) == 0 {
		if len(direnv) >= 1 {
			ttdir = direnv
		} else if len(home) > 0 {
			ttdir = fmt.Sprintf("%s/ttrack", home)
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

func OpenOutputFile(path string) (*os.File, error) {
	var out, err =  os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	return out, err
}

func OpenCurrentMonthOutputFile(t time.Time) (*os.File, error) {
	storageFolder, err := createStorageFolder("")
	if len(storageFolder) > 0 {
		monthPath := fmt.Sprintf("%s/%s", storageFolder, strings.ToLower(t.Format("Jan")))
		var out, err =  os.OpenFile(monthPath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
		return out, err
	}
	return & os.File{}, err
}

func lastTuple(tc tuples) tuple {
	length := len(tc.items)
	if length == 0 {
		return tuple{}
	}
	return tc.items[length-1]
}


func AppendLog(args Arguments) {
	stampsFile := Open(args)
	_, tuples, _ := ParseRecords(stampsFile)
	logPath := path.Dir(stampsFile.Name()) + "/" +  path.Base(stampsFile.Name()) + ".log"
	logFile, _ := OpenOutputFile(logPath)
	logFile.WriteString(strings.TrimSpace(FormatTuple(lastTuple(tuples))) + ": describe activity...\n")
	logFile.Close()
	logFile, _ = OpenOutputFile(logPath)
	lines := FileAsArray(logFile)
	for i := 3; i >= 1; i-- {
		if len(lines) >= i {
			fmt.Println(lines[len(lines)-i])
		}
	}
	stdErr.Println(logPath)
	cmd := exec.Command(os.Getenv("EDITOR"), logPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	exec.Command("reset")

}

func AddStamp(args Arguments) string {
	stampsFile := Open(args)
	stampLine := writeStamp(stampsFile, args.Stamp, args.Mark)
	stampsFile.Close()
	stdErr.Printf("%s -> %s\n", stampLine, stampsFile.Name())
	return stampLine
}

func MarkSession(args Arguments) {
	stampsFile := Open(args)
	writeStamp(stampsFile, args.Stamp, "out")
	writeStamp(stampsFile, args.Stamp, "in")
	stampsFile.Close()
	if args.DoLog {
		AppendLog(args)
	}
}

func writeStamp(out *os.File, stamp time.Time, mark string) string {
	var err error
	mark = normalizeMark(mark)
	if len(mark) == 0 {
		lastMark := identifyLastStamp(out.Name())
		mark = determineNextStamp(lastMark)
	} else if mark == "in" || mark == "out" {
		err = enforceSequence(mark, out)
		die(err)
	} else {
		panic(errors.New(fmt.Sprintf("Invalid stamp-mark %s", mark)))
	}

	stampLine := formatTime(stamp, mark)
	_, err = out.WriteString(stampLine+"\n")
	die(err)
	out.Sync()
	return stampLine
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

func normalizeMark(mark string) string {
	r := mark
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
		return errors.New(fmt.Sprintf("Unrecognized marker: %s", requestedMark))
	}
	if requestedMark == last {
		msg := fmt.Sprintf(
			"invalid in/out sequence %s -> %s;\n    - end the last session or start a new one first",
			last,
			requestedMark)
		return errors.New(msg)
	}
	return nil
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

