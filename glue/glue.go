package glue

import (
	"fmt"
	"log"
	"os"
)

func Debug(fmt string, arg ...interface{}) {
	log.New(os.Stderr, "", 0).
		Printf(fmt, arg...)
}

func Die(e error) {
	if e != nil {
		panic(e)
	}
}

var tfn = 0
var base = "/tmp/tt-stamps-test"
// parallel tests need individual output files
func TestStampFile() string {
	tfn++
	f := fmt.Sprintf("%s.%d", base, tfn)
	_ = os.Remove(f)
	return f
}
func WipeTestFiles() {
	_ = os.Remove(base)
	for i := 0; i < tfn ; i++ {
		_ = os.Remove(fmt.Sprintf("%s.%d", base, i))

	}
	tfn = 0
}
