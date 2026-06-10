package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/ncruces/go-sqlite3/driver"

	"uleam-transporte/internal/handlers"
	"uleam-transporte/internal/storage"
)

func main() {

	addr := flag.String("addr", ":8080", "direccion del servidor")
	dbPath := flag.String("db", "./data/transporte.db", "ruta a la base de datos SQLite")
	flag.Parse()

	almacen := initSQLite(*dbPath)

	srv := handlers.NewServer(almacen)

	log.Printf("servidor iniciado en %s", *addr)

	if err := http.ListenAndServe(*addr, srv.RegisterRoutes()); err != nil {
		log.Fatalf("error iniciando servidor: %v", err)
	}
}

func initSQLite(dbPath string) storage.Almacen {

	if err := os.MkdirAll("./data", 0755); err != nil {
		log.Fatalf("error creando directorio de datos: %v", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("error abriendo base de datos: %v", err)
	}

	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		log.Fatalf("error activando llaves foraneas: %v", err)
	}
	if _, err := db.Exec("PRAGMA journal_mode = WAL"); err != nil {
		log.Fatalf("error activando WAL: %v", err)
	}

	s := storage.NewSQLiteAlmacen(db)

	if err := s.Migrate(); err != nil {
		log.Fatalf("error ejecutando migraciones: %v", err)
	}

	if err := s.Seed(); err != nil {
		log.Fatalf("error sembrando datos: %v", err)
	}

	log.Printf("base de datos SQLite inicializada en %s", dbPath)
	return s
}
