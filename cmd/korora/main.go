package main

import (
	"net/http"
	"os"

	"github.com/korora-social/korora"

	sq "github.com/sleepdeprecation/squirrelly"
	_ "modernc.org/sqlite"
)

func main() {
	db, err := sq.Open("sqlite", "db/test.sqlite")
	if err != nil {
		panic(err)
	}

	// janky migration -- just run a single sql file because i don't want to figure out a migrations system yet
	rawSchema, err := os.ReadFile("db/schema.sql")
	if err != nil {
		panic(err)
	}

	_, err = db.DB.Exec(string(rawSchema))
	if err != nil {
		panic(err)
	}

	k := korora.New(db)
	http.ListenAndServe(":8800", k.Router())
}
