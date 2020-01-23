package stepsImpl

import (
	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/openshift-pipelines/release-tests/pkg/helper"
	"github.com/openshift-pipelines/release-tests/pkg/operator"
	"github.com/openshift-pipelines/release-tests/stepsImpl/flags"
)

var _ = gauge.Step("Create random namespace and clientset", func() {
	flags.Clients, flags.Namespace, flags.Cleanup = helper.NewClientSet()
})

var _ = gauge.Step("Delete namespace and clientset", func() {
	defer flags.Cleanup()
})

var _ = gauge.Step("Operator should be installed", func() {
	operator.ValidateOperatorInstall(flags.Clients)
})
