package datebase

import (
	"time"
	"log"
	"encoding/json"
	"encoding/binary"
	"fmt"

	"github.com/boltdb/bolt"
)

type Blog struct{
	ID int
	Title string
	Content string
	Label []string
	Owner string
	CreateTime time.Time
}

// type User struct{
// 	ID int
// 	Name string
// 	Password string
// }

func ce(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Sections(DB *bolt.DB) []string {
	var buckets []string
	err := DB.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(id []byte, _ *bolt.Bucket) error {
			buckets = append(buckets, string(id))
			return nil
		})
	})
	ce(err)
	return buckets
}

func SaveBlog(DB *bolt.DB, section string, blog Blog) error{
	err := DB.Update(func(tx *bolt.Tx) error{
		// Retrieve the blogs bucket.
        // This should be created when the DB is first opened.
		bucket := tx.Bucket([]byte(section))

		// b, err := tx.CreateBucketIfNotExists(bucket)
		// Generate ID for the blog.
        // This returns an error only if the Tx is closed or not writeable.
        // That can't happen in an Update() call so I ignore the error check.
		id, err := bucket.NextSequence()
		ce(err)
		blog.ID = int(id)

		// Marshal blog data into bytes.
		encoded, err := json.Marshal(blog)
		ce(err)

		return bucket.Put(itob(blog.ID) , encoded)
	})
	ce(err)
	return err
}

func LoadBlog(DB *bolt.DB, section string,id int) Blog{
	var blog Blog
	err := DB.View(func(tx *bolt.Tx)error {
		bucket := tx.Bucket([]byte(section))
		jsonStr := bucket.Get(itob(id))
		json.Unmarshal([]byte(jsonStr), &blog)
		return nil
	})
	ce(err)
	return blog
}

func Posts(DB *bolt.DB, section string, id int) []string {
	if section == ""{
		return []string{""}
	}

	var keys []string
	err := DB.View(func(tx *bolt.Tx) error{
		bucket := tx.Bucket([]byte(section))
		if bucket == nil {
			return fmt.Errorf("No such bucket")
		}
		bucket.ForEach(func(id []byte, _ []byte) error {
			keys = append(keys, string(id))
			return nil
		})
		return nil
	})
	ce(err)
	return keys
}

func itob(v int) []byte {
    b := make([]byte, 8)
    binary.BigEndian.PutUint64(b, uint64(v))
    return b
}
