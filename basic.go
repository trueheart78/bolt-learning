package main

import (
	"fmt"
	"os"
	"time"

	bolt "github.com/coreos/bbolt"
)

var db *bolt.DB
var bucketName = "gifs"

func main() {
	fmt.Println("Go")

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

	// save some data!
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		err := b.Put([]byte("answer"), []byte("42"))
		return err
	})
	if err != nil {
		fmt.Println(err)
	}

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		v := b.Get([]byte("answer"))
		fmt.Printf("The answer is: %s\n", v)
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Go")
}
