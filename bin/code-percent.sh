#!/bin/bash
total=$(find -name '*.go' -and -not -name '*_test.go' | xargs wc -l | awk '/total$/{ print $1}')
find -name '*.go' -and -not -name '*_test.go' | xargs wc -l | sed 's: \.::' | awk -F/ -v tot=$total '!/total$/{
	#printf "%.2f %s\n", $1/tot, $2 "/" $3 "/" $4 "/" $5 "/" $6
	dir[$2]+=$1
} END {
	for (i in dir) {
		printf "%.2f%% %s/\n", dir[i]/tot, i
	}
}' | sed s:/\\+\$::
