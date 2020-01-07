package olm

import (
	"log"
	"testing"
	"time"

	"gotest.tools/icmd"
)

func subscriptionPath() string {
	return "config/subscription.yaml"
}

func SubscribeOperator(t *testing.T) {
	t.Helper()

	res := icmd.RunCmd(icmd.Cmd{
		Command: []string{"oc", "apply", "-f", subscriptionPath()},
		Timeout: 10 * time.Minute})

	log.Printf("%s\n", res.Stdout())
	res.Assert(t, icmd.Expected{ExitCode: 0, Err: icmd.None})
}
