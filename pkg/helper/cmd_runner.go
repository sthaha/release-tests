package helper

import (
	"testing"

	"github.com/openshift-pipelines/release-tests/pkg/config"
	"gotest.tools/v3/icmd"
)

var t = &testing.T{}

func RunQuiet(cmd ...string) *icmd.Result {
	return icmd.RunCmd(icmd.Cmd{Command: cmd, Timeout: config.Timeout})
}

// CmdShouldPass runs a command and verfies exit code (0)
func CmdShouldPass(cmd ...string) *icmd.Result {
	res := RunQuiet(cmd...)
	res.Assert(t, icmd.Expected{ExitCode: 0, Err: icmd.None})
	return res
}

// CmdShouldFail runs a command and verifies exit code ~(0)
func CmdShouldFail(cmd ...string) *icmd.Result {
	res := RunQuiet(cmd...)
	res.Assert(t, icmd.Expected{ExitCode: 1})
	return res
}
