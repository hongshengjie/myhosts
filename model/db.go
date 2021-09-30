package model

import (
	"myhosts/app"
	"os"
	"path"

	"github.com/boltdb/bolt"
)

var database = "myhosts.db"
var db *bolt.DB
var bucketName = []byte("myhost")
var firstOpen bool

func FirstOpen() bool {
	return firstOpen
}

func Init() {
	var err error
	local := app.ConfigDir()
	if err != nil {
		panic(err)
	}

	os.MkdirAll(local, 0700)
	dbpath := path.Join(local, database)
	_, err = os.Stat(dbpath)
	if err != nil {
		if os.IsNotExist(err) {
			firstOpen = true
		}
	}

	db, err = bolt.Open(dbpath, 0600, nil)
	if err != nil {
		panic(err)
	}
	err = db.Update(func(t *bolt.Tx) error {
		if _, err := t.CreateBucketIfNotExists(bucketName); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func Close() {
	db.Close()
}
