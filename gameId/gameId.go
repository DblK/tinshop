package gameId

import (
	"github.com/dblk/tinshop/repository"
)

type gameId struct {
	fullId    string
	shortId   string
	extension string
}

func New(shortId string, fullId string, extension string) repository.GameId {
	return &gameId{
		fullId:    fullId,
		shortId:   shortId,
		extension: extension,
	}
}

func (game *gameId) SetFullId(fullId string) {
	game.fullId = fullId
}

func (game *gameId) FullId() string {
	return game.fullId
}
func (game *gameId) ShortId() string {
	return game.shortId
}
func (game *gameId) Extension() string {
	return game.extension
}
