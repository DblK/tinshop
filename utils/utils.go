package utils

import (
	"regexp"
	"strings"

	"github.com/dblk/tinshop/gameId"
	"github.com/dblk/tinshop/repository"
)

// ExtractGameId from fileName the id of game and version
func ExtractGameId(fileName string) repository.GameId {
	ext := strings.Split(fileName, ".")
	re := regexp.MustCompile(`\[(\w{16})\]\[(v\d+)\]`)
	matches := re.FindStringSubmatch(fileName)

	if len(matches) != 3 {
		return gameId.New("", "", "")

	}

	return gameId.New(strings.ToUpper(matches[1]), "["+strings.ToUpper(matches[1])+"]["+matches[2]+"]."+ext[len(ext)-1], ext[len(ext)-1])
}

func Search(length int, f func(index int) bool) int {
	for index := 0; index < length; index++ {
		if f(index) {
			return index
		}
	}
	return -1
}
