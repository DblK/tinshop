package config_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/dblk/tinshop/config"
)

var _ = Describe("Config", func() {
	It("Ensure to be able to set a RootShop", func() {
		config.GetConfig().SetRootShop("http://tinshop.example.com")
		cfg := config.GetConfig()

		Expect(cfg.RootShop()).To(Equal("http://tinshop.example.com"))
	})
})
