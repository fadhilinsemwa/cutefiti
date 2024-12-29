package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

// RoleHandler structure
type RoleHandler struct {
	DB     *sql.DB
	DBType string
}

// NewRoleHandler creates a new RoleHandler
func NewRoleHandler(db *sql.DB, dbType string) *RoleHandler {
	return &RoleHandler{DB: db, DBType: dbType}
}

// CreateDefaultRoles ensures that the default Admin role with role_id 'admin' is present.
func (handler *RoleHandler) CreateDefaultRoles() error {
	var insertAdminSQL string

	switch handler.DBType {
	case "sqlite":
		insertAdminSQL = `
		INSERT OR IGNORE INTO roles (role_id, role_name)
		VALUES ('admin', 'Admin');`
	case "mysql":
		insertAdminSQL = `
		INSERT IGNORE INTO roles (role_id, role_name)
		VALUES ('admin', 'Admin');`
	default:
		return errors.New("unsupported database type for inserting default roles")
	}

	_, err := handler.DB.Exec(insertAdminSQL)
	if err != nil {
		return err
	}

	log.Println("Default Admin role inserted or already exists.")
	return nil
}

// CreateRole creates a new role with a unique role_id and role_name.
func (handler *RoleHandler) CreateRole(roleID, roleName string) error {
	var createRoleSQL string

	switch handler.DBType {
	case "sqlite":
		createRoleSQL = `
		INSERT INTO roles (role_id, role_name)
		VALUES (?, ?);`
	case "mysql":
		createRoleSQL = `
		INSERT INTO roles (role_id, role_name)
		VALUES (?, ?);`
	default:
		return errors.New("unsupported database type for creating roles")
	}

	_, err := handler.DB.Exec(createRoleSQL, roleID, roleName)
	if err != nil {
		return err
	}

	log.Printf("Role '%s' created successfully with role_id '%s'.", roleName, roleID)
	return nil
}

// DeleteRole deletes a role based on id (the primary key).
func (handler *RoleHandler) DeleteRole(id int) error {
	var deleteRoleSQL string

	switch handler.DBType {
	case "sqlite":
		deleteRoleSQL = `DELETE FROM roles WHERE id = ?;`
	case "mysql":
		deleteRoleSQL = `DELETE FROM roles WHERE id = ?;`
	default:
		return errors.New("unsupported database type for deleting roles")
	}

	result, err := handler.DB.Exec(deleteRoleSQL, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no role found with the given id")
	}

	log.Printf("Role with id '%d' deleted successfully.", id)
	return nil
}

// GetRoleByID retrieves a role based on id (the primary key).
func (handler *RoleHandler) GetRoleByID(id int) (string, string, error) {
	var getRoleSQL string
	var roleID, roleName string

	switch handler.DBType {
	case "sqlite":
		getRoleSQL = `SELECT role_id, role_name FROM roles WHERE id = ?;`
	case "mysql":
		getRoleSQL = `SELECT role_id, role_name FROM roles WHERE id = ?;`
	default:
		return "", "", errors.New("unsupported database type for retrieving role by ID")
	}

	err := handler.DB.QueryRow(getRoleSQL, id).Scan(&roleID, &roleName)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", errors.New("no role found with the given id")
		}
		return "", "", err
	}

	log.Printf("Role retrieved: id = '%d', role_id = '%s', role_name = '%s'", id, roleID, roleName)
	return roleID, roleName, nil
}

// GetAllRoles retrieves all roles in the system.
func (handler *RoleHandler) GetAllRoles() (map[int]map[string]string, error) {
	var getAllRolesSQL string
	roles := make(map[int]map[string]string)

	switch handler.DBType {
	case "sqlite":
		getAllRolesSQL = `SELECT id, role_id, role_name FROM roles;`
	case "mysql":
		getAllRolesSQL = `SELECT id, role_id, role_name FROM roles;`
	default:
		return nil, errors.New("unsupported database type for retrieving all roles")
	}

	rows, err := handler.DB.Query(getAllRolesSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var roleID, roleName string
		if err := rows.Scan(&id, &roleID, &roleName); err != nil {
			return nil, err
		}
		roles[id] = map[string]string{
			"role_id":   roleID,
			"role_name": roleName,
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	log.Printf("Retrieved %d roles from the database.", len(roles))
	return roles, nil
}

// AssignUserRole assigns a role to a user by username and role_id.
func (handler *RoleHandler) AssignUserRole(username, roleID string) error {
	var getUserIDSQL string
	var insertUserRoleSQL string
	var userID int

	// Get the user's ID based on their username
	switch handler.DBType {
	case "sqlite":
		getUserIDSQL = `SELECT user_id FROM users WHERE username = ?;`
		insertUserRoleSQL = `
		INSERT OR IGNORE INTO user_roles (user_id, role_id)
		VALUES (?, ?);`
	case "mysql":
		getUserIDSQL = `SELECT user_id FROM users WHERE username = ?;`
		insertUserRoleSQL = `
		INSERT IGNORE INTO user_roles (user_id, role_id)
		VALUES (?, ?);`
	default:
		return errors.New("unsupported database type for assigning user roles")
	}

	// Retrieve the user ID
	err := handler.DB.QueryRow(getUserIDSQL, username).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("no user found with the given username")
		}
		return err
	}

	// Insert into the user_roles table
	_, err = handler.DB.Exec(insertUserRoleSQL, userID, roleID)
	if err != nil {
		return err
	}

	log.Printf("Assigned role_id '%s' to user '%s' (user_id: %d).", roleID, username, userID)
	return nil
}

// GetUserRoles retrieves the username, email, role_id, and role_name for all users with their roles.
func (handler *RoleHandler) GetUserRoles() (map[string]map[string]string, error) {
	var getUserRolesSQL string
	userRoles := make(map[string]map[string]string)

	switch handler.DBType {
	case "sqlite":
		getUserRolesSQL = `
		SELECT u.username, u.email, r.role_id, r.role_name 
		FROM users u 
		JOIN user_roles ur ON u.user_id = ur.user_id 
		JOIN roles r ON ur.role_id = r.role_id;`
	case "mysql":
		getUserRolesSQL = `
		SELECT u.username, u.email, r.role_id, r.role_name 
		FROM users u 
		JOIN user_roles ur ON u.user_id = ur.user_id 
		JOIN roles r ON ur.role_id = r.role_id;`
	default:
		return nil, errors.New("unsupported database type for retrieving user roles")
	}

	rows, err := handler.DB.Query(getUserRolesSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var username, email, roleID, roleName string
		if err := rows.Scan(&username, &email, &roleID, &roleName); err != nil {
			return nil, err
		}
		userRoles[username] = map[string]string{
			"email":     email,
			"role_id":   roleID,
			"role_name": roleName,
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	log.Printf("Retrieved roles for %d users from the database.", len(userRoles))
	return userRoles, nil
}

// HandleCreateRole handles the role creation request
func (handler *RoleHandler) HandleCreateRole(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		RoleID   string `json:"role_id"`
		RoleName string `json:"role_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := handler.CreateRole(reqBody.RoleID, reqBody.RoleName)
	if err != nil {
		http.Error(w, "Error creating role: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Role created successfully"})
}

// HandleGetAllRoles handles the request to retrieve all roles
func (handler *RoleHandler) HandleGetAllRoles(w http.ResponseWriter, r *http.Request) {
	roles, err := handler.GetAllRoles()
	if err != nil {
		http.Error(w, "Error retrieving roles: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(roles)
}

// HandleAssignUserRole handles the assignment of a role to a user
func (handler *RoleHandler) HandleAssignUserRole(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Username string `json:"username"`
		RoleID   string `json:"role_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := handler.AssignUserRole(reqBody.Username, reqBody.RoleID)
	if err != nil {
		http.Error(w, "Error assigning role: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Role assigned successfully"})
}

// HandleGetUserRoles handles the request to retrieve users and their roles
func (handler *RoleHandler) HandleGetUserRoles(w http.ResponseWriter, r *http.Request) {
	userRoles, err := handler.GetUserRoles()
	if err != nil {
		http.Error(w, "Error retrieving user roles: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userRoles)
}
