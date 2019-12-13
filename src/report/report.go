package report

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	A "../arguments"
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


func Count(file *os.File, args A.Arguments) {
	defer file.Close()
	scanner := bufio.NewScanner(file)
	dates := make(map[string]int)
	hoursPerDay := make(map[string]float32)

	var hours float32
	var totalHours float32
	var average float32
	var in int
	var out int

	var tuples int
	var entries int

	for scanner.Scan() {
		line := strings.ReplaceAll(scanner.Text(), "  ", " ")
		if len(line) == 0 {
			continue
		}
		fields := strings.Split(line, " ")
		mark := fields[0]
		date := fields[1]
		epoch := fields[3]
		dates[date]++

		if "in:" == mark {
			in, _ = strconv.Atoi(epoch)
		} else if "out:" == mark{
			out, _ = strconv.Atoi(epoch)
			hours += float32(out-in)/3600
			totalHours += hours
			hoursPerDay[date] += hours
			tuples++
			if ! args.SumPerDay {
				fmt.Printf("%10s %5.2f\n", date, hours)
			}
		}
	}
	for _, e := range dates { entries += e }

	average = totalHours / float32(tuples)

	if ! args.SumPerDay {
		fmt.Printf("    total: %5.2f\n", totalHours)
		fmt.Printf("  average: %5.2f\n", average)
	} else {
		for day, hours := range hoursPerDay {
			fmt.Printf("%10s %5.2f\n", day, hours)
		}
		fmt.Printf("    total: %5.2f\n", totalHours)
		fmt.Printf("  average: %5.2f\n", totalHours / float32(len(hoursPerDay)))
	}

	if entries %2 != 0 {
		log.New(os.Stderr, "", 0).Println("File contains open stamps.")
	}
}
