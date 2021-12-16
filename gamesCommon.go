package main

import (
	"log"
)

func AddNewGames(newGames []FileDesc) {
	log.Printf("\n\nAdd new games...\n")
	var gameList = make([]interface{}, 0)
	for _, file := range newGames {
		game := make(map[string]interface{})
		game["url"] = rootShop + "/games/" + file.gameId + "#" + file.gameInfo
		game["size"] = file.size
		gameList = append(gameList, game)

		if library[file.gameId] != nil {
			// Verify already present and not update nor dlc
			if Games["titledb"].(map[string]interface{})[file.gameId] != nil && library[file.gameId].(map[string]interface{})["iconUrl"] != nil {
				log.Println("Already added id!", file.gameId, file.path)
			} else {
				Games["titledb"].(map[string]interface{})[file.gameId] = library[file.gameId]
			}
		} else {
			log.Println("Game not found in database!", file.gameInfo, file.path)
		}
	}
	Games["files"] = append(Games["files"].([]interface{}), gameList...)
	log.Printf("Added %d games in your library\n", len(gameList))
}
