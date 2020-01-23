package pipeline

import (
	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/openshift-pipelines/release-tests/pkg/pipelines"
	"github.com/openshift-pipelines/release-tests/stepsImpl/flags"
)

var _ = gauge.Step("Create output pipeline", func() {
	pipelines.CreateSamplePipeline(flags.Clients, flags.Namespace)
})

var _ = gauge.Step("Run pipeline", func() {
	pipelines.RunSamplePipeline(flags.Clients, flags.Namespace)
})

var _ = gauge.Step("Validate pipelinerun for success status", func() {
	pipelines.ValidatePipelineRunStatus(flags.Clients, flags.Namespace)
})
