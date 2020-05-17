package main


// Represents a point in Time.
//
//     in:  2020-05-15 10:02:06 1589529726
//     out: 2020-05-15 15:40:53 1589550053
type Record struct {
	Mark  string
	Day   string
	Time  string
	Stamp uint32
}

// Record in and record out
type Tuple struct {
	Day     string
	Seconds float32
	in      *Record
	out     *Record
}

type Tuples struct {
	Items []Tuple
	// total        float32
	// totalHours   float32
	// averageHours float32
}

func (r *Record) isIn() bool      { return r.Mark == "in" }
func (r *Record) isOut() bool     { return r.Mark == "out" }
func (t *Tuple) isValid() bool    { return t.in != nil && t.out != nil && t.Seconds >= 0 }
func (tt *Tuples) Hours() float32 { return tt.Seconds() / 3600 }

func (tt *Tuples) Seconds() float32 {
	total := float32(0)
	for _, t := range tt.Items {
		if t.isValid() {
			total += t.Seconds
		}
	}
	return total
}

func (tt *Tuples) HoursAverage() float32 {
	validTuplesCount := float32(0)
	for _, t := range tt.Items {
		if t.isValid() {
			validTuplesCount++
		}
	}
	return tt.Seconds() / 3600 / validTuplesCount
}
