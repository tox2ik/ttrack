package glue

import (
	"fmt"
	"log"
	"os"
	"os/exec"
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


func IsExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}

	return false
}

func RunEditor(path string) error {
	cmd := exec.Command(os.Getenv("EDITOR"), path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
