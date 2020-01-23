package pipeline

import (
	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/openshift-pipelines/release-tests/pkg/pipelines"
	"github.com/openshift-pipelines/release-tests/stepsImpl/flags"
)

var _ = gauge.Step("Create Task", func() {
	pipelines.CreateTask(flags.Clients, flags.Namespace)
})

var _ = gauge.Step("Run Task with <serviceAccount> SA", func(serviceAccount string) {
	pipelines.CreateTaskRunWithSA(flags.Clients, flags.Namespace, serviceAccount)
})

var _ = gauge.Step("Validate TaskRun for failed status", func() {
	pipelines.ValidateTaskRunForFailedStatus(flags.Clients, flags.Namespace)
})

//=======================================================//

var _ = gauge.Step("Create pipeline", func() {
	pipelines.CreatePipeline(flags.Clients, flags.Namespace)
})

var _ = gauge.Step("Run pipeline with <serviceAccount> SA", func(serviceAccount string) {
	pipelines.CreatePipelineRunWithSA(flags.Clients, flags.Namespace, serviceAccount)
})

var _ = gauge.Step("Validate pipelineRun for failed status", func() {
	pipelines.ValidatePipelineRunForFailedStatus(flags.Clients, flags.Namespace)
})
