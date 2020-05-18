package rpcsrv_test

import (
	"context"
	"net/rpc"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/rhdedgar/container-info/rpcsrv"
)

var _ = Describe("Rpcsrv", func() {
	var sockPath = "@rpctest.sock"

	BeforeEach(func() {
		go RPCSrv(sockPath)
	})

	Context("Validate connection to the running RPC server", func() {
		It("Should be accepting connections", func() {
			_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			for {
				_, err := rpc.DialHTTP("unix", sockPath)
				if err == nil {
					Expect(err).To(BeNil())
					break
				}
			}
		})
	})
})
