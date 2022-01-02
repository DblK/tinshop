package stats

import (
	"fmt"

	"github.com/DblK/tinshop/repository"
	bolt "go.etcd.io/bbolt"
)

type stat struct {
}

// New create a new stats object
func New() repository.Stats {
	return &stat{}
}

func (*stat) Load() {
	db, err := bolt.Open("stats.db", 0666, nil)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
}
