package udp_test

import (
	"context"
	"golang-udp-server/udp"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UDP Client", func() {
	var client *udp.Client
	var server *udp.Server
	var err error
	var response string

	BeforeEach(func() {
		By("Creating UDP client")
		ctx := context.Background()
		client = udp.NewClient(ctx, "127.0.0.1:8080")
		By("Creating UDP Server")
		server = udp.NewServer(ctx, "0.0.0.0:8080")
		go server.Start()
		err, response = client.Request("/test", "hello world!")
	})

	Context("Make simple UDP request: Path: /test Body: hello World", func() {
		It("No errors", func() {
			Expect(err == nil).To(Equal(true))
		})

		It("Should not be empty", func() {
			Expect(response != "").To(Equal(true))
		})

		It("Result should be expected value", func() {
			Expect(strings.Compare(response, "This is a random response") == 1).To(Equal(true))
		})
	})
})
