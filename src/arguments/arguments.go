package arguments

import "time"

type Arguments struct {
	Stamp time.Time
	Mark string
	OutPath string
	SumPerDay bool
	DoCount bool
	DoMark bool
	DoLog bool
}
