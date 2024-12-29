package main

import (
	"log"
	"net/http"
	"os"

	"usersb/handlers"
	"usersb/middleware"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Initialize the database handler
	dbHandler, err := handlers.NewDBHandler()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbHandler.CloseDB()

	// Initialize the mux router
	r := mux.NewRouter()

	// Initialize the authentication handler without the secret field
	authHandler := handlers.NewAuthHandler(dbHandler.DB, dbHandler.DBType)

	// Authentication Endpoints
	r.HandleFunc("/register", authHandler.RegisterUser).Methods("POST")
	r.HandleFunc("/login", authHandler.LoginUser).Methods("POST")
	r.HandleFunc("/logout", authHandler.LogoutUser).Methods("POST")
	r.HandleFunc("/request-password-reset", authHandler.RequestPasswordReset).Methods("POST")
	r.HandleFunc("/reset-password", authHandler.ResetPassword).Methods("POST")
	r.HandleFunc("/verify-email", authHandler.VerifyEmail).Methods("POST")

	// Initialize the generic CRUD handler
	crudHandler := handlers.NewCRUDHandler(dbHandler.DB, dbHandler.DBType)

	// Register specific user CRUD endpoints
	r.HandleFunc("/user", crudHandler.HandleCRUD).Methods("POST", "GET", "PUT", "DELETE")  // Consolidated CRUD operations for a single user
	r.HandleFunc("/users", crudHandler.GetAllUsersHandler).Methods("GET")                  // Get all users
	r.HandleFunc("/delete-all-users", crudHandler.DeleteAllUsersHandler).Methods("DELETE") // Delete all users
	r.HandleFunc("/truncate-users", crudHandler.TruncateUsersHandler).Methods("POST")      // Truncate users table

	// Initialize the role handler
	roleHandler := handlers.NewRoleHandler(dbHandler.DB, dbHandler.DBType)

	// Role handlers
	r.HandleFunc("/api/roles", roleHandler.HandleCreateRole).Methods("POST")
	r.HandleFunc("/api/roles", roleHandler.HandleGetAllRoles).Methods("GET")
	r.HandleFunc("/api/assign-role", roleHandler.HandleAssignUserRole).Methods("POST")
	r.HandleFunc("/api/user-roles", roleHandler.HandleGetUserRoles).Methods("GET")

	// Initialize the roles components handler
	rolesComponentsHandler := handlers.NewRolesComponentsHandler(dbHandler.DB, dbHandler.DBType)

	// Components endpoints
	r.HandleFunc("/components", rolesComponentsHandler.CreateComponent).Methods("POST")
	r.HandleFunc("/components", rolesComponentsHandler.GetAllComponents).Methods("GET")
	r.HandleFunc("/components/{id}", rolesComponentsHandler.UpdateComponent).Methods("PUT")
	r.HandleFunc("/components/{id}", rolesComponentsHandler.DeleteComponent).Methods("DELETE")
	r.HandleFunc("/components/truncate", rolesComponentsHandler.TruncateComponents).Methods("POST")
	r.HandleFunc("/components/delete_all", rolesComponentsHandler.DeleteAllComponents).Methods("POST")

	// RolesComponents endpoints
	r.HandleFunc("/roles_components", rolesComponentsHandler.CreateRoleComponent).Methods("POST")
	r.HandleFunc("/roles_components", rolesComponentsHandler.GetAllRolesComponents).Methods("GET")
	r.HandleFunc("/roles_components/{id}", rolesComponentsHandler.UpdateRoleComponent).Methods("PUT")
	r.HandleFunc("/roles_components/{id}", rolesComponentsHandler.DeleteRoleComponent).Methods("DELETE")
	r.HandleFunc("/roles_components/truncate", rolesComponentsHandler.TruncateRolesComponents).Methods("POST")
	r.HandleFunc("/roles_components/delete_all", rolesComponentsHandler.DeleteAllRolesComponents).Methods("POST")

	// Wrap the entire router with the CORS middleware
	corsRouter := middleware.CORSMiddleware(r)

	// Start the server
	port := os.Getenv("PORT1")
	if port == "" {
		port = "8081" // Default port if not specified
	}
	log.Printf("Server started on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, corsRouter))
}
