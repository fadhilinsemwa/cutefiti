package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// RolesComponentsHandler structure
type RolesComponentsHandler struct {
	DB     *sql.DB
	DBType string
}

// NewRolesComponentsHandler creates a new RolesComponentsHandler
func NewRolesComponentsHandler(db *sql.DB, dbType string) *RolesComponentsHandler {
	return &RolesComponentsHandler{DB: db, DBType: dbType}
}

// CreateComponent creates a new component
func (handler *RolesComponentsHandler) CreateComponent(w http.ResponseWriter, r *http.Request) {
	var component struct {
		ComponentName string `json:"component_name"`
		ComponentPath string `json:"component_path"`
	}

	err := json.NewDecoder(r.Body).Decode(&component)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := "INSERT INTO components (component_name, component_path) VALUES (?, ?)"
	_, err = handler.DB.Exec(query, component.ComponentName, component.ComponentPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "component created successfully"})
}

// GetAllComponents returns all components
func (handler *RolesComponentsHandler) GetAllComponents(w http.ResponseWriter, r *http.Request) {
	rows, err := handler.DB.Query("SELECT id, component_name, component_path, created_at, updated_at FROM components")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var components []map[string]interface{}
	for rows.Next() {
		var id, componentName, componentPath, createdAt, updatedAt interface{}

		err := rows.Scan(&id, &componentName, &componentPath, &createdAt, &updatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		component := map[string]interface{}{
			"id":             id,
			"component_name": componentName,
			"component_path": componentPath,
			"created_at":     createdAt,
			"updated_at":     updatedAt,
		}

		components = append(components, component)
	}

	json.NewEncoder(w).Encode(components)
}

// UpdateComponent updates an existing component
func (handler *RolesComponentsHandler) UpdateComponent(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var component struct {
		ComponentName string `json:"component_name"`
		ComponentPath string `json:"component_path"`
	}

	err := json.NewDecoder(r.Body).Decode(&component)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := "UPDATE components SET component_name = ?, component_path = ? WHERE id = ?"
	_, err = handler.DB.Exec(query, component.ComponentName, component.ComponentPath, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "component updated successfully"})
}

// DeleteComponent deletes a component by ID
func (handler *RolesComponentsHandler) DeleteComponent(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	query := "DELETE FROM components WHERE id = ?"
	_, err := handler.DB.Exec(query, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "component deleted successfully"})
}

// DeleteAllComponents deletes all components
func (handler *RolesComponentsHandler) DeleteAllComponents(w http.ResponseWriter, r *http.Request) {
	_, err := handler.DB.Exec("DELETE FROM components")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "all components deleted successfully"})
}

// TruncateComponents truncates the components table
func (handler *RolesComponentsHandler) TruncateComponents(w http.ResponseWriter, r *http.Request) {
	var truncateSQL string
	switch handler.DBType {
	case "sqlite":
		truncateSQL = "DELETE FROM components; DELETE FROM sqlite_sequence WHERE name = 'components';"
	case "mysql":
		truncateSQL = "TRUNCATE TABLE components"
	default:
		http.Error(w, "unsupported database type", http.StatusBadRequest)
		return
	}

	_, err := handler.DB.Exec(truncateSQL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "components table truncated successfully"})
}

// CreateRoleComponent creates an entry in the roles_components table
func (handler *RolesComponentsHandler) CreateRoleComponent(w http.ResponseWriter, r *http.Request) {
	var assignment struct {
		RoleID      int `json:"role_id"`
		ComponentID int `json:"component_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&assignment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := "INSERT INTO roles_components (role_id, component_id) VALUES (?, ?)"
	_, err = handler.DB.Exec(query, assignment.RoleID, assignment.ComponentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "role assigned to component successfully"})
}

// GetAllRolesComponents returns all role-component mappings
func (handler *RolesComponentsHandler) GetAllRolesComponents(w http.ResponseWriter, r *http.Request) {
	rows, err := handler.DB.Query("SELECT id, role_id, component_id, created_at, updated_at FROM roles_components")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var mappings []map[string]interface{}
	for rows.Next() {
		var id, roleID, componentID, createdAt, updatedAt interface{}

		err := rows.Scan(&id, &roleID, &componentID, &createdAt, &updatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		mapping := map[string]interface{}{
			"id":           id,
			"role_id":      roleID,
			"component_id": componentID,
			"created_at":   createdAt,
			"updated_at":   updatedAt,
		}

		mappings = append(mappings, mapping)
	}

	json.NewEncoder(w).Encode(mappings)
}

// UpdateRoleComponent updates an existing role-component mapping
func (handler *RolesComponentsHandler) UpdateRoleComponent(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var assignment struct {
		RoleID      int `json:"role_id"`
		ComponentID int `json:"component_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&assignment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := "UPDATE roles_components SET role_id = ?, component_id = ? WHERE id = ?"
	_, err = handler.DB.Exec(query, assignment.RoleID, assignment.ComponentID, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "role-component mapping updated successfully"})
}

// DeleteRoleComponent deletes a role-component mapping by ID
func (handler *RolesComponentsHandler) DeleteRoleComponent(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	query := "DELETE FROM roles_components WHERE id = ?"
	_, err := handler.DB.Exec(query, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "role-component mapping deleted successfully"})
}

// DeleteAllRolesComponents deletes all role-component mappings
func (handler *RolesComponentsHandler) DeleteAllRolesComponents(w http.ResponseWriter, r *http.Request) {
	_, err := handler.DB.Exec("DELETE FROM roles_components")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "all role-component mappings deleted successfully"})
}

// TruncateRolesComponents truncates the roles_components table
func (handler *RolesComponentsHandler) TruncateRolesComponents(w http.ResponseWriter, r *http.Request) {
	var truncateSQL string
	switch handler.DBType {
	case "sqlite":
		truncateSQL = "DELETE FROM roles_components; DELETE FROM sqlite_sequence WHERE name = 'roles_components';"
	case "mysql":
		truncateSQL = "TRUNCATE TABLE roles_components"
	default:
		http.Error(w, "unsupported database type", http.StatusBadRequest)
		return
	}

	_, err := handler.DB.Exec(truncateSQL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "roles_components table truncated successfully"})
}
