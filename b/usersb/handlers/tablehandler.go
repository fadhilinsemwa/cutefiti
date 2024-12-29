package handlers

import (
	"database/sql"
	"errors"
	"log"
)

// TableHandler structure
type TableHandler struct {
	DB     *sql.DB
	DBType string
}

// NewTableHandler creates a new TableHandler
func NewTableHandler(db *sql.DB, dbType string) *TableHandler {
	return &TableHandler{DB: db, DBType: dbType}
}

// CreateUsersTable creates the users table based on the provided schema.
func (handler *TableHandler) CreateUsersTable() error {
	var createTableSQL string

	switch handler.DBType {
	case "sqlite":
		createTableSQL = `
		CREATE TABLE IF NOT EXISTS users (
			user_id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			email TEXT NOT NULL UNIQUE,
			firstname TEXT NOT NULL,
			lastname TEXT NOT NULL,
			city TEXT NOT NULL,
			password TEXT NOT NULL,
			reset_token TEXT,
			email_verified BOOLEAN DEFAULT 0,
			verification_token TEXT
		);`
	case "mysql":
		createTableSQL = `
		CREATE TABLE IF NOT EXISTS users (
			user_id INT AUTO_INCREMENT PRIMARY KEY,
			username VARCHAR(255) NOT NULL UNIQUE,
			email VARCHAR(255) NOT NULL UNIQUE,
			firstname VARCHAR(255) NOT NULL,
			lastname VARCHAR(255) NOT NULL,
			city VARCHAR(255) NOT NULL,
			password VARCHAR(255) NOT NULL,
			reset_token VARCHAR(255),
			email_verified BOOLEAN DEFAULT 0,
			verification_token VARCHAR(255)
		);`
	default:
		return errors.New("unsupported database type for table creation")
	}

	_, err := handler.DB.Exec(createTableSQL)
	if err != nil {
		return err
	}

	log.Println("Users table created or updated successfully.")
	return nil
}

// CreateRolesTable creates the roles table and the user_roles table to map users to roles.
func (handler *TableHandler) CreateRolesTable() error {
	var createRolesTableSQL string
	var createUserRolesTableSQL string

	switch handler.DBType {
	case "sqlite":
		createRolesTableSQL = `
		CREATE TABLE IF NOT EXISTS roles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			role_id TEXT NOT NULL UNIQUE,
			role_name TEXT NOT NULL UNIQUE
		);`

		createUserRolesTableSQL = `
		CREATE TABLE IF NOT EXISTS user_roles (
			user_id INTEGER NOT NULL,
			role_id TEXT NOT NULL,
			FOREIGN KEY(user_id) REFERENCES users(user_id),
			FOREIGN KEY(role_id) REFERENCES roles(role_id),
			PRIMARY KEY (user_id, role_id)
		);`
	case "mysql":
		createRolesTableSQL = `
		CREATE TABLE IF NOT EXISTS roles (
			id INT AUTO_INCREMENT PRIMARY KEY,
			role_id VARCHAR(255) NOT NULL UNIQUE,
			role_name VARCHAR(255) NOT NULL UNIQUE
		);`

		createUserRolesTableSQL = `
		CREATE TABLE IF NOT EXISTS user_roles (
			user_id INT NOT NULL,
			role_id VARCHAR(255) NOT NULL,
			FOREIGN KEY(user_id) REFERENCES users(user_id),
			FOREIGN KEY(role_id) REFERENCES roles(role_id),
			PRIMARY KEY (user_id, role_id)
		);`
	default:
		return errors.New("unsupported database type for table creation")
	}

	// Create roles table
	_, err := handler.DB.Exec(createRolesTableSQL)
	if err != nil {
		return err
	}

	// Create user_roles table
	_, err = handler.DB.Exec(createUserRolesTableSQL)
	if err != nil {
		return err
	}

	// Insert default role if the roles table is created for the first time
	_, err = handler.DB.Exec(`INSERT OR IGNORE INTO roles (role_id, role_name) VALUES ('1', 'Admin');`)
	if err != nil {
		return err
	}

	log.Println("Roles and UserRoles tables created or updated successfully, and default role 'Admin' inserted.")
	return nil
}

// CreateComponentsTable creates the components table.
func (handler *TableHandler) CreateComponentsTable() error {
	var createComponentsTableSQL string

	switch handler.DBType {
	case "sqlite":
		createComponentsTableSQL = `
		CREATE TABLE IF NOT EXISTS components (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			component_name TEXT NOT NULL,
			component_path TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`
	case "mysql":
		createComponentsTableSQL = `
		CREATE TABLE IF NOT EXISTS components (
			id INT AUTO_INCREMENT PRIMARY KEY,
			component_name VARCHAR(255) NOT NULL,
			component_path VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		);`
	default:
		return errors.New("unsupported database type for table creation")
	}

	_, err := handler.DB.Exec(createComponentsTableSQL)
	if err != nil {
		return err
	}

	log.Println("Components table created or updated successfully.")
	return nil
}

// CreateRolesComponentsTable creates the roles_components table.
func (handler *TableHandler) CreateRolesComponentsTable() error {
	var createRolesComponentsTableSQL string

	switch handler.DBType {
	case "sqlite":
		createRolesComponentsTableSQL = `
		CREATE TABLE IF NOT EXISTS roles_components (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			role_id TEXT NOT NULL,
			component_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(role_id) REFERENCES roles(role_id),  -- Changed to role_id
			FOREIGN KEY(component_id) REFERENCES components(id)
		);`
	case "mysql":
		createRolesComponentsTableSQL = `
		CREATE TABLE IF NOT EXISTS roles_components (
			id INT AUTO_INCREMENT PRIMARY KEY,
			role_id VARCHAR(255) NOT NULL,
			component_id INT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY(role_id) REFERENCES roles(role_id),  -- Changed to role_id
			FOREIGN KEY(component_id) REFERENCES components(id)
		);`
	default:
		return errors.New("unsupported database type for table creation")
	}

	_, err := handler.DB.Exec(createRolesComponentsTableSQL)
	if err != nil {
		return err
	}

	log.Println("RolesComponents table created or updated successfully.")
	return nil
}
