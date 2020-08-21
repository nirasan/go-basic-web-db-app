package main

import (
	"log"

	"upper.io/db.v3/mysql"
)

var settings = mysql.ConnectionURL{
	Database: `booktown`,
	Host:     `127.0.0.1:13306`,
	User:     `root`,
	Password: `root`,
}

type Book struct {
	ID        uint   `db:"id"`
	Title     string `db:"title"`
	AuthorID  uint   `db:"author_id"`
	SubjectID uint   `db:"subject_id"`
}

func main() {
	sess, err := mysql.Open(settings)
	if err != nil {
		log.Fatal("Open: ", err)
	}
	defer sess.Close()

	sess.SetLogging(true)

	booksTable := sess.Collection("books")

	if _, err := booksTable.Insert(Book{
		ID:        4267,
		Title:     "book1",
		AuthorID:  0,
		SubjectID: 0,
	}); err != nil {
		log.Fatal("Insert: ", err)
	}

	// This result set includes a single item.
	res := booksTable.Find(4267)

	// The item is retrieved with the given ID.
	var book Book
	err = res.One(&book)
	if err != nil {
		log.Fatal("Find: ", err)
	}

	log.Printf("Book: %#v", book)

	// A change is made to a property.
	book.Title = "New title"

	log.Printf("Book (modified): %#v", book)

	// The result set is updated.
	if err := res.Update(book); err != nil {
		log.Printf("Update: %v\n", err)
		log.Printf("This is OK, this is a restricted sandbox with a read-only database.")
	}

	// The result set is deleted.
	if err := res.Delete(); err != nil {
		log.Printf("Delete: %v\n", err)
		log.Printf("This is OK, this is a restricted sandbox with a read-only database.")
	}
}
