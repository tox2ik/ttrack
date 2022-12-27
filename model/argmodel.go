package model

import (
	"strings"
	"time"
)

type Arguments struct {
	DoDry         bool   `short:"n" long:"dry"     description:"Dry run"`
	DoDump        bool   `short:"d" long:"dump"    description:"Dump records file        (dump)"`
	DoEdit        bool   `short:"e" long:"edit"    description:"Edit with external editor (edit, e)"`
	DoInteractive bool   `short:"i" long:"tui"     description:"Start interactive editor (tui, i)"`
	DoCount       bool   `short:"c" long:"count"   description:"Count stamps             (count, c)"`
	DoSumPerDay   bool   `short:"s" long:"sum"     description:"Count average per day    (sum, s, cc)"`
	DoLog         bool   `short:"l" long:"log"     description:"Describe last time stamp (log)"`
	DoList        bool   `short:"t" long:"list"    description:"List stamps and days     (list)"`
	DoMark        bool   `short:"m" long:"mark"    description:"Sign out and back in     (mark)"`
	Mark          string `short:"r" long:"record"  description:"Sign in or out           (in, out)" choice:"in" choice:"out"`
	InStamp       string `short:"a" long:"date"    description:"Time to record           («argument»)"`
	Stamp         time.Time
	OutPath       string
}

func (a *Arguments) SetFlagFromAlias(head string) bool {
	switch head {
	case "tui", "ui", "i":
		a.DoInteractive = true
	case "e":
		fallthrough
	case "edit":
		a.DoEdit = true
	case "count", "c":
		a.DoCount = true
	case "sum", "s", "cc":
		a.DoCount = true
		a.DoSumPerDay = true
	case "mark":
		a.DoMark = true

	case "log", "l":
		a.DoLog = true
	case "in", "out":
		a.Mark = head

	case "dump", "d", "ll", "l1":
		a.DoDump = true

	case "list", "ls", "t":
		a.DoList = true
	default:
		return false
	}
	return true

}

// The default action is to add a stamp (in or out)
func (a Arguments) DoAddStamp() bool {
	return !(a.DoInteractive || a.DoCount || a.DoLog) && !a.Stamp.IsZero()
}

func (a Arguments) DoCountCurrentMonth() bool {
	return a.DoCount && len(a.OutPath) == 0

}

func (a Arguments) IsStamp(head string) bool {
	// : should be regex [0-9]:[0-9]
	cons := strings.Contains
	return !cons(head, "/") && (cons(head, ":") ||
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

func (a *Arguments) ResetTime() {
	a.Stamp = time.Time{}
	a.InStamp = ""
}
