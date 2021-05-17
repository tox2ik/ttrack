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
	Stamp uint32
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

// Record in and record out
type Tuple struct {
	Day     string
	Seconds uint32
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
	return fmt.Sprintf("%s   %.5s - %.5s   %#2d:%02d\n",
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


func (tt Tuples) ReportHours(writer io.StringWriter) (err error) {
	var total uint32
	for _, t := range tt.Items {
		_, err = writer.WriteString(t.FormatDur()+"\n")
		if err != nil {
			return
		}
		total += t.Seconds
	}
	_, err = writer.WriteString("" +
		fmt.Sprintf("    total: %5.2f\n", tt.Hours()) +
		fmt.Sprintf("  average: %5.2f\n", tt.HoursAverage()))
	return
}

func (tt Tuples) ReportHoursPerDay(writer io.StringWriter) (err error) {
	perDay := make(map[string]uint32)
	for _, t := range tt.Items {
		perDay[t.Day] += t.Seconds
	}

	i := 0
	keys := make([]string, len(perDay))
	for key := range perDay {
		keys[i] = key
		i++
	}
	sort.Strings(keys)

	for _, day := range keys {
		_, err = writer.WriteString(fmt.Sprintf("%10s %5.2f\n", day, float32(perDay[day])/3600))
		if err != nil {
			return
		}
	}
	itemc := float32(len(perDay))
	if itemc == 0 {
		itemc = 100000000
	}
	_, err = writer.WriteString("" +
		fmt.Sprintf("    total: %5.2f\n", tt.Hours()) +
		fmt.Sprintf("  average: %5.2f\n", tt.Hours()/itemc))
	return
}

func (tt Tuples) Seconds() int {
	total := uint32(0)
	for _, t := range tt.Items {
		if t.IsValid() {
			total += t.Seconds
		}
	}
	return int(total)
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
	return fmt.Sprintf("%#2d:%02d", h, m)
}
func (tt Tuples) HoursAverageH() string {
	s := tt.Seconds() / tt.validCount()
	h := s / 3600
	m := s % 3600 / 60
	return fmt.Sprintf("%#2d:%02d", h, m)

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
