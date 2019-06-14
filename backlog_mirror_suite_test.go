package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBacklogMirror(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "BacklogMirror Suite")
}
