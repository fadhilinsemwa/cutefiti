package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

// User struct represents the users table schema
type User struct {
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	City      string `json:"city"`
	Password  string `json:"password"`
}

// CRUDHandler structure for handling generic CRUD operations
type CRUDHandler struct {
	DB     *sql.DB
	DBType string // Added to distinguish between SQLite and MariaDB/MySQL
}

// NewCRUDHandler creates a new CRUDHandler
func NewCRUDHandler(db *sql.DB, dbType string) *CRUDHandler {
	return &CRUDHandler{DB: db, DBType: dbType}
}

// CreateUser inserts a new user into the database
func (handler *CRUDHandler) CreateUser(username, email, firstname, lastname, city, password string) (int64, error) {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	query := `INSERT INTO users (username, email, firstname, lastname, city, password) VALUES (?, ?, ?, ?, ?, ?)`
	result, err := handler.DB.Exec(query, username, email, firstname, lastname, city, string(hashedPassword))
	if err != nil {
		return 0, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return userID, nil
}

// GetUser retrieves a user by their ID
func (handler *CRUDHandler) GetUser(userID int64) (*User, error) {
	var user User
	query := `SELECT user_id, username, email, firstname, lastname, city, password FROM users WHERE user_id = ?`
	err := handler.DB.QueryRow(query, userID).Scan(&user.UserID, &user.Username, &user.Email, &user.Firstname, &user.Lastname, &user.City, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // User not found
		}
		return nil, err
	}

	return &user, nil
}

// GetAllUsers retrieves all users from the database
func (handler *CRUDHandler) GetAllUsers() ([]User, error) {
	query := `SELECT user_id, username, email, firstname, lastname, city, password FROM users`
	rows, err := handler.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.UserID, &user.Username, &user.Email, &user.Firstname, &user.Lastname, &user.City, &user.Password); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// UpdateUser updates an existing user's information
func (handler *CRUDHandler) UpdateUser(userID int64, username, email, firstname, lastname, city, password string) error {
	query := `UPDATE users SET username = ?, email = ?, firstname = ?, lastname = ?, city = ?, password = ? WHERE user_id = ?`
	_, err := handler.DB.Exec(query, username, email, firstname, lastname, city, password, userID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteUser deletes a user from the database
func (handler *CRUDHandler) DeleteUser(userID int64) error {
	query := `DELETE FROM users WHERE user_id = ?`
	_, err := handler.DB.Exec(query, userID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteAllUsers deletes all users from the database without resetting the primary key
func (handler *CRUDHandler) DeleteAllUsers() error {
	query := `DELETE FROM users`
	_, err := handler.DB.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

// TruncateUsersTable truncates the users table and resets the autoincrement
func (handler *CRUDHandler) TruncateUsersTable() error {
	var query string

	if handler.DBType == "sqlite" {
		// SQLite does not support TRUNCATE, so we use DELETE and reset autoincrement separately
		query = `DELETE FROM users`
		resetQuery := `DELETE FROM sqlite_sequence WHERE name='users'`
		_, err := handler.DB.Exec(query)
		if err != nil {
			return err
		}
		_, err = handler.DB.Exec(resetQuery)
		if err != nil {
			return err
		}
	} else {
		// In MariaDB/MySQL, we can use TRUNCATE which also resets the autoincrement
		query = `TRUNCATE TABLE users`
		_, err := handler.DB.Exec(query)
		if err != nil {
			return err
		}
	}

	return nil
}

// HandleCRUD handles all CRUD operations based on HTTP methods
func (handler *CRUDHandler) HandleCRUD(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handler.CreateUserHandler(w, r)
	case http.MethodGet:
		if r.URL.Query().Has("id") {
			handler.GetUserHandler(w, r)
		} else {
			handler.GetAllUsersHandler(w, r)
		}
	case http.MethodPut:
		handler.UpdateUserHandler(w, r)
	case http.MethodDelete:
		handler.DeleteUserHandler(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// CreateUserHandler - Handler for creating a new user
func (handler *CRUDHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	userID, err := handler.CreateUser(user.Username, user.Email, user.Firstname, user.Lastname, user.City, user.Password)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"user_id": userID})
}

// GetUserHandler - Handler for fetching a user
func (handler *CRUDHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	userIDParam := r.URL.Query().Get("id")
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := handler.GetUser(userID)
	if err != nil {
		http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// GetAllUsersHandler - Handler for fetching all users
func (handler *CRUDHandler) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := handler.GetAllUsers()
	if err != nil {
		// Log the error for debugging purposes
		log.Printf("Error fetching users: %v", err)
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

// UpdateUserHandler - Handler for updating user information
func (handler *CRUDHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	userIDParam := r.URL.Query().Get("id")
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = handler.UpdateUser(userID, user.Username, user.Email, user.Firstname, user.Lastname, user.City, user.Password)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "User updated"})
}

// DeleteUserHandler - Handler for deleting a user
func (handler *CRUDHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	userIDParam := r.URL.Query().Get("id")
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = handler.DeleteUser(userID)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "User deleted"})
}

// DeleteAllUsersHandler - Handler for deleting all users
func (handler *CRUDHandler) DeleteAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	err := handler.DeleteAllUsers()
	if err != nil {
		http.Error(w, "Failed to delete users", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "All users deleted"})
}

// TruncateUsersHandler - Handler for truncating the users table
func (handler *CRUDHandler) TruncateUsersHandler(w http.ResponseWriter, r *http.Request) {
	err := handler.TruncateUsersTable()
	if err != nil {
		http.Error(w, "Failed to truncate users table", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "Users table truncated"})
}
