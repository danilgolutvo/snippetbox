package main

import (
	"database/sql"
	"flag"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"snippetboxmod/internal/models"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *models.SnippetModel
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func main() {

	addr := flag.String("addr", ":4000", "http service address")
	flag.Parse()

	dsn := flag.String("postgres", "postgres://web:pass@localhost:5432/snippetbox?sslmode=disable", "MySQL data source name")

	infolog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorlog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorlog.Fatalln(err)
	}

	defer db.Close()

	app := &application{
		errorLog: errorlog,
		infoLog:  infolog,
		snippets: &models.SnippetModel{DB: db},
	}

	srv := &http.Server{
		Addr:     *addr,
		Handler:  app.routes(),
		ErrorLog: errorlog,
	}

	infolog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorlog.Fatal(err)
}
