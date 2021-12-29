package directory

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/DblK/tinshop/repository"
)

type directorySource struct {
	gameFiles  []repository.FileDesc
	collection repository.Collection
}

// New create a directory source
func New(collection repository.Collection) repository.Source {
	return &directorySource{
		gameFiles:  make([]repository.FileDesc, 0),
		collection: collection,
	}
}

func (src *directorySource) Download(w http.ResponseWriter, r *http.Request, game, path string) {
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

func (src *directorySource) Load(directories []string, uniqueSource bool) {
	for _, directory := range directories {
		err := src.loadGamesDirectory(directory)

		if err != nil {
			if strings.Contains(err.Error(), "no such file or directory") {
				if len(directories) == 1 && uniqueSource {
					log.Fatal("You must create a folder 'games' and put your games inside or use config.yaml to add sources!")
				}
			} else {
				log.Println(err)
			}
		}
	}
}

func (src *directorySource) Reset() {
	watcherDirectories = newWatcher()
	src.gameFiles = make([]repository.FileDesc, 0)
}

func (src *directorySource) UnWatchAll() {
	removeGamesWatcherDirectories()
}

func (src *directorySource) GetFiles() []repository.FileDesc {
	return src.gameFiles
}
