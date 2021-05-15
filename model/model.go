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

// Record in and record out
type Tuple struct {
	Day     string
	Seconds float32
	In      Record
	Out     Record
}

func (t Tuple) IsValid() bool { return (t.In.IsIn() || t.Out.IsOut()) && t.Seconds > 0 }

func (t Tuple) Format() string {
	return fmt.Sprintf("%10s %5.2f", t.Day, t.Seconds/3600)
}

func (t Tuple) String() string {
	return fmt.Sprintf("%s - %s (%.2fh)", t.In, t.Out, float64(t.Out.Stamp - t.In.Stamp) / 3600)
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
	var total float32
	for _, t := range tt.Items {
		_, err = writer.WriteString(t.Format())
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
	perDay := make(map[string]float32)
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
		_, err = writer.WriteString(fmt.Sprintf("%10s %5.2f\n", day, perDay[day]/3600))
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

func (tt Tuples) Hours() float32 { return tt.Seconds() / 3600 }

func (tt Tuples) Seconds() float32 {
	total := float32(0)
	for _, t := range tt.Items {
		if t.IsValid() {
			total += t.Seconds
		}
	}
	return total
}

func (tt Tuples) HoursAverage() float32 {
	validTuplesCount := float32(0)
	for _, t := range tt.Items {
		if t.IsValid() {
			validTuplesCount++
		}
	}
	return tt.Seconds() / 3600 / validTuplesCount
}

// todo: move isvalid check to an AddItem function on Tuples
