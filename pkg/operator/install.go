package operator

import (
	"log"
	"testing"

	"github.com/openshift-pipelines/release-tests/pkg/cli"
	"github.com/openshift-pipelines/release-tests/pkg/client"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	"github.com/openshift-pipelines/release-tests/pkg/helper"
	"gotest.tools/v3/icmd"

	. "github.com/smartystreets/goconvey/convey"
	op "github.com/tektoncd/operator/pkg/apis/operator/v1alpha1"
	"github.com/tektoncd/pipeline/pkg/names"
	knativetest "knative.dev/pkg/test"
)

func ValidateClusterCR(t *testing.T, cs *client.Clients) *op.Config {
	return helper.WaitForClusterCR(t, cs, config.ClusterCRName)
}

func VerifyPipelineVersion(t *testing.T, cs *client.Clients, version string) {
	cr := ValidateClusterCR(t, cs)
	So(cr.Status.Conditions[0].Version, ShouldStartWith, version)
}

func ValidateSCC(t *testing.T, cs *client.Clients) {
	cr := ValidateClusterCR(t, cs)
	helper.ValidateSCCAdded(t, cs, cr.Spec.TargetNamespace, config.PipelineControllerName)
}

func ValidatePipelineDeployments(t *testing.T, cs *client.Clients) {
	cr := ValidateClusterCR(t, cs)
	helper.ValidateDeployments(t, cs, cr.Spec.TargetNamespace,
		config.PipelineControllerName, config.PipelineWebhookName)
}
func ValidateTriggerDeployments(t *testing.T, cs *client.Clients) {
	cr := ValidateClusterCR(t, cs)
	helper.ValidateDeployments(t, cs, cr.Spec.TargetNamespace,
		config.TriggerControllerName, config.TriggerWebhookName)
}

func ValidateOperatorInstall(t *testing.T, cs *client.Clients) {
	log.Printf("Waiting for operator to be up and running....\n")

	ValidatePipelineDeployments(t, cs)
	ValidateTriggerDeployments(t, cs)

	// Refresh Cluster CR
	cr := ValidateClusterCR(t, cs)

	if code := cr.Status.Conditions[0].Code; code != op.InstalledStatus {
		t.Errorf("Expected code to be %s but got %s", op.InstalledStatus, code)
	}
	log.Printf("Operator is up\n")

}

func NewClientset() (*client.Clients, string, func()) {
	ns := names.SimpleNameGenerator.RestrictLengthWithRandomSuffix("testrelease")
	cs := client.NewClients(knativetest.Flags.Kubeconfig, knativetest.Flags.Cluster, ns)
	helper.CreateNamespace(cs.KubeClient, ns)

	cleanup := func() { helper.DeleteNamespace(cs.KubeClient, ns) }
	return cs, ns, cleanup
}

func Subscribe(t *testing.T, subsPath string) (*client.Clients, string, func()) {
	t.Helper()

	cli.Run("pwd")
	res := cli.Run("oc", "apply", "-f", subsPath)
	res.Assert(t, icmd.Expected{ExitCode: 0, Err: icmd.None})

	cs, ns, cleanupNs := NewClientset()

	ValidateOperatorInstall(t, cs)
	helper.VerifyServiceAccountExists(cs.KubeClient, ns)

	cleanup := func() {
		DeleteSubscription(t)
		cleanupNs()
	}

	return cs, ns, cleanup
}
