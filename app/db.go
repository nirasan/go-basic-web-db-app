package app

import (
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
	"upper.io/db.v3/mysql"
)

type User struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

type Book struct {
	ID     uint   `db:"id" form:"id"`
	Title  string `db:"title" form:"title"`
	Author string `db:"author" form:"author"`
	Price  int    `db:"price" form:"price"`
}

var settings = mysql.ConnectionURL{
	Database: `booktown`,
	Host:     `127.0.0.1:13307`,
	User:     `root`,
	Password: `root`,
}

type DBClient struct {
	sess  sqlbuilder.Database
	books db.Collection
	users db.Collection
}

func NewDBClient() (*DBClient, error) {
	sess, err := mysql.Open(settings)
	if err != nil {
		return nil, err
	}
	return &DBClient{
		sess:  sess,
		books: sess.Collection("books"),
		users: sess.Collection("users"),
	}, nil
}

func (c *DBClient) Close() error {
	return c.sess.Close()
}
