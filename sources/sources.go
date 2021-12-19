package sources

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

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
	test := filepath.Clean(directory)
	fmt.Println(test)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					// fmt.Println("Channel is closed!")
					return
				}
				fmt.Println("New event!")
				if event.Name != "" {
					fmt.Println(event.Name)
				}
				if event.Op.String() != "" {
					fmt.Println(event.Op.String())
				}
				switch {
				case event.Op&fsnotify.Write == fsnotify.Write:
					log.Printf("Write:  %s: %s", event.Op, event.Name)
				case event.Op&fsnotify.Create == fsnotify.Create:
					log.Printf("Create: %s: %s", event.Op, event.Name)
				case event.Op&fsnotify.Remove == fsnotify.Remove:
					log.Printf("Remove: %s: %s", event.Op, event.Name)
				case event.Op&fsnotify.Rename == fsnotify.Rename:
					log.Printf("Rename: %s: %s", event.Op, event.Name)
				case event.Op&fsnotify.Chmod == fsnotify.Chmod:
					log.Printf("Chmod:  %s: %s", event.Op, event.Name)
				}
			case err := <-watcher.Errors:
				if err != nil {
					log.Println("watcher error:", err)
				}
			}
		}
	}()

	err = watcher.Add(directory + "/")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("end of watch function")
}
