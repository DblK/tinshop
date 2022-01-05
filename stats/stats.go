package stats

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/DblK/tinshop/repository"
	bolt "go.etcd.io/bbolt"
)

type stat struct {
	db *bolt.DB
}

// New create a new stats object
func New() repository.Stats {
	return &stat{}
}

func (s *stat) initDB() {
	_ = s.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("switch"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
}

func (s *stat) Load() {
	db, err := bolt.Open("stats.db", 0600, nil)
	if err != nil {
		fmt.Println(err)
	}
	s.db = db

	s.initDB()
}

func (s *stat) Close() error {
	return s.db.Close()
}

// Summary return the summary of all stats
func (s *stat) Summary() repository.StatsSummary {
	// fmt.Println(b.Stats().KeyN) // Num of element
	return repository.StatsSummary{
		NumberVisit: 42,
	}
}

// DownloadAsked compute stats when we download a game
func (s *stat) DownloadAsked(IP string, gameID string) error {
	fmt.Println("[Stats] DownloadAsked", IP, gameID)
	// TODO: Add in global download games stats
	// TODO: Add in global IP download stats

	return nil
}

// ListVisit count every visit to the listing page (either root or filter)
func (s *stat) ListVisit(console *repository.Switch) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		// Retrieve the users bucket.
		// This should be created when the DB is first opened.
		b := tx.Bucket([]byte("switch"))

		// Generate ID for the user.
		// This returns an error only if the Tx is closed or not writeable.
		// That can't happen in an Update() call so I ignore the error check.
		id, _ := b.NextSequence()
		console.ID = int(id)

		// Marshal user data into bytes.
		buf, err := json.Marshal(console)
		if err != nil {
			return err
		}

		// Persist bytes to users bucket.
		return b.Put(itob(console.ID), buf)
	})
}

// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
