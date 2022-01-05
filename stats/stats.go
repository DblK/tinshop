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
		_, err := tx.CreateBucketIfNotExists([]byte("global"))
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
func (s *stat) Summary() (repository.StatsSummary, error) {
	var visit uint64
	var uniqueSwitch int
	var consoles map[string]interface{}
	var download uint64
	var downloadDetails map[string]interface{}

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("global"))
		visit = byteToUint64(b.Get([]byte("visit")))

		var errConsoles error
		consoles, errConsoles = byteToMap(b.Get([]byte("switch")))
		if errConsoles != nil {
			return errConsoles
		}
		uniqueSwitch = len(consoles)

		download = byteToUint64(b.Get([]byte("download")))

		var errDownloadDetails error
		downloadDetails, errDownloadDetails = byteToMap(b.Get([]byte("downloadDetails")))
		if errDownloadDetails != nil {
			return errDownloadDetails
		}

		return nil
	})
	if err != nil {
		return repository.StatsSummary{}, err
	}

	return repository.StatsSummary{
		Visit:           visit,
		UniqueSwitch:    uint64(uniqueSwitch),
		VisitPerSwitch:  consoles,
		DownloadAsked:   download,
		DownloadDetails: downloadDetails,
	}, nil
}

// DownloadAsked compute stats when we download a game
func (s *stat) DownloadAsked(IP string, gameID string) error {
	fmt.Println("[Stats] DownloadAsked", IP, gameID)
	// TODO: Add in global IP download stats

	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("global"))

		// Handle download
		download := byteToUint64(b.Get([]byte("download")))
		errDownload := b.Put([]byte("download"), itob(download+1))
		if errDownload != nil {
			return errDownload
		}

		// Handle download per IP
		allDownloads, err := byteToMap(b.Get([]byte("downloadDetails")))
		if err != nil {
			return err
		}
		if allDownloads[IP] == nil {
			allDownloads[IP] = make([]interface{}, 0)
		}
		allDownloads[IP] = append(allDownloads[IP].([]interface{}), gameID)
		buf, err := json.Marshal(allDownloads)
		if err != nil {
			return err
		}
		return b.Put([]byte("downloadDetails"), buf)
	})
}

// ListVisit count every visit to the listing page (either root or filter)
func (s *stat) ListVisit(console *repository.Switch) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("global"))

		// Handle visit
		visit := byteToUint64(b.Get([]byte("visit")))
		errVisit := b.Put([]byte("visit"), itob(visit+1))
		if errVisit != nil {
			return errVisit
		}

		// Handle visit per switch
		consoles, err := byteToMap(b.Get([]byte("switch")))
		if err != nil {
			return err
		}
		currentID := console.UID
		if currentID == "" {
			currentID = "Unknown-" + console.IP
		}

		if consoles[currentID] == nil {
			consoles[currentID] = float64(0)
		}
		consoles[currentID] = uint64(consoles[currentID].(float64)) + 1
		buf, err := json.Marshal(consoles)
		if err != nil {
			return err
		}
		return b.Put([]byte("switch"), buf)
	})
}

func byteToMap(bytes []byte) (map[string]interface{}, error) {
	val := make(map[string]interface{})
	if len(bytes) > 0 {
		err := json.Unmarshal(bytes, &val)
		if err != nil {
			return make(map[string]interface{}), err
		}
	}
	return val, nil
}

func byteToUint64(bytes []byte) uint64 {
	num := uint64(0)
	if len(bytes) > 0 {
		num = binary.BigEndian.Uint64(bytes)
	}
	return num
}

// itob returns an 8-byte big endian representation of v.
func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}
