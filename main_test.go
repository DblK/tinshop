package main_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	main "github.com/DblK/tinshop"
	"github.com/DblK/tinshop/mock_repository"
	"github.com/DblK/tinshop/repository"
)

var _ = Describe("Main", func() {
	Describe("HomeHandler", func() {
		var (
			req              *http.Request
			handler          http.Handler
			writer           *httptest.ResponseRecorder
			myMockCollection *mock_repository.MockCollection
			myMockSources    *mock_repository.MockSources
			myMockConfig     *mock_repository.MockConfig
			ctrl             *gomock.Controller
			myShop           *main.TinShop
		)

		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())
			myMockCollection = mock_repository.NewMockCollection(ctrl)
			myMockSources = mock_repository.NewMockSources(ctrl)
			myMockConfig = mock_repository.NewMockConfig(ctrl)
			myShop = &main.TinShop{}
		})

		JustBeforeEach(func() {
			myShop.Shop = repository.Shop{}
			myShop.Shop.Config = myMockConfig
			myShop.Shop.Collection = myMockCollection
			myShop.Shop.Sources = myMockSources
		})

		Context("With empty collection", func() {
			BeforeEach(func() {
				handler = http.HandlerFunc(myShop.HomeHandler)
				req = httptest.NewRequest(http.MethodGet, "/", nil)
				writer = httptest.NewRecorder()
			})
			It("Verify response without data", func() {
				myShop.Shop.Collection = nil
				handler.ServeHTTP(writer, req)
				Expect(writer.Code).To(Equal(http.StatusNotFound))
			})
			It("Verify empty response", func() {
				emptyCollection := &repository.GameType{}

				myMockCollection.EXPECT().
					Games().
					Return(*emptyCollection).
					AnyTimes()

				// tinshop.ResetTinshop(myShop)
				handler.ServeHTTP(writer, req)
				Expect(writer.Code).To(Equal(http.StatusOK))

				var list repository.GameType
				err := json.NewDecoder(writer.Body).Decode(&list)

				Expect(err).To(BeNil())
				Expect(list.Files).To(HaveLen(0))
				Expect(list.ThemeBlackList).To(BeNil())
				Expect(list.Success).To(BeEmpty())
				Expect(list.Titledb).To(HaveLen(0))
			})
		})
		Context("With collection", func() {
			JustBeforeEach(func() {
				handler = http.HandlerFunc(myShop.HomeHandler)
				req = httptest.NewRequest(http.MethodGet, "/", nil)
				writer = httptest.NewRecorder()
				// tinshop.ResetTinshop(myShop)
			})
			It("Verify status code", func() {
				smallCollection := &repository.GameType{}
				smallCollection.Files = make([]repository.GameFileType, 0)
				oneFile := &repository.GameFileType{
					Size: 42,
					URL:  "http://test.tinshop.io",
				}
				smallCollection.Files = append(smallCollection.Files, *oneFile)
				smallCollection.Success = "Welcome to your own shop!"
				smallCollection.Titledb = make(map[string]repository.TitleDBEntry)
				oneEntry := &repository.TitleDBEntry{
					ID: "0000000000000001",
				}
				smallCollection.Titledb["0000000000000001"] = *oneEntry

				myMockCollection.EXPECT().
					Games().
					Return(*smallCollection).
					AnyTimes()

				handler.ServeHTTP(writer, req)
				Expect(writer.Code).To(Equal(http.StatusOK))

				var list repository.GameType
				err := json.NewDecoder(writer.Body).Decode(&list)

				Expect(err).To(BeNil())
				Expect(list.Files).To(HaveLen(1))
				Expect(list.ThemeBlackList).To(BeNil())
				Expect(list.Success).To(Equal("Welcome to your own shop!"))
				Expect(list.Titledb).To(HaveLen(1))
			})
		})
	})
})
