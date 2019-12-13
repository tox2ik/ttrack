package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)


/*

function uniq_dates(dates) {
	c=0; for (i in dates) c++;
	return c;
}

$2 ~ /^[0-9]{4}-[0-9]{2}-[0-9]{2}$/ { dates[$2]++ }

/^\s*inn:/ { inn=$NF }
/^\s*ut:/ { ut=$NF ;
	timer = (ut-inn) / 3600
	total +=timer
	avg = total / uniq_dates(dates)
	printf "%10s %5.2f\n", $2, timer
}

END {
	printf "%10s %5.2f\n  average: %5.2f\n", "total:", total, avg

}
*/

type record struct {
	mark string
	day string
	time string
	stamp int
}

type tuple struct{
	day string
	seconds float32
}

type tuples struct {
	items []tuple
	total float32
	totalHours float32
	averageHours float32
}

func FileAsArray(file *os.File) []string {
	lines := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func ParseRecords(file *os.File) ([]record, tuples, error) {
	file.Seek(0, io.SeekStart)
	scanner := bufio.NewScanner(file)
	records := make([]record, 0)
	for scanner.Scan() {
		line := strings.Trim(strings.ReplaceAll(scanner.Text(), "  ", " "), " " )
		fields := strings.Split(line, " ")
		if len(fields) >= 4 {
			ts, _ := strconv.Atoi(fields[3])
			rec := record{
				mark: strings.Replace(fields[0], ":","", 1),
				day: fields[1],
				time: fields[2],
				stamp: ts}
			records = append(records, rec)
		}
	}
	var in int
	var day string
	tuplesSlice := make([]tuple, 0)
	for i, rec := range records {
		if "in" == rec.mark {
			day = rec.day
			in = rec.stamp
		} else {
			if day != rec.day { log.New(os.Stderr, "", 0).Printf("Day mismatch for record: %d (%s,%s)", i, day, rec.day) }
			tuplesSlice = append(tuplesSlice, tuple{day: day, seconds: float32(rec.stamp - in) })
		}
	}
	total := float32(0)
	for _, t := range tuplesSlice {
		total += t.seconds
	}
	tupleSummary := tuples{
		items: tuplesSlice,
		total: total,
		totalHours: total / 3600,
		averageHours: total / 3600 / float32(len(tuplesSlice))}

	if len(records) % 2 != 0 {
		return records, tupleSummary, errors.New("file contains unfinished stamps")
	}
	return records, tupleSummary, nil
}

func FormatTuple(t tuple) string {
	return fmt.Sprintf("%10s %5.2f", t.day, t.seconds/3600)
}

func Count(file *os.File) error {
	var total float32
	_, tuples, parseError := ParseRecords(file)
	for _, t := range tuples.items {
		fmt.Println(FormatTuple(t))
		total += t.seconds
	}
	fmt.Printf("    total: %5.2f\n", tuples.totalHours)
	fmt.Printf("  average: %5.2f\n", tuples.averageHours)
	return parseError
}

func CountPerDay(file *os.File) error {
	//var total float32
	secondsPerDay := make(map[string]float32)
	_, tuples, parseError := ParseRecords(file)
	for _, t := range tuples.items {
		secondsPerDay[t.day] += t.seconds
		//total += t.seconds
	}
	keys := make([]string, 0, len(secondsPerDay))
	for key := range secondsPerDay { keys = append(keys, key) }
	sort.Strings(keys)
	for _, day := range keys {
		fmt.Printf("%10s %5.2f\n", day, secondsPerDay[day] / 3600)
	}
	fmt.Printf("    total: %5.2f\n", tuples.totalHours )
	fmt.Printf("  average: %5.2f\n", tuples.total / 3600 / float32(len(secondsPerDay)))
	return parseError
}


