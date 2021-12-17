package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/dblk/tinshop/utils"
)

func loadGamesDirectories(singleSource bool) {
	for _, directory := range configServer.Directories() {
		err := loadGamesDirectory(directory)

		if err != nil {
			fmt.Println(err.Error())
			if len(configServer.Directories()) == 1 && err.Error() == "lstat ./games: no such file or directory" && singleSource {
				log.Fatal("You must create a folder 'games' and put your games inside or use config.yml to add sources!")
			} else {
				log.Println(err)
			}
		}
	}
}

func loadGamesDirectory(directory string) error {
	log.Printf("Loading games from directory '%s'...\n", directory)
	var newGameFiles []FileDesc
	// Walk through games directory
	err := filepath.Walk(directory,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				newFile := FileDesc{size: info.Size(), path: path}
				names := utils.ExtractGameID(path)

				if names.ShortID() != "" {
					newFile.gameID = names.ShortID()
					newFile.gameInfo = names.FullID()
					newFile.hostType = LocalFile
					newGameFiles = append(newGameFiles, newFile)
				} else {
					log.Println("Ignoring file because parsing failed", path)
				}
			}
			return nil
		})
	if err != nil {
		return err
	}
	gameFiles = append(gameFiles, newGameFiles...)

	// Add all files
	if len(newGameFiles) > 0 {
		AddNewGames(newGameFiles)
	}

	return nil
}

func downloadLocalFile(w http.ResponseWriter, r *http.Request, game, path string) {
	f, err := os.Open(path)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer f.Close()

	fi, err := f.Stat()

	if err == nil {
		http.ServeContent(w, r, game, fi.ModTime(), f)
	} else {
		http.ServeContent(w, r, game, time.Time{}, f)
	}
}
