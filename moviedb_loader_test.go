package moviedb_loader_test

import (
	. "localflix-server-"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type DatabaseFake interface{}

var _ = Describe("MoviedbLoader", func() {
	Describe("LoadTmdb", func() {
		It("should not return an error", func() {
			err := LoadTmdb()
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
