package operator

import (
	"log"
	"testing"

	"github.com/openshift-pipelines/release-tests/pkg/cli"
	"github.com/openshift-pipelines/release-tests/pkg/client"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	"github.com/openshift-pipelines/release-tests/pkg/helper"
	"gotest.tools/v3/icmd"
)

func DeleteClusterCR(t *testing.T) {
	t.Helper()

	res := cli.Run("oc", "delete", "config.operator.tekton.dev", "cluster")
	res.Assert(t, icmd.Expected{ExitCode: 0, Err: icmd.None})
}

func DeleteInstallPlan(t *testing.T) {
	t.Helper()

	installPlan := cli.RunQuiet(
		"oc", "get", "-n", "openshift-operators",
		"subscripton", "openshift-pipelines-operator",
		`-o=jsonpath={.status.installplan.name}`,
	).Stdout()

	res := cli.Run("oc", "delete", "-n", "openshift-operators", "installplan", installPlan)
	res.Assert(t, icmd.Expected{ExitCode: 0, Err: icmd.None})
	log.Printf("Deleted install plan %s\n", installPlan)
}

func DeleteSubscription(t *testing.T) {
	res := cli.Run(
		"oc", "delete", "-n", "openshift-operators",
		"subscription", "openshift-pipelines-operator",
	)
	res.Assert(t, icmd.Expected{ExitCode: 0, Err: icmd.None})
	log.Printf("Deleted Subscription %s\n", res.Stdout())
}

func DeleteOperator(t *testing.T, cs *client.Clients) {

	cr := helper.WaitForClusterCR(t, cs, config.ClusterCRName)

	helper.DeleteClusterCR(t, cs, config.ClusterCRName)

	ns := cr.Spec.TargetNamespace
	helper.ValidateDeploymentDeletion(t, cs,
		ns,
		config.PipelineControllerName,
		config.PipelineWebhookName,
		config.TriggerControllerName,
		config.TriggerWebhookName,
	)

	helper.ValidateSCCRemoved(t, cs, ns, config.PipelineControllerName)
	DeleteInstallPlan(t)
	DeleteSubscription(t)
}
