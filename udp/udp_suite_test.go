package udp_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestUdp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Udp Suite")
}
