package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
)

func debug(fmt string, arg ...interface{}) {
	log.New(os.Stderr, "", 0).
		Printf(fmt, arg...)
}


func FileAsArray(file *os.File) []string {
	lines := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}


func FormatTuple(t Tuple) string {
	return fmt.Sprintf("%10s %5.2f", t.Day, t.Seconds/3600)
}

func Count(file *os.File) error {
	var total float32
	_, tuples, parseError := ParseRecords(file)

	if parseError == nil && len(tuples.Items) > 0 {
		for _, t := range tuples.Items {
			fmt.Println(FormatTuple(t))
			total += t.Seconds
		}

		fmt.Printf("    total: %5.2f\n", tuples.Hours())
		fmt.Printf("  average: %5.2f\n", tuples.HoursAverage())
	}
	return parseError
}

func CountPerDay(file *os.File) error {
	secondsPerDay := make(map[string]float32)
	var tuples Tuples
	var parseError error
	if _, tuples, parseError = ParseRecords(file); parseError != nil {
		return parseError
	}
	for _, t := range tuples.Items {
		secondsPerDay[t.Day] += t.Seconds
	}
	keys := make([]string, 0, len(secondsPerDay))
	for key := range secondsPerDay {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, day := range keys {
		fmt.Printf("%10s %5.2f\n", day, secondsPerDay[day]/3600)
	}

	fmt.Printf("    total: %5.2f\n", tuples.Hours())
	fmt.Printf("  average: %5.2f\n", tuples.Hours()/float32(len(secondsPerDay)))
	return parseError
}
