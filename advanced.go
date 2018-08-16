package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
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
	gifs = append(gifs, Gif{"abc", "taylor.gif", 123, 0})
	gifs = append(gifs, Gif{"def", "swift.gif", 234, 0})
	gifs = append(gifs, Gif{"ghi", "rocks.gif", 345, 1})

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

	a, _ := find("abc")
	fmt.Println(a)
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

func find(checksum string) (Gif, error) {
	var gif Gif
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		v := b.Get([]byte(checksum))
		json.Unmarshal(v, &gif)
		return nil
	})
	return gif, err
}
