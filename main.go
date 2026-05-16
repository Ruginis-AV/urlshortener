package main

import (
	"URL/internal/handler"
	"URL/internal/service"
	"URL/internal/storage"
	"embed"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/joho/godotenv"
)

//go:embed migrations
var migrations embed.FS

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	url := os.Getenv("DB_URL")
	port := os.Getenv("PORT")

	// _, filename, _, _ := runtime.Caller(0)
	// root := filepath.Join(filepath.Dir(filename), "..")
	// migPath := "file:///" + filepath.ToSlash(filepath.Join(root, "migrations"))

	// log.Println(migPath)

	// m, err := migrate.New(migPath, url)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	d, err := iofs.New(migrations, "migrations")
	if err != nil {
		log.Fatal(err)
	}
	m, err := migrate.NewWithSourceInstance("iofs", d, url)
	if err != nil {
		log.Fatal(err)
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	db := storage.New(url)
	srv := service.New(db)
	handle := handler.New(srv)

	r := chi.NewRouter()
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	r.Post("/shorten", handle.ShortenUrl)
	r.Get("/{code}", handle.Redirect)
	r.Get("/stats/{code}", handle.GetStats)
	r.Delete("/{code}", handle.Delete)

	err = http.ListenAndServe(port, r)
	if err != nil {
		log.Fatal(err)
	}
}
