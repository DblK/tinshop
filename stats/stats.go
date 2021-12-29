package stats

import (
	"fmt"

	bolt "go.etcd.io/bbolt"
)

func Load() {
	db, err := bolt.Open("stats.db", 0666, nil)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
}
