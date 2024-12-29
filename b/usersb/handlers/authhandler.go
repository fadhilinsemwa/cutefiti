package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4" // Updated import for jwt-go
	"golang.org/x/crypto/bcrypt"
)

// AuthHandler structure to manage authentication logic
type AuthHandler struct {
	DB     *sql.DB
	DBType string
}

// NewAuthHandler initializes a new AuthHandler
func NewAuthHandler(db *sql.DB, dbType string) *AuthHandler {
	return &AuthHandler{DB: db, DBType: dbType}
}

// JSONResponse is a helper function to write JSON responses
func JSONResponse(w http.ResponseWriter, statusCode int, message interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(message)
}

// RegisterUser registers a new user with the provided information
func (handler *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user struct {
		Username  string `json:"username"`
		Email     string `json:"email"`
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
		City      string `json:"city"`
		Password  string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		JSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid input"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
		return
	}

	query := `INSERT INTO users (username, email, firstname, lastname, city, password) VALUES (?, ?, ?, ?, ?, ?)`
	_, err = handler.DB.Exec(query, user.Username, user.Email, user.Firstname, user.Lastname, user.City, string(hashedPassword))
	if err != nil {
		JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to register user"})
		return
	}

	log.Println("User registered successfully.")
	JSONResponse(w, http.StatusOK, map[string]string{"message": "User registered successfully"})
}

// LoginUser handles user login and issues a JWT token along with the user_id
func (handler *AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		JSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid input"})
		return
	}

	log.Println("Attempting login for username:", credentials.Username)

	var storedPassword string
	var userID int
	query := `SELECT user_id, password FROM users WHERE username = ?`
	err := handler.DB.QueryRow(query, credentials.Username).Scan(&userID, &storedPassword)
	if err != nil {
		log.Printf("Error fetching user from DB: %v", err)
		JSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Invalid username or password"})
		return
	}

	log.Println("User found with ID:", userID)

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(credentials.Password))
	if err != nil {
		log.Printf("Password comparison failed for username: %v, error: %v", credentials.Username, err)
		JSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Invalid username or password"})
		return
	}

	log.Println("Password match successful for username:", credentials.Username)

	// Create JWT token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": credentials.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign token with a secret key (you can use a constant or environment variable)
	tokenString, err := token.SignedString([]byte("your_secret_key"))
	if err != nil {
		log.Printf("Error generating JWT token for username: %v, error: %v", credentials.Username, err)
		JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
		return
	}

	log.Println("User logged in successfully, issuing token.")
	JSONResponse(w, http.StatusOK, map[string]interface{}{
		"token":   tokenString,
		"user_id": userID,
	})
}

// LogoutUser invalidates the user's session by removing the JWT token and user_id
func (handler *AuthHandler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	log.Println("User logged out successfully.")
	http.SetCookie(w, &http.Cookie{
		Name:   "token",
		MaxAge: -1, // Expire the token immediately
	})
	http.SetCookie(w, &http.Cookie{
		Name:   "user_id",
		MaxAge: -1, // Expire the user_id immediately
	})
	JSONResponse(w, http.StatusOK, map[string]string{"message": "Logged out successfully"})
}

// RequestPasswordReset generates a reset token for the user and returns it
func (handler *AuthHandler) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		JSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid input"})
		return
	}

	resetToken := generateToken()
	query := `UPDATE users SET reset_token = ? WHERE email = ?`
	_, err := handler.DB.Exec(query, resetToken, request.Email)
	if err != nil {
		JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error generating reset token"})
		return
	}

	log.Println("Password reset requested successfully.")
	JSONResponse(w, http.StatusOK, map[string]string{"reset_token": resetToken})
}

// ResetPassword updates the user's password if the reset token is valid
func (handler *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var reset struct {
		ResetToken  string `json:"reset_token"`
		NewPassword string `json:"new_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reset); err != nil {
		JSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid input"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reset.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
		return
	}

	query := `UPDATE users SET password = ?, reset_token = NULL WHERE reset_token = ?`
	_, err = handler.DB.Exec(query, string(hashedPassword), reset.ResetToken)
	if err != nil {
		JSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Invalid or expired reset token"})
		return
	}

	log.Println("Password reset successfully.")
	JSONResponse(w, http.StatusOK, map[string]string{"message": "Password reset successfully"})
}

// VerifyEmail verifies the user's email address with a verification token
func (handler *AuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		JSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid input"})
		return
	}

	query := `UPDATE users SET email_verified = 1, verification_token = NULL WHERE verification_token = ?`
	_, err := handler.DB.Exec(query, request.Token)
	if err != nil {
		JSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Invalid or expired verification token"})
		return
	}

	log.Println("Email verified successfully.")
	JSONResponse(w, http.StatusOK, map[string]string{"message": "Email verified successfully"})
}

// Helper function to generate a simple token (could be more secure)
func generateToken() string {
	return fmt.Sprintf("%d", time.Now().Unix()) // Simple token generation
}
