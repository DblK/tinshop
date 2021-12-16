package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	tin "github.com/dblk/tinshop"
)

var _ = Describe("Tinshop", func() {
	Describe("AddNewGames", func() {
		BeforeEach(func() {
			tin.Games = make(map[string]interface{})
			tin.Games["success"] = "Welcome to your own shop!"
			tin.Games["titledb"] = make(map[string]interface{})
			tin.Games["files"] = make([]interface{}, 0)
		})

		It("Add empty table", func() {
			var newGameFiles []tin.FileDesc

			tin.AddNewGames(newGameFiles)
			Expect(tin.Games["files"]).To(HaveLen(0))
		})
	})
})
