package model

import (
	"fmt"
	"io"
	"sort"
)

// Represents a point in Time.
//
//     in:  2020-05-15 10:02:06 1589529726
//     out: 2020-05-15 15:40:53 1589550053
type Record struct {
	Mark  string
	Day   string
	Time  string
	Stamp int64
}

func (r Record) IsIn() bool   { return r.Mark == "in" }
func (r Record) IsOut() bool  { return r.Mark == "out" }
func (r Record) String() string { return fmt.Sprintf("%s %s %s %d", r.Mark, r.Day, r.Time, r.Stamp)}

func (r Record) IsValid() bool {
	return (r.Mark == "in" || r.Mark == "out") &&
		len(r.Day) == 10 &&
		len(r.Time) == 8 &&
		r.Stamp >= 1
}

func (r Record) Equals(rm Record) bool {
	return r.Day == rm.Day &&
		r.Mark == rm.Mark &&
		r.Stamp == rm.Stamp &&
		r.Time == rm.Time

}

func (r Record) FormatText() string {
	return fmt.Sprintf("%-4s %s %s %d\n", r.Mark+":", r.Day, r.Time, r.Stamp)
}

// Record in and record out
type Tuple struct {
	Day     string
	Seconds int64
	In      Record
	Out     Record
}

func (t Tuple) IsValid() bool { return (t.In.IsIn() || t.Out.IsOut()) && t.Seconds > 0 }

func (t Tuple) FormatDur() string {
	return fmt.Sprintf("%10s %5.2f", t.Day, float32(t.Seconds)/3600)
}
func (t Tuple) FormatHuman() string {
	h := t.Seconds / 3600
	m := t.Seconds % 3600 / 60
	return fmt.Sprintf("%s   %.5s - %.5s   %2d:%02d\n",
		t.In.Day,
		t.In.Time,
		t.Out.Time,
		h,
		m)
}

func (t Tuple) String() string {
	return fmt.Sprintf("%s - %s (%.2fh)", t.In, t.Out, float64(t.Out.Stamp - t.In.Stamp) / 3600)
}

func (t *Tuple) Swap() {
	in := t.In
	out := t.Out

	t.Out = in
	t.In = out
	t.Day = t.In.Day

	t.In.Mark = "in"
	t.Out.Mark = "out"

	if t.In.IsValid() {
		//t.In.Stamp = uint32(time.Now().Unix())
		if t.Out.Stamp > t.In.Stamp {
			t.Seconds = t.Out.Stamp - t.In.Stamp
		}

	}

	if !(t.In.IsValid() && t.Out.IsValid()) {
		t.Seconds = 0
	}

	if t.Out.Stamp < t.In.Stamp {
		t.Seconds = 0
	}


}

func (t Tuple) Len() int {
	c := 0
	if t.In.IsValid() {
		c++
	}
	if t.Out.IsValid() {
		c++
	}
	return c
}

type Tuples struct {
	Items []Tuple
	// total        float32
	// totalHours   float32
	// averageHours float32
}

func (tt Tuples) Last() Tuple {
	if len(tt.Items) == 0 {
		return Tuple{}
	}
	return tt.Items[len(tt.Items)-1]
}


func (tt Tuples) ReportHours(ow io.Writer) (err error) {
	var total int64
	for _, t := range tt.Items {
		_, err = fmt.Fprintf(ow, t.FormatDur()+"\n")
		if err != nil {
			return
		}
		total += t.Seconds
	}
	_, err = fmt.Fprintf(ow, "    total: %5.2f\n", tt.Hours())
	_, err = fmt.Fprintf(ow, "  average: %5.2f\n", tt.HoursAverage())
	return
}


func (tt Tuples) ReportHoursHuman(ow io.Writer) (err error) {
	for _, t := range tt.Items {
		fmt.Fprintf(ow, t.FormatHuman())
	}
	fmt.Fprintf(ow, "%19s   total: %s\n", " ", tt.HoursH())
	fmt.Fprintf(ow, "%19s average: %s\n", " ", tt.HoursAverageH())
	return

}

func (tt Tuples) ReportRecords(ow io.Writer) (err error) {
	for _, t := range tt.Items {
		_, err = fmt.Fprintf(ow, t.In.FormatText())
		if err != nil {
			return
		}
		_, err = fmt.Fprintf(ow, t.Out.FormatText())
		if err != nil {
			return
		}
	}
	return
}

func (tt Tuples) ReportHoursPerDay(ow io.Writer) (err error) {
	orderedDays, secPerDay := tt.SecondsPerDay()
	for _, ymd := range orderedDays {
		_, err = fmt.Fprintf(ow, "%10s %5.2f\n", ymd, float32(secPerDay[ymd])/3600)
		if err != nil {
			return
		}
	}
	itemc := float32(len(secPerDay))
	if itemc == 0 {
		itemc = 100000000
	}
	fmt.Fprintf(ow, "    total: %5.2f\n", tt.Hours())
	fmt.Fprintf(ow, "  average: %5.2f\n", tt.Hours()/itemc)
	return
}

func (tt Tuples) Seconds() int {
	total := int64(0)
	for _, t := range tt.Items {
		if t.IsValid() {
			total += t.Seconds
		}
	}
	return int(total)
}

func (tt Tuples) SecondsPerDay() (orderedDays []string, secPerDay map[string]int64) {
	secPerDay = make(map[string]int64)
	for _, t := range tt.Items {
		secPerDay[t.Day] += t.Seconds
	}

	i := 0
	orderedDays = make([]string, len(secPerDay))
	for ymd := range secPerDay {
		orderedDays[i] = ymd
		i++
	}
	sort.Strings(orderedDays)
	return
}

func (tt Tuples) Hours() float32 { return float32(tt.Seconds()) / 3600 }

func (tt Tuples) validCount() (v int) {
	for _, t := range tt.Items {
		if t.IsValid() {
			v++
		}
	}
	return
}

func (tt Tuples) HoursAverage() float32 {
	return float32(tt.Seconds()) / 3600 / float32(tt.validCount())
}

func (tt Tuples) HoursH() string {

	s := tt.Seconds()
	h := s / 3600
	m := s % 3600 / 60
	return fmt.Sprintf("%2d:%02d", h, m)
}
func (tt Tuples) HoursAverageH() string {
	s := tt.Seconds() / tt.validCount()
	h := s / 3600
	m := s % 3600 / 60
	return fmt.Sprintf("%2d:%02d", h, m)

}

func (tt *Tuples) Remove(cy int) Tuple {
	if cy < 0 || len(tt.Items) == 0 && cy == 0 || cy >= len(tt.Items) {
		return Tuple{}
	}

	rm := tt.Items[cy]
	ni := []Tuple{}
	ni = append(ni, tt.Items[0:cy]...)
	ni = append(ni, tt.Items[cy+1:]...)
	tt.Items = ni
	return rm


}


// todo: move isvalid check to an AddItem function on Tuples
