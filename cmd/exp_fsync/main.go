package main

import (
	"fmt"
	"strconv"

	"go.etcd.io/bbolt"

	log "github.com/sirupsen/logrus"
)

// This program will run a test, in which database with a single number
// will be increased a lot of times. This allows to test performance.
func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)
	log.SetLevel(log.DebugLevel)

	databaseFile := "data.db"

	db, err := bbolt.Open(databaseFile, 0666, nil)
	if err != nil {
		log.WithError(err).Fatal("failed to open db")
	}
	defer db.Close()

	var key = []byte("value")
	_ = db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("MyBucket"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		_ = b.Put(key, []byte("0"))
		return nil
	})

	for j := 0; j < 100; j++ {
		var newVal int
		for i := 0; i < 100; i++ {
			var err error
			newVal, err = increment(db, key)
			if err != nil {
				log.WithError(err).Fatal("failed to increment")
			}
		}

		log.WithField("val", newVal).WithField("it", j).Info("done increments")
	}

	_ = db.Sync()
}

func increment(db *bbolt.DB, key []byte) (int, error) {
	var i int
	var err error

	err = db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("MyBucket"))

		b := bucket.Get(key)

		i, err = strconv.Atoi(string(b))
		if err != nil {
			return err
		}

		i++

		str := strconv.Itoa(i)
		err = bucket.Put(key, []byte(str))
		return err
	})

	return i, err
}
