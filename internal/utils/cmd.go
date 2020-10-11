package utils

import (
	"github.com/go-cmd/cmd"
)

// CMD executes a linux command
func CMD(name string, args ...string) (err error, stdout, stderr []string) {
	c := cmd.NewCmd(name, args...)
	s := <-c.Start()
	stdout = s.Stdout
	stderr = s.Stderr
	return
}
