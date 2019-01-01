// Package bolt contains the types and methods needed for using a urlshort instance backed in a Bolt database.
// There's a convenience method New to avoid the need of creating the store and then the handler
// in the client code.
package bolt

import (
	"net/http"

	"github.com/pkg/errors"
	"github.com/roberveral/gophercises/urlshort2/urlshort"
	log "github.com/sirupsen/logrus"
	bbolt "go.etcd.io/bbolt"
)

// This type implements the urlshort.Store interface using the provided
// bbolt.DB to perform queries/updates.
type store struct {
	db *bbolt.DB
}

// The name of the Bolt bucket where the path mappings will be stored.
const bucketName string = "PathsToUrls"

// NewStore creates a new urlshort.Store which performs the queries/updates
// against the given Bolt database.
//
// This constructor initializes all the required buckets for the internal store
// operation. If there's a failure initializing the bucket, an error is returned.
func NewStore(db *bbolt.DB) (urlshort.Store, error) {
	err := db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		return err
	})

	if err != nil {
		return nil, errors.Wrap(err, "Unable to initialize Bolt store")
	} else {
		return &store{db}, nil
	}
}

// New returns an http.Handler configured to perform the urlshort tasks using a
// Bolt Store created against the given database.
//
// This method just initializes the store and calls urlshort.New to obtain the handler.
func New(db *bbolt.DB) (http.Handler, error) {
	store, err := NewStore(db)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to initialize Bolt-configured urlshort handler")
	}

	return urlshort.New(store), nil
}

// Get is implemented by a simple read from the bucket using the path as key.
func (s *store) Get(path string) (string, bool) {
	var url []byte
	err := s.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		value := bucket.Get([]byte(path))

		// Byte slice has no meaning outside the transaction so we need to
		// copy it while maintaining the 'nil' meaning.
		if value != nil {
			url = make([]byte, len(value))
			copy(url, value)
		} else {
			url = nil
		}

		return nil
	})

	if err != nil {
		log.Errorf("There was an error getting path '%s' from Bolt store: %+v", path, err)
	}

	return string(url), url != nil
}

// Put is implemented as a write to the bucket using the path as the key and the url as the value.
// NOTE: this operation overwrittes the previous path value in case of collsion.
func (s *store) Put(path string, url string) error {
	err := s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		err := b.Put([]byte(path), []byte(url))
		return err
	})

	return errors.Wrapf(err, "Unable to put key '%s' in Bolt store", path)
}
