package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	var err error

	for i := 0; i < 30; i++ {
		DB, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Println("Retrying DB connection...")
			time.Sleep(2 * time.Second)
			continue
		}

		err = DB.Ping()
		if err == nil {
			log.Println("Connected to PostgreSQL")
			break
		}

		log.Println("Waiting for DB to be ready...")
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("Could not connect to database:", err)
	}

	// ------------------------
	// CREATE PRODUCTS TABLE
	// ------------------------
	createProductsTable := `
	CREATE TABLE IF NOT EXISTS products (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		price DOUBLE PRECISION,
		mood TEXT,
		stock INT
	);`

	_, err = DB.Exec(createProductsTable)
	if err != nil {
		log.Fatal("Failed to create products table:", err)
	}

	log.Println("Products table ready")

	// ------------------------
	// CREATE BUNDLES TABLE
	// ------------------------
	createBundlesTable := `
	CREATE TABLE IF NOT EXISTS bundles (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		mood TEXT
	);`

	_, err = DB.Exec(createBundlesTable)
	if err != nil {
		log.Fatal("Failed to create bundles table:", err)
	}

	// ------------------------
	// CREATE BUNDLE_PRODUCTS TABLE
	// ------------------------
	createBundleProductsTable := `
	CREATE TABLE IF NOT EXISTS bundle_products (
		bundle_id TEXT REFERENCES bundles(id) ON DELETE CASCADE,
		product_id TEXT REFERENCES products(id) ON DELETE CASCADE,
		quantity INT NOT NULL,
		PRIMARY KEY (bundle_id, product_id)
	);`

	_, err = DB.Exec(createBundleProductsTable)
	if err != nil {
		log.Fatal("Failed to create bundle_products table:", err)
	}

	log.Println("Bundles tables ready")
}
