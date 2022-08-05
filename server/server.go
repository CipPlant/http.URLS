package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/google/uuid"

	_ "github.com/lib/pq"
)

var (
	db *sql.DB
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

var ID = uuid.New()

func main() {

	var err error
	connStr := "user=postgres password=200112 dbname=forURL sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(db.Stats().OpenConnections)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handlerF)
	mux.HandleFunc("/query", handlerF1)

	http.ListenAndServe("localhost:8080", mux)
}

type URLse struct {
	fullURL string
	ID      string
}

func handlerF1(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HELLO FROM HF1")
	switch {
	case r.Method == http.MethodGet:
		q := r.URL.Query().Get("id")
		fmt.Println(q)
		rows, err := db.Query("SELECT fullURL, ID FROM URLS")
		if err != nil {
			log.Fatal(err)
		}
		var URLs = URLse{}
		for rows.Next() {
			err = rows.Scan(&URLs.fullURL, &URLs.ID)
			if err != nil {
				log.Fatal(err)
			}

			if q == URLs.ID {
				w.Header().Add("Location", URLs.fullURL)
				w.WriteHeader(http.StatusTemporaryRedirect)
				return
			}
		}
		http.Error(w, "No such get URL", http.StatusNotFound)
	default:
		http.Error(w, "just for GET requests", http.StatusUnauthorized)
	}

}
func handlerF(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HELLO FROM hF")
	switch {
	case r.Method == http.MethodPost:
		rows, err := db.Query("SELECT fullURL, ID FROM URLS")

		if err != nil {
			log.Fatal(err)
		}
		var URLs = URLse{}

		for rows.Next() {
			err = rows.Scan(&URLs.fullURL, &URLs.ID)
			if err != nil {
				log.Fatal(err)
			}
			if URLs.fullURL == r.Header.Get("URL") {
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte("http://localhost:8080/" + URLs.ID))
				return
			}
		}

		res, err := db.Exec("INSERT INTO URLS (fullURL, ID) VALUES ($1, $2)",
			r.Header.Get("URL"),
			RandStringBytes(10),
		)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(res)

		rows1, err := db.Query("SELECT ID FROM URLS")
		for rows1.Next() {
			err = rows1.Scan(&URLs.ID)
		}
		if err != nil {
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("http://localhost:8080/" + URLs.ID))

	default:
		http.Error(w, "Only GET requests are allowed!", http.StatusUnauthorized)
	}
}
