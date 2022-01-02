package main_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Security", func() {
	Describe("TinfoilMiddleware", func() {
		Context("No security", func() {
			It("Dummy", func() {
				Expect(true).To(BeTrue())
			})
		})
	})
})
