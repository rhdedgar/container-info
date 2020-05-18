package chroot_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestChroot(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Chroot Suite")
}
