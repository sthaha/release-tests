package pipeline

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPipelinesInstall(t *testing.T) {
	Convey("Given that the Operator is installed", t, func() {
		Convey("It should have installed pipelines controller", nil)
		Convey("It should have installed pipelines webhook", nil)
		Convey("It should have configured privileged SCC", nil)
		Convey("The Pipeline version should reflect in config CR", nil)
	})
}
