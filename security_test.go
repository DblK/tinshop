package main_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	main "github.com/DblK/tinshop"
	"github.com/DblK/tinshop/mock_repository"
	"github.com/DblK/tinshop/repository"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Security", func() {
	Describe("TinfoilMiddleware", func() {
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

		Context("No security", func() {
			BeforeEach(func() {
				r := mux.NewRouter()
				r.Use(myShop.TinfoilMiddleware)
				r.HandleFunc("/", myShop.HomeHandler)
				handler = r
				req = httptest.NewRequest(http.MethodGet, "/", nil)
				writer = httptest.NewRecorder()
			})
			It("without any headers", func() {
				emptyCollection := &repository.GameType{}

				myMockCollection.EXPECT().
					Games().
					Return(*emptyCollection).
					AnyTimes()

				myMockConfig.EXPECT().
					DebugNoSecurity().
					Return(true).
					AnyTimes()

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
		Context("With security", func() {
			BeforeEach(func() {
				r := mux.NewRouter()
				r.Use(myShop.TinfoilMiddleware)
				r.HandleFunc("/", myShop.HomeHandler)
				handler = r
				req = httptest.NewRequest(http.MethodGet, "/", nil)
				writer = httptest.NewRecorder()
			})
			It("test for blacklisted switch", func() {
				emptyCollection := &repository.GameType{}

				myMockCollection.EXPECT().
					Games().
					Return(*emptyCollection).
					AnyTimes()

				myMockConfig.EXPECT().
					DebugNoSecurity().
					Return(false).
					AnyTimes()

				myMockConfig.EXPECT().
					IsBlacklisted(gomock.Any()).
					Return(true).
					AnyTimes()

				shopTemplateData := &repository.ShopTemplate{
					ShopTitle: "Unit Test",
				}
				myMockConfig.EXPECT().
					ShopTemplateData().
					Return(*shopTemplateData).
					AnyTimes()

				handler.ServeHTTP(writer, req)

				Expect(writer.Code).To(Equal(http.StatusOK))

				var list repository.GameType
				err := json.NewDecoder(writer.Body).Decode(&list)
				Expect(err).To(Not(BeNil()))
			})
			It("test for banned theme switch", func() {
				emptyCollection := &repository.GameType{}

				myMockCollection.EXPECT().
					Games().
					Return(*emptyCollection).
					AnyTimes()

				myMockConfig.EXPECT().
					DebugNoSecurity().
					Return(false).
					AnyTimes()

				myMockConfig.EXPECT().
					IsBlacklisted(gomock.Any()).
					Return(false).
					AnyTimes()

				myMockConfig.EXPECT().
					IsBannedTheme(gomock.Any()).
					Return(true).
					AnyTimes()

				shopTemplateData := &repository.ShopTemplate{
					ShopTitle: "Unit Test",
				}
				myMockConfig.EXPECT().
					ShopTemplateData().
					Return(*shopTemplateData).
					AnyTimes()

				handler.ServeHTTP(writer, req)

				Expect(writer.Code).To(Equal(http.StatusOK))

				var list repository.GameType
				err := json.NewDecoder(writer.Body).Decode(&list)
				Expect(err).To(Not(BeNil()))
			})
			It("test for an existing user agent", func() {
				emptyCollection := &repository.GameType{}

				myMockCollection.EXPECT().
					Games().
					Return(*emptyCollection).
					AnyTimes()

				myMockConfig.EXPECT().
					DebugNoSecurity().
					Return(false).
					AnyTimes()

				myMockConfig.EXPECT().
					IsBlacklisted(gomock.Any()).
					Return(false).
					AnyTimes()

				myMockConfig.EXPECT().
					IsBannedTheme(gomock.Any()).
					Return(false).
					AnyTimes()

				shopTemplateData := &repository.ShopTemplate{
					ShopTitle: "Unit Test",
				}
				myMockConfig.EXPECT().
					ShopTemplateData().
					Return(*shopTemplateData).
					AnyTimes()

				req.Header.Set("User-Agent", "Tinshop testing!")
				handler.ServeHTTP(writer, req)

				Expect(writer.Code).To(Equal(http.StatusOK))

				var list repository.GameType
				err := json.NewDecoder(writer.Body).Decode(&list)
				Expect(err).To(Not(BeNil()))
			})
			It("test for with missing mandatory headers", func() {
				emptyCollection := &repository.GameType{}

				myMockCollection.EXPECT().
					Games().
					Return(*emptyCollection).
					AnyTimes()

				myMockConfig.EXPECT().
					DebugNoSecurity().
					Return(false).
					AnyTimes()

				myMockConfig.EXPECT().
					IsBlacklisted(gomock.Any()).
					Return(false).
					AnyTimes()

				myMockConfig.EXPECT().
					IsBannedTheme(gomock.Any()).
					Return(false).
					AnyTimes()

				shopTemplateData := &repository.ShopTemplate{
					ShopTitle: "Unit Test",
				}
				myMockConfig.EXPECT().
					ShopTemplateData().
					Return(*shopTemplateData).
					AnyTimes()

				handler.ServeHTTP(writer, req)

				Expect(writer.Code).To(Equal(http.StatusOK))

				var list repository.GameType
				err := json.NewDecoder(writer.Body).Decode(&list)
				Expect(err).To(Not(BeNil()))
			})
		})
	})
})
