package storage

import (
	"log"
	"testing"

	badger "github.com/dgraph-io/badger/v3"
)

func TestNew(t *testing.T) {
	// Open the Badger database located in the /tmp/badger directory.
	// It will be created if it doesn't exist.
	db, err := badger.Open(badger.DefaultOptions("./badger"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}
