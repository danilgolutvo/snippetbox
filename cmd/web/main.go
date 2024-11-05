package main

import (
	"database/sql"
	"flag"
	"github.com/go-playground/form/v4"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
	"os"
	"snippetboxmod/internal/models"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
	formDecoder   *form.Decoder
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

	dsn := flag.String("postgres", "postgres://postgres:mysecretpassword@localhost:5432/postgres?sslmode=disable", "MySQL data source name")

	infolog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorlog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorlog.Fatalln(err)
	}

	defer db.Close()

	formDecoder := form.NewDecoder()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorlog.Fatal(err)
	}

	app := &application{
		errorLog:      errorlog,
		infoLog:       infolog,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
		formDecoder:   formDecoder,
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
