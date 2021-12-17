package repository

import "net/http"

// GamesSource interface
type GamesSource interface {
	load(sources []string)
	download(w http.ResponseWriter, r *http.Request, game string, path string)
}

// GameID interface
type GameID interface {
	FullID() string
	ShortID() string
	Extension() string
}

// Config interface
type Config interface {
	RootShop() string
	SetRootShop(string)
	Host() string
	Protocol() string
	Port() int

	DebugNfs() bool
	Directories() []string
	NfsShares() []string
	ShopTitle() string
	ShopTemplateData() ShopTemplate
	SetShopTemplateData(ShopTemplate)
}

// ShopTemplate contains all variables used for shop template
type ShopTemplate struct {
	ShopTitle string
}
