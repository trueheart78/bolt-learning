package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	bolt "github.com/coreos/bbolt"
)

var db *bolt.DB
var bucketName = "gifs"

type Gif struct {
	ID       string `json:"checksum"`
	BaseName string `json:"base_name"`
	FileSize int    `json:"file_size"`
}

func main() {
	var gifs []Gif
	gifs = append(gifs, Gif{"abc", "taylor.gif", 123})
	gifs = append(gifs, Gif{"def", "swift.gif", 234})
	gifs = append(gifs, Gif{"ghi", "rocks.gif", 345})

	os.Remove("data/test.db")
	db, err := bolt.Open("data/test.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		fmt.Println(err)
	}

	// create a bucket!
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}

	for _, g := range gifs {
		data, err := json.Marshal(g)
		if err != nil {
			fmt.Println(err)
		}
		// save some data!
		err = db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(bucketName))
			err := b.Put([]byte(g.ID), data)
			return err
		})
		if err != nil {
			fmt.Println(err)
		}
	}

	for _, g := range gifs {
		err = db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(bucketName))
			v := b.Get([]byte(g.ID))
			fmt.Printf("%s\n", v)
			return nil
		})
		if err != nil {
			fmt.Println(err)
		}
	}
}
