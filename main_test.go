package main_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	tinshop "github.com/DblK/tinshop"
	"github.com/DblK/tinshop/repository"
)

var _ = Describe("Main", func() {
	Describe("HomeHandler", func() {
		var (
			req     *http.Request
			handler http.Handler
			writer  *httptest.ResponseRecorder
		)

		Context("With empty collection", func() {
			BeforeEach(func() {
				handler = http.HandlerFunc(tinshop.HomeHandler)
				req = httptest.NewRequest(http.MethodGet, "/", nil)
				writer = httptest.NewRecorder()
			})
			It("Verify response without data", func() {
				handler.ServeHTTP(writer, req)
				Expect(writer.Code).To(Equal(http.StatusNotFound))
			})
			It("Verify empty response", func() {
				tinshop.ResetTinshop()
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
		XContext("With collection", func() {
			BeforeEach(func() {
				handler = http.HandlerFunc(tinshop.HomeHandler)
				req = httptest.NewRequest(http.MethodGet, "/", nil)
				writer = httptest.NewRecorder()
				tinshop.ResetTinshop()
				// Add mock collection
			})
			It("Verify status code", func() {

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
	})
})
