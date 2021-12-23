package gamescollection_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	collection "github.com/DblK/tinshop/gamescollection"
	"github.com/DblK/tinshop/mock_repository"
	"github.com/DblK/tinshop/repository"
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
	Describe("RemoveGame", func() {
		BeforeEach(func() {
			collection.ResetGamesCollection()
		})
		It("Removing existing ID", func() {
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

			Expect(collection.Games().Files).To(HaveLen(1))
			collection.RemoveGame("0000000000000001")
			Expect(collection.Games().Files).To(HaveLen(0))
		})
		It("Removing not existing ID", func() {
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

			Expect(collection.Games().Files).To(HaveLen(1))
			collection.RemoveGame("0000000000000002")
			Expect(collection.Games().Files).To(HaveLen(1))
		})
	})
	Describe("Filter", func() {
		var (
			myMockConfig *mock_repository.MockConfig
			ctrl         *gomock.Controller
		)
		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())
		})
		JustBeforeEach(func() {
			myMockConfig = mock_repository.NewMockConfig(ctrl)
			customDB := make(map[string]repository.CustomDBEntry)
			custom1 := repository.CustomDBEntry{
				ID:              "0000000000000001",
				Languages:       []string{"FR", "EN", "US"},
				NumberOfPlayers: 1,
			}
			customDB["0000000000000001"] = custom1
			custom2 := repository.CustomDBEntry{
				ID:              "0000000000000002",
				Languages:       []string{"JP"},
				NumberOfPlayers: 2,
			}
			customDB["0000000000000001"] = custom1
			customDB["0000000000000002"] = custom2

			myMockConfig.EXPECT().
				Host().
				Return("tinshop.example.com").
				AnyTimes()
			myMockConfig.EXPECT().
				CustomDB().
				Return(customDB).
				AnyTimes()
			myMockConfig.EXPECT().
				BannedTheme().
				Return(nil).
				AnyTimes()

			collection.OnConfigUpdate(myMockConfig)

			newGames := make([]repository.FileDesc, 0)
			newFile1 := repository.FileDesc{
				Size:     1,
				Path:     "here",
				GameID:   "0000000000000001",
				GameInfo: "[0000000000000001][v0].nsp",
				HostType: repository.LocalFile,
			}
			newFile2 := repository.FileDesc{
				Size:     22,
				Path:     "here",
				GameID:   "0000000000000002",
				GameInfo: "[0000000000000002][v0].nsp",
				HostType: repository.LocalFile,
			}
			newGames = append(newGames, newFile1)
			newGames = append(newGames, newFile2)
			collection.AddNewGames(newGames)
		})
		It("Filtering world", func() {
			filteredGames := collection.Filter("WORLD")
			Expect(len(filteredGames.Titledb)).To(Equal(2))
			Expect(filteredGames.Titledb["0000000000000001"]).To(Not(BeNil()))
			Expect(filteredGames.Titledb["0000000000000002"]).To(Not(BeNil()))
			Expect(len(filteredGames.Files)).To(Equal(2))
		})
		It("Filtering US", func() {
			filteredGames := collection.Filter("US")
			Expect(len(filteredGames.Titledb)).To(Equal(1))
			Expect(filteredGames.Titledb["0000000000000001"]).To(Not(BeNil()))
			Expect(filteredGames.Titledb["0000000000000002"]).To(BeNil())
			Expect(len(filteredGames.Files)).To(Equal(1))
		})
		It("Filtering non existing language entry (HK)", func() {
			filteredGames := collection.Filter("HK")
			Expect(len(filteredGames.Titledb)).To(Equal(0))
			Expect(len(filteredGames.Files)).To(Equal(0))
		})
		It("Filtering multi", func() {
			filteredGames := collection.Filter("MULTI")
			Expect(len(filteredGames.Titledb)).To(Equal(1))
			Expect(filteredGames.Titledb["0000000000000001"]).To(BeNil())
			Expect(filteredGames.Titledb["0000000000000002"]).To(Not(BeNil()))
			Expect(len(filteredGames.Files)).To(Equal(1))
		})
	})
})
