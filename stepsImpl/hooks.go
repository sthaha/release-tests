package stepsImpl

import (
	"github.com/getgauge-contrib/gauge-go/gauge"
	. "github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/olm"
	"github.com/openshift-pipelines/release-tests/stepsImpl/flags"
)

var _ = gauge.BeforeSuite(func() {
	flags.Clients, _, flags.CleanupSuite = olm.Subscribe(flags.OperatorVersion)
}, []string{}, AND)

var _ = gauge.AfterSuite(func() {
	defer flags.CleanupSuite()
}, []string{}, AND)
