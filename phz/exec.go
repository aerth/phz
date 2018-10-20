package phz

import (
	logg "log"
	"os"
	"os/exec"

	"github.com/google/shlex"
)

var log = logg.New(os.Stderr, "phz: ", 0)

var DefaultFuncMap = map[string]interface{}{
	"exec": execstring,
}

func execstring(s string) string {
	cmdline, err := shlex.Split(os.ExpandEnv(s))
	if err != nil {
		return "error 5"
	}
	return execslice(cmdline)
}

func execslice(cmdline []string) string {
	cmd := exec.Command(cmdline[0])
	cmd.Args = cmdline[0:]
	b, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
	}
	return string(b)
}
