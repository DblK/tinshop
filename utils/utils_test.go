package utils_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/dblk/tinshop/utils"
)

var _ = Describe("ExtractGameId", func() {
	Context("Should succeed", func() {
		It("Nicely separated groups", func() {
			game := utils.ExtractGameId("Paw Patrol Mighty Pups Save Adventure Bay [01001F201121E800][v131072] (1.58 GB).nsz")

			Expect(game.Extension()).To(Equal("nsz"))
			Expect(game.ShortId()).To(Equal("01001F201121E800"))
			Expect(game.FullId()).To(Equal("[01001F201121E800][v131072].nsz"))
		})
		It("Make upper of Game Id", func() {
			game := utils.ExtractGameId("Game [01001f201121e800][v131072] (1.58 GB).nsz")

			Expect(game.Extension()).To(Equal("nsz"))
			Expect(game.ShortId()).To(Equal("01001F201121E800"))
			Expect(game.FullId()).To(Equal("[01001F201121E800][v131072].nsz"))
		})
		It("Group tied with parenthesis group", func() {
			game := utils.ExtractGameId("Paw Patrol Mighty Pups Save Adventure Bay [01001F201121E800][v131072](1.58 GB).nsz")

			Expect(game.Extension()).To(Equal("nsz"))
			Expect(game.ShortId()).To(Equal("01001F201121E800"))
			Expect(game.FullId()).To(Equal("[01001F201121E800][v131072].nsz"))
		})
		It("Nice filename with nsp file", func() {
			game := utils.ExtractGameId("Super Mario Odyssey [0100000000010000][v0].nsp")

			Expect(game.Extension()).To(Equal("nsp"))
			Expect(game.ShortId()).To(Equal("0100000000010000"))
			Expect(game.FullId()).To(Equal("[0100000000010000][v0].nsp"))
		})
		It("Nice separated DLC information", func() {
			game := utils.ExtractGameId("The Legend of Zelda Breath of the Wild [DLC Pack 1 The Master Trials] [01007EF00011F001][v196608].nsp")

			Expect(game.Extension()).To(Equal("nsp"))
			Expect(game.ShortId()).To(Equal("01007EF00011F001"))
			Expect(game.FullId()).To(Equal("[01007EF00011F001][v196608].nsp"))
		})
		It("Tied DLC info to game id and version", func() {
			game := utils.ExtractGameId("Fake - The Legend of Zelda Breath of the Wild [DLC Pack 1 The Master Trials][01007EF00011F001][v196608].nsp")

			Expect(game.Extension()).To(Equal("nsp"))
			Expect(game.ShortId()).To(Equal("01007EF00011F001"))
			Expect(game.FullId()).To(Equal("[01007EF00011F001][v196608].nsp"))
		})
		It("Tied DLC info with no space to game id and version", func() {
			game := utils.ExtractGameId("Fake - The Legend of Zelda Breath of the Wild [DLCPack1TheMasterTrials][01007EF00011F001][v196608].nsp")

			Expect(game.Extension()).To(Equal("nsp"))
			Expect(game.ShortId()).To(Equal("01007EF00011F001"))
			Expect(game.FullId()).To(Equal("[01007EF00011F001][v196608].nsp"))
		})
		It("Game inside sub directory", func() {
			game := utils.ExtractGameId("Fake - My Directory/Fake - [0100152000022800][v655360].nsz")

			Expect(game.Extension()).To(Equal("nsz"))
			Expect(game.ShortId()).To(Equal("0100152000022800"))
			Expect(game.FullId()).To(Equal("[0100152000022800][v655360].nsz"))
		})
	})
	Context("Should Fail", func() {
		It("Test with not size valid game id", func() {
			game := utils.ExtractGameId("Fake - My Game [NSP]/Fake - My Own Game [1231231][v0].nsz")

			Expect(game.Extension()).To(BeEmpty())
			Expect(game.ShortId()).To(BeEmpty())
			Expect(game.FullId()).To(BeEmpty())
		})
		It("Test with bad number of version", func() {
			game := utils.ExtractGameId("Fake - My Game [NSP]/Fake - My Own Game [0100152000022800][0].nsz")

			Expect(game.Extension()).To(BeEmpty())
			Expect(game.ShortId()).To(BeEmpty())
			Expect(game.FullId()).To(BeEmpty())
		})
		It("Test with no game id no version", func() {
			game := utils.ExtractGameId("Fake - Bad name.txt")

			Expect(game.Extension()).To(BeEmpty())
			Expect(game.ShortId()).To(BeEmpty())
			Expect(game.FullId()).To(BeEmpty())
		})
		It("Test with double extension", func() {
			game := utils.ExtractGameId("Fake - Bad name.old.txt")

			Expect(game.Extension()).To(BeEmpty())
			Expect(game.ShortId()).To(BeEmpty())
			Expect(game.FullId()).To(BeEmpty())
		})
	})
})
