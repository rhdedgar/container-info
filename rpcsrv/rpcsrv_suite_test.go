package rpcsrv_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRpcsrv(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rpcsrv Suite")
}
