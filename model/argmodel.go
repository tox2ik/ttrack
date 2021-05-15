package model

import (
	"time"
)

type Arguments struct {
	DoDry     bool   `short:"n" long:"dry"     description:"Dry run"`
	DoCount   bool   `short:"c" long:"count"   description:"Count stamps             (count, c)"`
	DoLog     bool   `short:"l" long:"log"     description:"Describe last time stamp (log)"`
	DoMark    bool   `short:"m" long:"mark"    description:"Sign out and back in     (mark)"`
	SumPerDay bool   `short:"s" long:"sum"     description:"Count average per day    (sum, cc)"`
	Mark      string `short:"r" long:"record"  description:"Sign in or out           (in, out)" choice:"in" choice:"out"`
	InStamp   string `short:"d" long:"date"    description:"Time to record           («argument»)"`
	Stamp     time.Time
	OutPath   string
}

// The default action is to add a stamp (in or out)
func (a Arguments) DoAddStamp() bool {

	return ! (a.DoCount || a.DoLog) && ! a.Stamp.IsZero()
}
