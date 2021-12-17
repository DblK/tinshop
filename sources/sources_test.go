package sources_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/dblk/tinshop/sources"
)

var _ = Describe("Sources", func() {
	It("Return list of game files", func() {
		files := sources.GetFiles()

		Expect(len(files)).To(Equal(0))
	})
})
