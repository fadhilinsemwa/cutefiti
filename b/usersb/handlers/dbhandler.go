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
		dbType = "sqlite" // Default to SQLite if not set in .env
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "./users.db" // Default value if not set in .env
	}

	var dsn string
	var db *sql.DB
	var err error

	switch dbType {
	case "mysql":
		// For MySQL, format the DSN as user:password@tcp(host:port)/dbname
		dsn = os.Getenv("MYSQL_DSN")
		if dsn == "" {
			dsn = "user:password@tcp(localhost:3306)/dbname"
		}
		// Open MySQL connection
		db, err = sql.Open(dbType, dsn)

	case "sqlite":
		// For SQLite, use the dbName as the DSN
		dsn = dbName
		// Open SQLite connection explicitly with "sqlite3" driver
		db, err = sql.Open("sqlite3", dsn)

	default:
		return nil, errors.New("unsupported database type")
	}

	// Debug logging before testing the connection
	log.Printf("Attempting to open database connection:")
	log.Printf("Database Type: %s", dbType)
	log.Printf("DSN: %s", dsn)

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

	// Enable foreign key constraints for SQLite
	if dbType == "sqlite" {
		_, err = db.Exec("PRAGMA foreign_keys = ON;")
		if err != nil {
			log.Printf("Error setting foreign key constraints: %v", err)
			return nil, err
		}
	}

	log.Println("Connected to", dbType, "database successfully!")

	// Create the tables
	tableHandler := NewTableHandler(db, dbType)

	// Create users table
	err = tableHandler.CreateUsersTable()
	if err != nil {
		log.Printf("Error creating users table: %v", err)
		return nil, err
	}

	// Create roles and user_roles tables
	err = tableHandler.CreateRolesTable()
	if err != nil {
		log.Printf("Error creating roles table: %v", err)
		return nil, err
	}

	// Create components table
	err = tableHandler.CreateComponentsTable()
	if err != nil {
		log.Printf("Error creating components table: %v", err)
		return nil, err
	}

	// Create roles_components table
	err = tableHandler.CreateRolesComponentsTable()
	if err != nil {
		log.Printf("Error creating roles_components table: %v", err)
		return nil, err
	}

	log.Println("All tables created successfully.")

	// Initialize database with default values
	initHandler := NewInitHandler(db, dbType)
	err = initHandler.InitializeDatabase()
	if err != nil {
		log.Printf("Error initializing database: %v", err)
		return nil, err
	}

	log.Println("Database initialized with default values successfully.")

	return &DBHandler{DB: db, DBType: dbType}, nil
}

// Close the database connection
func (handler *DBHandler) CloseDB() {
	if handler.DB != nil {
		handler.DB.Close()
		log.Println("Database connection closed.")
	}
}
