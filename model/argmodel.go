package model

import (
	"time"
)

type Arguments struct {
	DoCount   bool `short:"c" long:"count" description:"Count stamps"`
	DoLog     bool `short:"l" long:"log" description:"Describe last time stamp"`
	DoMark    bool `short:"m" long:"mark" description:"Sign out and back in"`
	DoDry     bool `short:"n" long:"dry" description:"Dry run"`
	SumPerDay bool `short:"s" long:"sum" description:"Count average per day"`

	Stamp   time.Time `short:"d" long:"date" description:"Time to record"`
	OutPath string
	Mark    string `short:"r" long:"record" description:"Sign in or out" choice:"in" choice:"out"`
}
