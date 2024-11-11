package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	db, err := sql.Open("pgx", "host=localhost port=5432 user=root password=secret dbname=test_db sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected!")

	tpl, err := template.ParseFiles("test.gohtml")
	if err != nil {
		panic(err)
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		err := tpl.Execute(w, nil)
		if err != nil {
			log.Printf("executing template: %v", err)
			http.Error(w, "There was an error executing the template", http.StatusInternalServerError)
			return
		}
	})

	fmt.Println("Starting the server on :3000...")

	http.ListenAndServe(":3000", r)

}
