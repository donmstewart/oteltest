package oteltest_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestOteltest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Oteltest Suite")
}
