package chroot_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/rhdedgar/container-info/chroot"
)

var _ = Describe("Chroot", func() {
	var chrootPath = "/"

	Context("Validate chroot can return without errors", func() {
		It("Should chroot to the providied directory", func() {
			_, err := ChrootPath(chrootPath)
			Expect(err).To(BeNil())
		})
	})
})
