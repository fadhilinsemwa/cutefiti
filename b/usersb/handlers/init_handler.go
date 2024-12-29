package handlers

import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// InitHandler structure
type InitHandler struct {
	DB     *sql.DB
	DBType string
}

// NewInitHandler creates a new InitHandler
func NewInitHandler(db *sql.DB, dbType string) *InitHandler {
	return &InitHandler{DB: db, DBType: dbType}
}

// InitializeDatabase inserts initial data into the database.
func (handler *InitHandler) InitializeDatabase() error {
	// Check if user already exists
	var userID int
	err := handler.DB.QueryRow(`SELECT user_id FROM users WHERE username = ?`, "fadhilisn").Scan(&userID)

	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error checking existing user: %v", err)
		return err
	}

	// If user does not exist, insert user
	if userID == 0 {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte("1234"), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			return err
		}

		_, err = handler.DB.Exec(`INSERT INTO users (username, email, firstname, lastname, city, password) VALUES (?, ?, ?, ?, ?, ?)`,
			"fadhilisn", "fadhilinsemwa@gmail.com", "Fadhili", "Nsemwa", "Mwanza", passwordHash)
		if err != nil {
			log.Printf("Error inserting user: %v", err)
			return err
		}

		log.Println("User 'fadhilisn' inserted successfully.")
	} else {
		log.Println("User 'fadhilisn' already exists. Skipping insertion.")
	}

	// Insert role if it does not exist
	var roleID int
	err = handler.DB.QueryRow(`SELECT id FROM roles WHERE role_id = ?`, "2").Scan(&roleID)

	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error checking existing role: %v", err)
		return err
	}

	// If role does not exist, insert role
	if roleID == 0 {
		_, err = handler.DB.Exec(`INSERT INTO roles (role_id, role_name) VALUES (?, ?)`, "2", "SuperAdmin")
		if err != nil {
			log.Printf("Error inserting role: %v", err)
			return err
		}
		log.Println("Role 'SuperAdmin' inserted successfully.")
	} else {
		log.Println("Role 'SuperAdmin' already exists. Skipping insertion.")
	}

	// Retrieve user ID for 'fadhilisn' again for the user_roles table insertion
	err = handler.DB.QueryRow(`SELECT user_id FROM users WHERE username = ?`, "fadhilisn").Scan(&userID)
	if err != nil {
		log.Printf("Error retrieving user ID for user 'fadhilisn': %v", err)
		return err
	}

	// Insert into user_roles table
	_, err = handler.DB.Exec(`INSERT OR IGNORE INTO user_roles (user_id, role_id) VALUES (?, ?)`, userID, "2")
	if err != nil {
		log.Printf("Error inserting into user_roles table: %v", err)
		return err
	}
	log.Println("User role mapping inserted successfully.")

	// Define components to be inserted
	components := []struct {
		Name string
		Path string
	}{
		{"Settings", "src/pages/Settings.tsx"},
		{"Roles", "/src/components/Roles.tsx"},
		{"Roles Components", "/src/components/RolesComponents.tsx"},
		{"Users", "/src/components/Users.tsx"},
	}

	// Loop through components and insert them into the database if they don't exist
	for _, component := range components {
		var componentID int
		err = handler.DB.QueryRow(`SELECT id FROM components WHERE component_name = ?`, component.Name).Scan(&componentID)

		if err != nil && err != sql.ErrNoRows {
			log.Printf("Error checking existing component '%s': %v", component.Name, err)
			return err
		}

		// If component does not exist, insert component
		if componentID == 0 {
			_, err = handler.DB.Exec(`INSERT INTO components (component_name, component_path) VALUES (?, ?)`,
				component.Name, component.Path)
			if err != nil {
				log.Printf("Error inserting component '%s': %v", component.Name, err)
				return err
			}
			log.Printf("Component '%s' inserted successfully.", component.Name)
		} else {
			log.Printf("Component '%s' already exists. Skipping insertion.", component.Name)
		}

		// Retrieve component ID for roles_components table insertion
		err = handler.DB.QueryRow(`SELECT id FROM components WHERE component_name = ?`, component.Name).Scan(&componentID)
		if err != nil {
			log.Printf("Error retrieving component ID for component '%s': %v", component.Name, err)
			return err
		}

		// Insert into roles_components table
		_, err = handler.DB.Exec(`INSERT OR IGNORE INTO roles_components (role_id, component_id) VALUES (?, ?)`, "2", componentID)
		if err != nil {
			log.Printf("Error inserting role-component mapping for '%s': %v", component.Name, err)
			return err
		}
		log.Printf("Role-component mapping for '%s' inserted successfully.", component.Name)
	}

	return nil
}
