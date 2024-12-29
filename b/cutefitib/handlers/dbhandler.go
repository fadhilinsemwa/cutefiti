package handlers

import (
	"database/sql"
	"errors"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql" // Import MySQL driver
	_ "github.com/mattn/go-sqlite3"    // Import SQLite driver
)

// DBHandler structure
type DBHandler struct {
	DB     *sql.DB
	DBType string
}

// Initialize the database connection
func NewDBHandler() (*DBHandler, error) {
	// Get the database type and name from environment variables
	dbType := os.Getenv("DB_TYPE")
	if dbType == "" {
		dbType = "sqlite" // Default to SQLite if not set
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "./cutefiti.db" // Default to SQLite file if not set
	}

	var dsn string
	var db *sql.DB
	var err error

	switch dbType {
	case "mysql":
		// MySQL DSN format: user:password@tcp(host:port)/dbname
		user := os.Getenv("MYSQL_USER")
		password := os.Getenv("MYSQL_PASSWORD")
		if user == "" || password == "" {
			return nil, errors.New("MySQL username or password is not set")
		}
		dsn = user + ":" + password + "@tcp(localhost:3306)/" + dbName
		db, err = sql.Open("mysql", dsn)
	case "sqlite":
		dsn = dbName
		db, err = sql.Open("sqlite3", dsn)
	default:
		return nil, errors.New("unsupported database type")
	}

	if err != nil {
		log.Printf("Error opening database connection: %v", err)
		return nil, err
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Printf("Error pinging database: %v", err)
		return nil, err
	}

	if dbType == "sqlite" {
		_, err = db.Exec("PRAGMA foreign_keys = ON;")
		if err != nil {
			log.Printf("Error enabling foreign keys: %v", err)
			return nil, err
		}
	}

	log.Printf("Connected to %s database successfully!", dbType)

	// Initialize table creation
	tableHandler := NewTableHandler(db, dbType)
	err = tableHandler.CreateAllTables()
	if err != nil {
		log.Printf("Error creating tables: %v", err)
		return nil, err
	}

	return &DBHandler{DB: db, DBType: dbType}, nil
}

// Close the database connection
func (handler *DBHandler) CloseDB() {
	if handler.DB != nil {
		handler.DB.Close()
		log.Println("Database connection closed.")
	}
}
