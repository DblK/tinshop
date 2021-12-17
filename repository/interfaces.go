package repository

import "net/http"

type GamesSource interface {
	load(sources []string)
	download(w http.ResponseWriter, r *http.Request, game string, path string)
}

type GameID interface {
	FullID() string
	ShortID() string
	Extension() string
}

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

type ShopTemplate struct {
	ShopTitle string
}
