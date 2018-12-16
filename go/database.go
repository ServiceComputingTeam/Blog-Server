package swagger

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

const (
	DBPATH = "./data/blog.db"
)

func init() {
	// if _, err := os.Stat("./data"); os.IsNotExist(err) {
	// 	err = os.Mkdir("./data", 0777)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
	err := os.MkdirAll("./data", 0777)
	if err != nil {
		log.Fatal(err)
	}
	db, err := bolt.Open(DBPATH, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	db.Batch(func(tx *bolt.Tx) error {
		var err error
		_, err = tx.CreateBucketIfNotExists([]byte("Users"))
		_, err = tx.CreateBucketIfNotExists([]byte("Blogs"))
		if err != nil {
			return err
		}
		return nil
	})

	defer db.Close()
}

func DBcheakUsernameAndPassWord(username, password string) (bool, error) {
	user, err := DBgetUserByUsername(username)
	if user == nil || err != nil {
		return false, err
	}
	return user.Password == password, err
}

func DBgetBlogByAuthor(username string) (*[]Blog, error) {
	user, err := DBgetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	if user == nil || user.BlogId == nil {
		return nil, errors.New("user inexits or user has not blog")
	}
	blogs := make([]Blog, len(user.BlogId))
	db, err := bolt.Open(DBPATH, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Blogs"))
		for index, id := range user.BlogId {
			var temp Blog
			err := json.Unmarshal(b.Get(itob(id)), &temp)
			if err != nil {
				return err
			}
			blogs[index] = temp
		}
		return nil
	})

	return &blogs, err
}

func DBgetAllBlog(offset, page int) (*[]Blog, error) {
	blogs := make([]Blog, page)
	db, err := bolt.Open(DBPATH, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Blogs"))
		c := b.Cursor()
		min := itob(uint64(offset))
		max := itob(uint64(offset + page - 1))

		var index = 0
		for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
			var temp Blog
			err := json.Unmarshal(v, &temp)
			if err != nil {
				return err
			}
			blogs[index] = temp
			index++
		}
		blogs = blogs[0:index]
		return nil
	})
	return &blogs, err
}

func DBCreateUser(user *User) (*User, error) {
	db, err := bolt.Open(DBPATH, 0600, nil)
	if err != nil {
		log.Fatal(err)
		return user, err
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Users"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var temp User
			err := json.Unmarshal(v, &temp)
			if err != nil {
				return err
			}
			if user.Username == temp.Username {
				return errors.New("username exits")
			}
		}
		return nil
	})

	if err != nil {
		return user, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Users"))
		id, _ := b.NextSequence()
		user.Id = id
		buf, err := json.Marshal(user)
		if err != nil {
			return err
		}
		return b.Put(itob(user.Id), buf)
	})
	return user, err
}

func DBCreateBlog(blog *Blog) (*Blog, error) {
	temp, err := DBgetBolgByBlogTitleAndAuthor(blog.Title, blog.Owner)
	if temp != nil || err != nil {
		return nil, errors.New("blog exits")
	}

	db, err := bolt.Open(DBPATH, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Blogs"))
		id, _ := b.NextSequence()
		blog.Id = id
		buf, err := json.Marshal(blog)
		if err != nil {
			return err
		}
		err2 := b.Put(itob(blog.Id), buf)
		if err2 != nil {
			return err2
		}

		b2 := tx.Bucket([]byte("Users"))
		c2 := b2.Cursor()

		for k, v := c2.First(); k != nil; k, v = c2.Next() {
			var temp User
			json.Unmarshal(v, &temp)
			if temp.Username == blog.Owner {
				temp.BlogId = append(temp.BlogId, blog.Id)
				buf, err := json.Marshal(temp)
				if err != nil {
					return err
				}
				return b2.Put(k, buf)
			}
		}
		return nil

	})
	return blog, err
}

func DBCreateReview(review *Review) (*Review, error) {
	blog, err := DBgetBolgByBlogTitleAndAuthor(review.Blogtitle, review.Blogowner)
	if err != nil || blog == nil {
		return nil, err
	}

	review.Id = uint64(len(blog.Review))
	blog.Review = append(blog.Review, *review)
	db, err := bolt.Open(DBPATH, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	return review, db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Blogs"))
		buf, err := json.Marshal(blog)
		if err != nil {
			return err
		}
		return b.Put(itob(blog.Id), buf)
	})
}

func DBgetReviewByBlogTitleAndAuthor(title, author string) (*[]Review, error) {
	blog, err := DBgetBolgByBlogTitleAndAuthor(title, author)
	return &blog.Review, err
}

func DBgetBlogsByLabelname(labelname string) (*[]Blog, error) {
	var blogs []Blog
	db, err := bolt.Open(DBPATH, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Blogs"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var temp Blog
			err := json.Unmarshal(v, &temp)
			if err != nil {
				return err
			}
			for _, v := range temp.Label {
				if v == labelname {
					blogs = append(blogs, temp)
					break
				}
			}
		}
		return nil
	})

	return &blogs, err
}

func DBgetBolgByBlogTitle(title string) (*[]Blog, error) {
	db, err := bolt.Open(DBPATH, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var blogs []Blog
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Blogs"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var temp Blog
			json.Unmarshal(v, &temp)
			if temp.Title == title {
				blogs = append(blogs, temp)
			}
		}
		return nil
	})
	return &blogs, err
}

func DBgetBolgByBlogTitleAndAuthor(title, author string) (*Blog, error) {
	db, err := bolt.Open(DBPATH, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var blog *Blog
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Blogs"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var temp Blog
			err := json.Unmarshal(v, &temp)
			if err != nil {
				return err
			}
			if temp.Title == title && temp.Owner == author {
				fmt.Println(temp)
				blog = &temp
				break
			}
		}
		return nil
	})
	return blog, err
}

func DBgetUserByUsername(username string) (*User, error) {
	db, err := bolt.Open(DBPATH, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var user *User
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Users"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var temp User
			json.Unmarshal(v, &temp)
			if temp.Username == username {
				user = &temp
				break
			}
		}
		return nil
	})
	return user, err

}

func btoint64(v []byte) uint64 {
	bits := binary.BigEndian.Uint64(v)
	return bits
}

func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
