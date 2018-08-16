package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	bolt "github.com/coreos/bbolt"
)

var bucketName = "gifs"
var db = initDB()

type Gif struct {
	ID       string `json:"checksum"`
	BaseName string `json:"base_name"`
	FileSize int    `json:"file_size"`
	Count    int    `json:"count"`
}

func main() {
	var gifs []Gif
	var abcs = strings.Split("abcdefghijklmnopqrstuvwxyz1234567890", "")
	for _, a := range abcs {
		gifs = append(gifs, Gif{a + "-12345", "taylor-" + a + ".gif", 123, 0})
	}

	// save all data
	for _, g := range gifs {
		_, err := g.save()
		if err != nil {
			fmt.Println(err)
		}
		for x := 0; x < 10; x++ {
			_, err = g.increment()
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	// TODO: See how many records exist?
	// retrieve all data
	for _, g := range gifs {
		a, err := find(g.ID)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(a)
	}

	_, err := find("abc")
	if err != nil {
		fmt.Println(err)
	}

	s, _ := count()
	fmt.Println(s, "records stored")
}

func initDB() *bolt.DB {
	os.Remove("data/test.db")
	db, err := bolt.Open("data/test.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	} else {
		initBucket(db)
	}
	return db
}

// create the bucket
func initBucket(db *bolt.DB) {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		return err
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (g Gif) json() []byte {
	data, _ := json.Marshal(g)
	return data
}

func (g *Gif) increment() (bool, error) {
	g.Count++
	return g.save()
}

func (g Gif) save() (bool, error) {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		err := b.Put([]byte(g.ID), g.json())
		return err
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func find(checksum string) (gif Gif, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		v := b.Get([]byte(checksum))
		if v != nil {
			json.Unmarshal(v, &gif)
		} else {
			return fmt.Errorf("Unable to find id \"%s\"", checksum)
		}
		return nil
	})
	return
}

func count() (int, error) {
	var s bolt.BucketStats
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		s = b.Stats()
		return nil
	})
	return s.KeyN, nil
}
