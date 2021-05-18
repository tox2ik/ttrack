package model

import (
	"testing"
)

func TestTuple_IsValid(t *testing.T) {

	bad := [][]Record{
		{

			{Mark: "in:", Day: "", Time: "", Stamp: 1},
			{Mark: "out", Day: "", Time: "", Stamp: 1},
		},
	}

	for i, rr := range bad {

		tuple := Tuple{
			Day:     rr[0].Day,
			Seconds: rr[1].Stamp - rr[0].Stamp,
			In:      rr[0],
			Out:     rr[1],
		}

		if tuple.IsValid() {
			t.Errorf("All tuples are bad in this test!\n(%d) %s", i, tuple)

		}


	}

}
