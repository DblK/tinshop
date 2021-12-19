package sources

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/dblk/tinshop/repository"
	"github.com/dblk/tinshop/utils"
	"gopkg.in/fsnotify.v1"
)

var gameFiles []repository.FileDesc

// OnConfigUpdate from all sources
func OnConfigUpdate(cfg repository.Config) {
	log.Println("Sources loading...")
	gameFiles = make([]repository.FileDesc, 0)
	loadGamesDirectories(cfg.Directories(), len(cfg.NfsShares()) == 0)
	loadGamesNfsShares(cfg.NfsShares())

	// Add watcher for directories
	if len(cfg.Directories()) != 0 {
		// TODO: Make loop
		watchDirectory(cfg.Directories()[0])
	}
}

// BeforeConfigUpdate from all sources
func BeforeConfigUpdate(cfg repository.Config) {
	// TODO: Stop watching previous directories
	fmt.Println("Code this!")
}

// GetFiles returns all games files in various sources
func GetFiles() []repository.FileDesc {
	return gameFiles
}

// AddFiles add files to global sources
func AddFiles(files []repository.FileDesc) {
	gameFiles = append(gameFiles, files...)
}

// DownloadGame method provide the file based on the source storage
func DownloadGame(gameID string, w http.ResponseWriter, r *http.Request) {
	idx := utils.Search(len(GetFiles()), func(index int) bool {
		return GetFiles()[index].GameID == gameID
	})

	if idx == -1 {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("Game '%s' not found!", gameID)
		return
	}
	log.Println("Retrieving from location '" + GetFiles()[idx].Path + "'")
	switch GetFiles()[idx].HostType {
	case repository.LocalFile:
		downloadLocalFile(w, r, gameID, GetFiles()[idx].Path)
	case repository.NFSShare:
		downloadNfsFile(w, r, GetFiles()[idx].Path)

	default:
		w.WriteHeader(http.StatusNotImplemented)
		log.Printf("The type '%s' is not implemented to download game", GetFiles()[idx].HostType)
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
						fmt.Println("Remove file")
						eventsWG.Done()
						return
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
		watcher.Add(directory)
		initWG.Done()   // done initializing the watch in this go routine, so the parent routine can move on...
		eventsWG.Wait() // now, wait for event loop to end in this go-routine...
	}()
	initWG.Wait() // make sure that the go routine above fully ended before returning
	fmt.Println("end of watch function")
}
