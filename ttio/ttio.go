package ttio

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	. "genja.org/ttrack/glue"
	. "genja.org/ttrack/model"
)

func ParseRecords(args Arguments, ew io.Writer) ([]Record, Tuples, error) {
	return ParseRecordsFile(Open(args), ew)
}

func ParseRecordsFile(file *os.File, ew io.Writer) ([]Record, Tuples, error) {
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return nil, Tuples{Items: make([]Tuple, 0)}, err
	}
	scanner := bufio.NewScanner(file)
	records := make([]Record, 0)
	for scanner.Scan() {
		rec := ParseRecord(scanner.Text())
		if rec.IsValid() {
			records = append(records, rec)
		}

	}
	day := "1970-01-01"
	allTuples := make([]Tuple, 0)
	in := int64(0)
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
				fmt.Fprintf(ew, "warn: day mismatch for Record: %d (%s,%s)\n", i, day, rec.Day)
			}
			allTuples = append(allTuples, Tuple{
				Day:     day,
				Seconds: rec.Stamp - in,
				In:      last, // Unnecessary to keep reference here.
				Out:     rec,  // Should just set isValid if both present.
			})
			day = ""
		}
	}

	if len(records)%2 != 0 {
		fmt.Fprintf(ew, "file contains unfinished stamps")
		lastRec := records[len(records)-1]
		if lastRec.IsValid() && lastRec.IsIn() {
			lastTup := Tuple{
				In:      lastRec,
				Out:     Record{Mark: "out", Day: lastRec.Day, Time: " ... "},
				Seconds: time.Now().Unix() - lastRec.Stamp,
				Day:     lastRec.Day }
			allTuples = append(allTuples, lastTup)
		}
	}
	return records, Tuples{Items: allTuples}, nil
}

func ParseRecord(line string) Record {
	line = strings.ReplaceAll(line, "  ", " ")
	line = strings.Trim(line, " ")
	fields := strings.Split(line, " ")
	if len(fields) >= 4 {
		ts := int64(0)
		if u64, err := strconv.ParseUint(fields[3], 10, 32); err == nil {
			ts = int64(u64)
		} else {
			Debug("Bad integer in line %s: %s. Error; %s", line, fields[3], err)
		}
		return Record{
			Mark:  strings.ToLower(strings.Replace(fields[0], ":", "", 1)),
			Day:   fields[1],
			Time:  fields[2],
			Stamp: ts}
	}
	return Record{}
}

// Open a file with records
func Open(args Arguments) *os.File {
	var out *os.File
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
		} else if len(ttdir)+len(home) == 0 {
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

// write to an existing file. Ask to create a new file
func OpenOutputFile(fPath string) (*os.File, error) {

	isTest := strings.HasPrefix(path.Base(os.Args[0]), "___Test")
	isInteractive := os.Getenv("tt_yes") == ""

	if isInteractive && ! isTest {
		if ! IsExists(fPath) {
			yes := "no"
			fmt.Printf("Really write to %s?\nyes/no: ", fPath)
			rd := bufio.NewReader(os.Stdin)
			input, _ := rd.ReadString('\n')
			yes = strings.ToLower(strings.TrimRight(input, "\r\n"))

			if ! (yes == "yes" || yes == "y") {
				println("aborting. (set tt_yes=1 to skip this question)")
				os.Exit(0)
			}
		}
	}

	var out, err = os.OpenFile(fPath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	return out, err
}

func OpenCurrentMonthOutputFile(t time.Time) (*os.File, error) {
	storageFolder, err := createStorageFolder("")
	if len(storageFolder) > 0 {
		monthPath := fmt.Sprintf("%s/%s", storageFolder, strings.ToLower(t.Format("Jan")))
		var out, err = os.OpenFile(monthPath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
		return out, err
	}
	return &os.File{}, err
}

func AppendLog(args Arguments, ew io.Writer) {
	_, tuples, _ := ParseRecords(args, ew)

	records := Open(args)
	logPath := path.Join(path.Dir(records.Name()), path.Base(records.Name())) + ".log"

	logFile, _ := OpenOutputFile(logPath)
	logFile.WriteString(fmt.Sprintf("%s: describe activity...\n", tuples.Last().FormatDur()))
	lines := fileAsArray(logFile)
	for i := 3; i >= 1; i-- {
		if len(lines) >= i {
			fmt.Println(lines[len(lines)-i])
		}
	}

	fmt.Fprintf(ew, "%s\n", logPath)
	_ = RunEditor(logPath)
	_ = exec.Command("reset").Run()

}

func fileAsArray(file *os.File) []string {
	lines := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func AddStamp(args Arguments, w io.Writer) string {
	stampsFile := Open(args)
	if args.DoDry {
		Die(fmt.Errorf("add dry not implemented"))
	}
	stampLine := writeStamp(stampsFile, args.Stamp, args.Mark)
	stampsFile.Close()
	fmt.Fprintf(w, "%s -> %s\n", stampLine, stampsFile.Name())
	return stampLine
}


func MarkSession(args Arguments, ew io.Writer) {
	stampsFile := Open(args)
	if args.DoDry {
		Die(fmt.Errorf("mark dry not implemented"))
	}
	writeStamp(stampsFile, args.Stamp, "out")
	writeStamp(stampsFile, args.Stamp, "in")
	stampsFile.Close()
	if args.DoLog {
		AppendLog(args, ew)
	}
}

func writeStamp(out *os.File, stamp time.Time, mark string) string {
	// fmt.Printf("writeStamp(%s, %s, %s)\n", out.Name(), stamp.Format(time.RFC3339), mark)
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
	_, err = out.WriteString(stampLine + "\n")
	Die(err)
	out.Sync()
	return stampLine
}

func RemoveFileRecord(args Arguments, rm Record) {
	stampsFile := Open(args)
	scanner := bufio.NewScanner(stampsFile)
	j := -1
	i := 0
	for scanner.Scan() {
		txt := scanner.Text()
		rec := ParseRecord(txt)
		i++
		if rec.IsValid() && rec.Equals(rm) {
			j = i
		}
	}
	stampsFile.Close()
	if j >= 0 {
		err := removeLines(stampsFile.Name(), j, 1)
		Die(err)

	}
}

func removeLines(fn string, start, n int) (err error) {
	skipLines := func(b []byte, n int) ([]byte, bool) {
		for ; n > 0; n-- {
			if len(b) == 0 {
				return nil, false
			}
			x := bytes.IndexByte(b, '\n')
			if x < 0 {
				x = len(b)
			} else {
				x++
			}
			b = b[x:]
		}
		return b, true
	}
	if start < 1 {
		return errors.New("invalid request.  line numbers start at 1.")
	}
	if n < 0 {
		return errors.New("invalid request.  negative number to remove.")
	}
	var f *os.File
	if f, err = os.OpenFile(fn, os.O_RDWR, 0); err != nil {
		return
	}
	defer func() {
		if cErr := f.Close(); err == nil {
			err = cErr
		}
	}()
	var b []byte
	if b, err = ioutil.ReadAll(f); err != nil {
		return
	}
	cut, ok := skipLines(b, start-1)
	if !ok {
		return fmt.Errorf("less than %d lines", start)
	}
	if n == 0 {
		return nil
	}
	tail, ok := skipLines(cut, n)
	if !ok {
		return fmt.Errorf("less than %d lines after line %d", n, start)
	}
	t := int64(len(b) - len(cut))
	if err = f.Truncate(t); err != nil {
		return
	}
	if len(tail) > 0 {
		_, err = f.WriteAt(tail, t)
	}
	return
}

// in : 2007-03-04 12:20:00 1173010800
// out: 2007-03-04 12:20:00 1173010800
func formatTime(t time.Time, mark string) string {
	if mark == "" {
		mark = "in"
	}
	return fmt.Sprintf("%-4s %d-%02d-%02d %02d:%02d:%02d %d",
		mark+":",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second(),
		t.Unix())
}

func normalizeMark(mark string) string {
	r := mark
	switch mark {
	case "inn":
		r = "in"
	case "ut":
		r = "out"
	case "out":
		r = "out"
	case "in":
		r = "in"
	}
	return r
}

func determineNextStamp(mark string) string {
	var r string
	switch mark {
	case "in":
		r = "out"
	case "out":
		r = "in"
	case "":
		r = "in"
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
	var out *os.File
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
	lastIndex := int(stat.Size() - 1)

	if lastIndex <= 0 {
		return "" //  first Record is always in? ... maybe
	}
	if lastIndex >= 2 {
		for lastLn = lastIndex; lastLn > -1; {
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
