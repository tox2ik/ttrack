package glue

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
)

func Debug(fmt string, arg ...interface{}) {
	log.New(os.Stderr, "", 0).
		Printf(fmt, arg...)
}

func Die(e error) {

	if e != nil {
		pc, _, _, ok := runtime.Caller(1)
		if !ok {
			fmt.Fprintf(os.Stderr, "???fn()\n")
		} else {
			pcfn := runtime.FuncForPC(pc)

			fnName := path.Base(pcfn.Name())
			lastDot := strings.LastIndex(fnName, ".")
			fnShort := fnName[lastDot+1:]
			fpath, l := pcfn.FileLine(pc)
			base := path.Base(fpath)

			fmt.Fprintf(os.Stderr, "[%s:%d].%s()\n", base, l, fnShort)
		}

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
	for i := 0; i < tfn; i++ {
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
	ed := strings.Split(os.Getenv("EDITOR"), " ")
	args := []string{ path }
	if len(ed) > 1 {
		args = append(ed[1:], path)
	}
	cmd := exec.Command(ed[0], args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
