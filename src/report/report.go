package report

import (
	"bufio"
	"errors"
	"fmt"
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
	day     string
	seconds float32
}


func parseRecords(file *os.File) ([]record, []tuple, error) {
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
	tuples := make([]tuple, 0)
	for i, rec := range records {
		if "in" == rec.mark {
			day = rec.day
			in = rec.stamp
		} else {
			if day != rec.day { log.New(os.Stderr, "", 0).Printf("Day mismatch for record: %d (%s,%s)", i, day, rec.day) }
			tuples = append(tuples, tuple{day: day, seconds: float32(rec.stamp - in) })
		}
	}

	if len(records) % 2 != 0 {
		return records, tuples, errors.New("file contains unfinished stamps")
	}
	return records, tuples, nil
}


func Count(file *os.File) error {
	var total float32
	_, tuples, parseError := parseRecords(file)
	for _, t := range tuples {
		fmt.Printf("%10s %5.2f\n", t.day, t.seconds/3600)
		total += t.seconds
	}
	fmt.Printf("    total: %5.2f\n", total / 3600 )
	fmt.Printf("  average: %5.2f\n", total / 3600 / float32(len(tuples)))
	return parseError
}

func CountPerDay(file *os.File) error {
	var total float32
	secondsPerDay := make(map[string]float32)
	_, tuples, parseError := parseRecords(file)
	for _, t := range tuples {
		secondsPerDay[t.day] += t.seconds
		total += t.seconds
	}
	keys := make([]string, 0, len(secondsPerDay))
	for key := range secondsPerDay { keys = append(keys, key) }
	sort.Strings(keys)
	for _, day := range keys {
		fmt.Printf("%10s %5.2f\n", day, secondsPerDay[day] / 3600)
	}
	fmt.Printf("    total: %5.2f\n", total / 3600 )
	fmt.Printf("  average: %5.2f\n", total / 3600 / float32(len(secondsPerDay)))
	return parseError
}

