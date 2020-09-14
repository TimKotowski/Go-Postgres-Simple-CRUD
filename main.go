package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"net/http"
)

// the ordering matters and name of the fields matter

type Book struct {
	isbn   string
	title  string
	author string
	price  float32
}

// package level scope
var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", "postgres://postgres:password@localhost/bookstore?sslmode=disable")
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("you're connected to the database")
}

func main() {
	http.HandleFunc("/books", booksIndex)
	http.ListenAndServe(":8080", nil)
}

func booksIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(500), 500)
	}

	selectBooks := "SELECT * FROM books;"
	rows, err := db.Query(selectBooks)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	bks := make([]Book, 0)
	for rows.Next() {
		// create a Book
		// And a book is a indentifer a value of type Book our struct and i'm using an empty composite literal to do that
		bk := Book{}
		err := rows.Scan(&bk.isbn, &bk.title, &bk.author, &bk.price) // order matters
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		bks = append(bks, bk)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	for _, bk := range bks {
		fmt.Printf("%s, %s, %s, $%.2f\n", bk.isbn, bk.title, bk.author, bk.price)
	}
}
