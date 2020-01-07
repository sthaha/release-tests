package cli

import (
	"log"
	"time"

	"gotest.tools/v3/icmd"
)

func RunQuiet(cmd ...string) *icmd.Result {
	return icmd.RunCmd(icmd.Cmd{Command: cmd, Timeout: 10 * time.Minute})
}

// Run runs a command and prints the output
func Run(cmd ...string) *icmd.Result {
	res := RunQuiet(cmd...)
	log.Printf("Output Stream:\n%s\n-----------------------\n", res.Stdout())
	log.Printf("Error Stream:\n%s\n-----------------------\n", res.Stderr())
	return res
}
