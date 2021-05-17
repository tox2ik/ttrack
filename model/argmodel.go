package model

import (
	"strings"
	"time"
)

type Arguments struct {
	DoDry         bool   `short:"n" long:"dry"     description:"Dry run"`
	DoInteractive bool   `short:"i" long:"tui"     description:"Start interactive editor (edit, i)"`
	DoCount       bool   `short:"c" long:"count"   description:"Count stamps             (count, c)"`
	DoSumPerDay   bool   `short:"s" long:"sum"     description:"Count average per day    (sum, s, cc)"`
	DoLog         bool   `short:"l" long:"log"     description:"Describe last time stamp (log)"`
	DoMark        bool   `short:"m" long:"mark"    description:"Sign out and back in     (mark)"`
	Mark          string `short:"r" long:"record"  description:"Sign in or out           (in, out)" choice:"in" choice:"out"`
	InStamp       string `short:"d" long:"date"    description:"Time to record           («argument»)"`
	Stamp         time.Time
	OutPath       string
}

// The default action is to add a stamp (in or out)
func (a Arguments) DoAddStamp() bool {

	return ! (a.DoInteractive || a.DoCount || a.DoLog) && ! a.Stamp.IsZero()
}

func (a Arguments) DoCountCurrentMonth() bool {
	return a.DoCount && len(a.OutPath) == 0

}

func (a *Arguments) SetFlagFromAlias(head string) bool {
	switch head {
	case "i": fallthrough
	case "edit":
		a.DoInteractive = true
	case "count", "c":
		a.DoCount = true
	case "sum", "s", "cc":
		a.DoCount = true
		a.DoSumPerDay = true
	case "mark":
		a.DoMark = true
	case "log":
		a.DoLog = true
	case "in", "out":
		a.Mark = head
	default:
		return false
	}
	return true

}

func (a Arguments) IsStamp(head string) bool {
	// : should be regex [0-9]:[0-9]
	cons := strings.Contains
	return ! cons(head, "/") && (cons(head, ":") ||
		cons(head, "now") ||
		cons(head, "hour") ||
		cons(head, "hours") ||
		cons(head, "second") ||
		cons(head, "seconds") ||
		cons(head, "tomorrow") ||
		cons(head, "yesterday") ||
		cons(head, "minute") ||
		cons(head, "years") ||
		cons(head, "year") ||
		cons(head, "days") ||
		cons(head, "day"))

}
