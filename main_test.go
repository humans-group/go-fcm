package fcm

import (
	"testing"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
)

func TestHttpClient(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("GoFcm.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "GoFcm", []Reporter{junitReporter})
}
