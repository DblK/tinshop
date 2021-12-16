package repository

import "net/http"

type GamesSource interface {
	load(sources []string)
	download(w http.ResponseWriter, r *http.Request, game string, path string)
}

type GameId interface {
	FullId() string
	ShortId() string
	Extension() string
}
