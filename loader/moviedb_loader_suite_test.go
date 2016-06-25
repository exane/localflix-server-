package loader_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestLocalflixServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Moviedb_loader Suite")
}
