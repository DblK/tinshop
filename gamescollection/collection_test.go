package gamescollection_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	collection "github.com/dblk/tinshop/gamescollection"
	"github.com/dblk/tinshop/repository"
)

var _ = Describe("Collection", func() {
	It("Return list of games", func() {
		games := collection.Games()

		Expect(games.Files).To(HaveLen(0))
	})
	Describe("AddNewGames", func() {
		BeforeEach(func() {
			collection.ResetGamesCollection()
		})
		It("Add an empty array", func() {
			newGames := make([]repository.FileDesc, 0)
			collection.AddNewGames(newGames)

			games := collection.Games()
			Expect(games.Files).To(HaveLen(0))
		})
		It("Add a game", func() {
			newGames := make([]repository.FileDesc, 0)
			newFile := repository.FileDesc{
				Size:     42,
				Path:     "here",
				GameID:   "0000000000000001",
				GameInfo: "[0000000000000001][v0].nsp",
				HostType: repository.LocalFile,
			}
			newGames = append(newGames, newFile)
			collection.AddNewGames(newGames)

			games := collection.Games()
			Expect(games.Files).To(HaveLen(1))
		})
	})
})
