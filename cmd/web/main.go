package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
	"os"
	"snippetboxmod/internal/models"
	"time"
)

type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	snippets       models.SnippetModelInterface
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
	users          models.UserModelInterface
	debug          bool
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

	debug := flag.Bool("debug", false, "debug mode")

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
	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	app := &application{
		errorLog:       errorlog,
		infoLog:        infolog,
		snippets:       &models.SnippetModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
		users:          &models.UserModel{DB: db},
		debug:          *debug,
	}
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:         *addr,
		Handler:      app.routes(),
		ErrorLog:     errorlog,
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infolog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorlog.Fatal(err)
}
