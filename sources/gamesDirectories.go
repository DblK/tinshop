package sources

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	collection "github.com/DblK/tinshop/gamescollection"
	"github.com/DblK/tinshop/repository"
	"github.com/DblK/tinshop/utils"
	"gopkg.in/fsnotify.v1"
)

func loadGamesDirectories(directories []string, singleSource bool) {
	for _, directory := range directories {
		err := loadGamesDirectory(directory)

		if err != nil {
			if len(directories) == 1 && err.Error() == "lstat ./games: no such file or directory" && singleSource {
				log.Fatal("You must create a folder 'games' and put your games inside or use config.yml to add sources!")
			} else {
				log.Println(err)
			}
		}
	}
}

func loadGamesDirectory(directory string) error {
	log.Printf("Loading games from directory '%s'...\n", directory)

	// Add watcher for directories
	watchDirectory(directory)

	var newGameFiles []repository.FileDesc
	// Walk through games directory
	err := filepath.Walk(directory,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				extension := filepath.Ext(info.Name())
				if extension == ".nsp" || extension == ".nsz" {
					newFile := repository.FileDesc{Size: info.Size(), Path: path}
					names := utils.ExtractGameID(path)

					if names.ShortID() != "" {
						newFile.GameID = names.ShortID()
						newFile.GameInfo = names.FullID()
						newFile.HostType = repository.LocalFile
						newGameFiles = append(newGameFiles, newFile)
					} else {
						log.Println("Ignoring file because parsing failed", path)
					}
				}
			} else if info.IsDir() && path != directory {
				// TODO: Add watcher for the sub-directory
				fmt.Println("Need to watch sub directory", path)
			}
			return nil
		})
	if err != nil {
		return err
	}
	AddFiles(newGameFiles)

	// Add all files
	if len(newGameFiles) > 0 {
		collection.AddNewGames(newGameFiles)
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

func watchDirectory(directory string) {
	fmt.Println("Watching directory", directory)

	initWG := sync.WaitGroup{}
	initWG.Add(1)
	go func() {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		defer watcher.Close()

		eventsWG := sync.WaitGroup{}
		eventsWG.Add(1)
		go func() {
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok { // 'Events' channel is closed
						eventsWG.Done()
						return
					}
					const writeOrCreateMask = fsnotify.Write | fsnotify.Create
					if event.Op&writeOrCreateMask != 0 {
						fmt.Println("Changes", event)
					} else if event.Op&fsnotify.Remove != 0 {
						fmt.Println("Remove file", event)
					} else if event.Op&fsnotify.Rename == fsnotify.Rename {
						log.Printf("Rename: %s: %s", event.Op, event.Name)
					}

				case err, ok := <-watcher.Errors:
					if ok { // 'Errors' channel is not closed
						log.Printf("watcher error: %v\n", err)
					}
					eventsWG.Done()
					return
				}
			}
		}()
		errWatcher := watcher.Add(directory)
		initWG.Done()   // done initializing the watch in this go routine, so the parent routine can move on...
		eventsWG.Wait() // now, wait for event loop to end in this go-routine...
		if errWatcher != nil {
			eventsWG.Done()
		}
	}()
	initWG.Wait() // make sure that the go routine above fully ended before returning
	fmt.Println("end of watch function")
}
