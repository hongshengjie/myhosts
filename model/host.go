package model

import (
	"encoding/binary"
	"encoding/json"
	"errors"

	"github.com/boltdb/bolt"
)

type Host struct {
	ID      uint64 `storm:"increment" json:"id"`
	Title   string `storm:"unique" json:"title"`
	Content string `json:"content"`
	Enable  bool   `json:"enable"`
}

func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

func HostCreate(h *Host) error {
	return db.Update(func(t *bolt.Tx) error {
		b := t.Bucket(bucketName)
		autoid, err := b.NextSequence()
		if err != nil {
			return err
		}
		h.ID = autoid
		v, err := json.Marshal(h)
		if err != nil {
			return err
		}
		return b.Put(itob(h.ID), v)
	})
}

var NotFound = errors.New("Not Found")

func HostGet(id uint64) (*Host, error) {
	var h Host
	err := db.View(func(t *bolt.Tx) error {
		b := t.Bucket(bucketName)
		v := b.Get(itob(id))
		if v == nil {
			return NotFound
		}
		if err := json.Unmarshal(v, &h); err != nil {
			return err
		}
		return nil

	})
	if err != nil {
		return nil, err
	}
	return &h, nil
}

func HostListAll() ([]*Host, error) {
	var list []*Host

	err := db.View(func(t *bolt.Tx) error {
		b := t.Bucket(bucketName)
		err := b.ForEach(func(k, v []byte) error {
			var h Host
			if err := json.Unmarshal(v, &h); err != nil {
				return err
			}
			list = append(list, &h)
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return list, nil
}

func HostUpdate(h *Host) error {
	return db.Update(func(t *bolt.Tx) error {
		b := t.Bucket(bucketName)
		content, err := json.Marshal(h)
		if err != nil {
			return err
		}
		return b.Put(itob(h.ID), content)
	})
}

func HostDelete(h *Host) error {
	return db.Update(func(t *bolt.Tx) error {
		b := t.Bucket(bucketName)
		return b.Delete(itob(h.ID))
	})
}
