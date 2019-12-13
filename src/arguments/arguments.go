package arguments

import "time"

type Arguments struct {
	Stamp time.Time
	Mark string
	OutPath string
	DoCount bool
	SumPerDay bool
}
