package db

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	var err error

	for i := 0; i < 5; i++ {
		DB, err = sql.Open("postgres", databaseURL)
		if err != nil {
			log.Println("Retrying DB connection...")
			time.Sleep(2 * time.Second)
			continue
		}

		err = DB.Ping()
		if err == nil {
			log.Println("Connected to Neon PostgreSQL")
			break
		}

		log.Println("Waiting for DB to be ready...")
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("Could not connect to database:", err)
	}

	createTables()
}

func createTables() {

	// ------------------------
	// PRODUCTS TABLE
	// ------------------------
	createProductsTable := `
	CREATE TABLE IF NOT EXISTS products (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    price DOUBLE PRECISION,
    mood TEXT,
    category TEXT,
    image TEXT,
    rating DOUBLE PRECISION DEFAULT 0,
    featured BOOLEAN DEFAULT FALSE,
    stock INT
);`
	_, err := DB.Exec(createProductsTable)
	if err != nil {
		log.Fatal("Failed to create products table:", err)
	}

	log.Println("Products table ready")

	// ------------------------
	// BUNDLES TABLE
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
	// BUNDLE_PRODUCTS TABLE
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
