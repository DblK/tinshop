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
	DebugNfs() bool
	Directories() []string
	NfsShares() []string
	ShopTitle() string
	ShopTemplateData() ShopTemplate
}

type ShopTemplate struct {
	ShopTitle string
}
