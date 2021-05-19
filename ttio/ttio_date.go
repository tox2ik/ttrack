package ttio

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func ParseDate(in string) (time.Time, error) {

	if t, e := parseGo(in); e == nil {
		return t, nil
	}
	if t, e := parseGnuDate(in); e == nil {
		return t, nil
	}
	return time.Time{}, fmt.Errorf("failed to parse time: %s", in)
}

func parseGo(inputDate string) (time.Time, error) {
	inputDate = strings.TrimSpace(inputDate)
	// ref-time: Mon Jan 2 15:04:05 -0700 MST 2006
	layoFull := []string{
		"2006-01-02 15:04",
		"2006-01-02 15:04:05",
	}
	layoShort := []string{
		"15:04:05",
		"15:04",
		"1504",
		"15",
		"today 15",
		"today 1504",
		"today 15:04",
		"today 15:04:05",
		"yest 15",
		"yest 1504",
		"yest 15:04",
		"yest 15:04:05",
		"yesterday 15",
		"yesterday 1504",
		"yesterday 15:04",
		"yesterday 15:04:05",
	}
	for _, s := range layoShort {
		t, err := time.Parse(s, inputDate)
		if nil == err {
			if strings.Contains(inputDate, "today") {
				return todayTime(t, 0), nil
			} else
			if strings.Contains(inputDate, "yest") {
				return todayTime(t, -1), nil
			} else {
				return todayTime(t, 0), nil
			}
		}
	}
	for _, s := range layoFull {
		t, err := time.Parse(s, inputDate)
		if nil == err {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unable to parse date")
}

func todayTime(parsed time.Time, dayOffset int) time.Time {
	now := time.Now()
	return time.Date(
		now.Year(),
		now.Month(),
		now.Day()+dayOffset,
		parsed.Hour(),
		parsed.Minute(),
		parsed.Second(),
		0,
		now.Location())
}

func parseGnuDate(inputDate string) (time.Time, error) {
	var t time.Time
	var err error
	var out []byte

	comment := `
maybe-todo: handle schmuck-os date and winders.
The semantics of GNU 'date -d' are very useful.
For more info read 'info date'; section 29.7 Relative Items in date strings
https://www.gnu.org/software/coreutils/manual/html_node/Relative-items-in-date-strings.html#Relative-items-in-date-strings
The intro-quote of section 29 Date input formats is also worth a read.
`
	comment += ""

	date := exec.Command("date", "--rfc-email", "-d", inputDate)
	//date.Env = []string{ "TZ=UTC"}

	out, err = date.Output()
	t, err = time.Parse(time.RFC1123Z, strings.Trim(string(out), "\r\n"))
	if nil == err {
		return t, nil
	}
	return time.Time{}, err
}
