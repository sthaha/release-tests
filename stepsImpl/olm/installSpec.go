package olm

import (
	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	"github.com/openshift-pipelines/release-tests/pkg/helper"
	"github.com/openshift-pipelines/release-tests/pkg/operator"
	"github.com/openshift-pipelines/release-tests/stepsImpl/flags"
)

var _ = gauge.Step("Wait for Cluster CR availability", func() {
	helper.WaitForClusterCR(flags.Clients, config.ClusterCRName)
})

var _ = gauge.Step("Validate SCC", func() {
	operator.ValidateSCC(flags.Clients)
})

var _ = gauge.Step("Validate pipelines deployment into target namespace (openshift-pipelines)", func() {
	operator.ValidatePipelineDeployments(flags.Clients)
})

var _ = gauge.Step("Validate pipeline version <version>", func(version string) {
	operator.VerifyPipelineVersion(flags.Clients, version)
})

var _ = gauge.Step("Validate Triggers deployment into target namespace (openshift-pipelines)", func() {
	operator.ValidateTriggerDeployments(flags.Clients)
})

var _ = gauge.Step("Validate opeartor setup status", func() {
	operator.ValidateOperatorInstalledStatus(flags.Clients)
})
