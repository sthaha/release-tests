package olm

import (
	"testing"

	"github.com/openshift-pipelines/release-tests/pkg/operator"
	. "github.com/smartystreets/goconvey/convey"
)

func TestFreshInstall(t *testing.T) {

	//defer operator.DeleteOperator(t, clients)

	Convey("Given a new cluster", t, func() {
		Convey("When I subscribe to the Pipelines Operator", func() {
			clients, _, cleanup := operator.Subscribe(t)
			defer cleanup()

			Convey("validate Cluster CR", func() {
				operator.ValidateClusterCR(t, clients)
				So(true, ShouldEqual, true)
			})

			Convey("validate SCC", func() {
				operator.ValidateSCC(t, clients)
				So(true, ShouldEqual, true)
			})

			Convey("installs Pipelines 0.9", func() {
				Convey("Validate pipelines deployment into target namespace (openshift-pipelines)", func() {
					operator.ValidatePipelineDeployments(t, clients)
					So(true, ShouldEqual, true)
				})

				Convey("Validate pipelines version", func() {
					operator.VerifyPipelineVersion(t, clients, "v0.9")
				})
			})

			Convey("installs Triggers 0.1", func() {
				operator.ValidateTriggerDeployments(t, clients)
				So(true, ShouldEqual, true)
			})

			SkipConvey("installs the following cluster tasks", func() {
				Convey("s2i", func() {
				})
				Convey("s2i-java-8", func() {
				})
				Convey("s2i-java-11", func() {
				})
				Convey("s2i-python-2", func() {
				})
				Convey("s2i-python-3", func() {
				})
				Convey("openshift-client", func() {
				})
				So(true, ShouldEqual, true)
			})
		})
	})
}
