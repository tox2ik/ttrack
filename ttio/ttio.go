package ttio

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	. "genja.org/ttrack/glue"
	. "genja.org/ttrack/model"
)


func parseGnuDate(inputDate string) (time.Time, error) {
	var t time.Time
	var err error
	var out []byte
	formats := [...] string{
		"2018-01-20 04:35:11",
		"12:59:59",
	}
	for _, e := range formats {
		t, err = time.Parse(e, inputDate)
		if err == nil {
			break
		}
	}
	if err != nil {
		// maybe-todo: handle schmuck-os date and winders.
		// The semantics of GNU `date -d` are most useful and you should consider installing GNU coreutils.
		// For more info read `info date`; section 29.7 Relative Items in date strings
		// https://www.gnu.org/software/coreutils/manual/html_node/Relative-items-in-date-strings.html#Relative-items-in-date-strings
		// The intro-quote of section 29 Date input formats is also worth a read.
		out, err = exec.Command("date", "--rfc-email", "-d", inputDate).Output()
		t, err = time.Parse(time.RFC1123Z, strings.Trim(string(out), "\n"))
	}
	if nil == err {
		return t, nil
	}
	Debug("Failed to parse time: %s", inputDate)
	return time.Time{}, err
}


func ParseRecords(file *os.File) ([]Record, Tuples, error) {
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return nil, Tuples{Items: make([]Tuple, 0)}, err
	}
	scanner := bufio.NewScanner(file)
	records := make([]Record, 0)
	for scanner.Scan() {
		line := strings.Trim(strings.ReplaceAll(scanner.Text(), "  ", " "), " ")
		fields := strings.Split(line, " ")
		if len(fields) >= 4 {
			ts := uint32(0)
			if u64, err := strconv.ParseUint(fields[3], 10, 32); err == nil {
				ts = uint32(u64)
			} else {
				Debug("Bad integer in line %s: %s. Error; %s", line, fields[3], err)
			}
			rec := Record{
				Mark:  strings.ToLower(strings.Replace(fields[0], ":", "", 1)),
				Day:   fields[1],
				Time:  fields[2],
				Stamp: ts}
			records = append(records, rec)
		}
	}
	day := "1970-01-01"
	allTuples := make([]Tuple, 0)
	in := uint32(0)
	var last Record
	for i, rec := range records {
		if rec.IsIn() {
			last = rec
			day = rec.Day
			in = rec.Stamp
		}

		if rec.IsOut() {
			if day != rec.Day {
				// todo: should support over-midnight stamps
				Debug("warn: day mismatch for Record: %d (%s,%s)", i, day, rec.Day)
			}
			allTuples = append(
				allTuples,
				Tuple{
					Day: day,
					Seconds: float32(rec.Stamp - in),
					In: last, // Unnecessary to keep reference here.
					Out: rec, //  Should just set isValid if both present.
				})
			day = ""
		}
	}

	if len(records)%2 != 0 {
		Debug("file contains unfinished stamps")
	}
	return records, Tuples{Items: allTuples}, nil
}

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



func AppendLog(args Arguments) {
	stampsFile := Open(args)
	_, tuples, _ := ParseRecords(stampsFile)
	logPath := path.Join( path.Dir(stampsFile.Name()), path.Base(stampsFile.Name()) + ".log")
	logFile, _ := OpenOutputFile(logPath)


	logFile.WriteString(fmt.Sprintf("%s: describe activity...\n", tuples.Last().Format()))
	logFile.Close()
	logFile, _ = OpenOutputFile(logPath)
	lines := fileAsArray(logFile)
	for i := 3; i >= 1; i-- {
		if len(lines) >= i {
			fmt.Println(lines[len(lines)-i])
		}
	}
	os.Stderr.WriteString(logPath)
	cmd := exec.Command(os.Getenv("EDITOR"), logPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	exec.Command("reset")

}


func fileAsArray(file *os.File) []string {
	lines := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func AddStamp(args Arguments) string {
	stampsFile := Open(args)
	if args.DoDry {
		Die(fmt.Errorf("add dry not implemented"))
	}
	stampLine := writeStamp(stampsFile, args.Stamp, args.Mark)
	stampsFile.Close()
	os.Stderr.WriteString(fmt.Sprintf("%s -> %s\n", stampLine, stampsFile.Name()))
	return stampLine
}

func MarkSession(args Arguments) {
	stampsFile := Open(args)
	if args.DoDry {
		Die(fmt.Errorf("mark dry not implemented"))
	}
	writeStamp(stampsFile, args.Stamp, "out")
	writeStamp(stampsFile, args.Stamp, "in")
	stampsFile.Close()
	if args.DoLog {
		AppendLog(args)
	}
}

func writeStamp(out *os.File, stamp time.Time, mark string) string {
	//fmt.Printf("writeStamp(%s, %s, %s)\n", out.Name(), stamp.Format(time.RFC3339), mark)
	var err error
	mark = normalizeMark(mark)
	if len(mark) == 0 {
		lastMark := identifyLastStamp(out.Name())
		mark = determineNextStamp(lastMark)
	} else if mark == "in" || mark == "out" {
		err = enforceSequence(mark, out)
		Die(err)
	} else {
		panic(errors.New(fmt.Sprintf("Invalid Stamp-Mark %s", mark)))
	}

	stampLine := formatTime(stamp, mark)
	_, err = out.WriteString(stampLine+"\n")
	Die(err)
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
		Die(err)
	}
	stat, _ := os.Stat(path)
	o := make([]byte, stat.Size())
	_, err = out.Read(o)
	lastIndex := int(stat.Size()-1)


	if lastIndex <= 0 {
		return "" //  first Record is always in? ... maybe
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

