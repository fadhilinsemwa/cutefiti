package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// Manufacturer represents the Manufacturers table structure
type Manufacturer struct {
	ManufacturerID int       `json:"manufacturerID"`
	Name           string    `json:"name"`
	ContactInfo    string    `json:"contactInfo"`
	Address        string    `json:"address"`
	CreatedAt      time.Time `json:"createdAt"`
}

// ManufacturerHandler structure
type ManufacturerHandler struct {
	DB *sql.DB
}

// NewManufacturerHandler creates a new ManufacturerHandler
func NewManufacturerHandler(db *sql.DB) *ManufacturerHandler {
	return &ManufacturerHandler{DB: db}
}

// CreateManufacturer creates a new manufacturer
func (h *ManufacturerHandler) CreateManufacturer(w http.ResponseWriter, r *http.Request) {
	var manufacturer Manufacturer
	if err := json.NewDecoder(r.Body).Decode(&manufacturer); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Printf("Decoded manufacturer: %+v", manufacturer)

	query := `INSERT INTO Manufacturers (Name, ContactInfo, Address) VALUES (?, ?, ?)`
	result, err := h.DB.Exec(query, manufacturer.Name, manufacturer.ContactInfo, manufacturer.Address)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		http.Error(w, "Failed to create manufacturer", http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	manufacturer.ManufacturerID = int(id)
	manufacturer.CreatedAt = time.Now() // Set the current time for CreatedAt

	log.Printf("Created manufacturer: %+v", manufacturer)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(manufacturer)
}

// GetManufacturer retrieves a single manufacturer by ID
func (h *ManufacturerHandler) GetManufacturer(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing ID parameter", http.StatusBadRequest)
		return
	}

	log.Printf("Retrieving manufacturer with ID: %s", id)

	var manufacturer Manufacturer
	var createdAt string // Temporary variable to hold the raw string value of CreatedAt

	query := `SELECT ManufacturerID, Name, ContactInfo, Address, CreatedAt FROM Manufacturers WHERE ManufacturerID = ?`
	row := h.DB.QueryRow(query, id)

	// Use Scan to fill the manufacturer struct and the raw string for CreatedAt
	if err := row.Scan(&manufacturer.ManufacturerID, &manufacturer.Name, &manufacturer.ContactInfo, &manufacturer.Address, &createdAt); err != nil {
		log.Printf("Error scanning row: %v", err)
		if err == sql.ErrNoRows {
			http.Error(w, "Manufacturer not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve manufacturer", http.StatusInternalServerError)
		}
		return
	}

	// Try to parse CreatedAt if it's a string
	parsedTime, err := time.Parse("2006-01-02 15:04:05", createdAt)
	if err != nil {
		log.Printf("Error parsing CreatedAt: %v", err)
		http.Error(w, "Failed to parse CreatedAt field", http.StatusInternalServerError)
		return
	}

	// Set the parsed time into the manufacturer struct
	manufacturer.CreatedAt = parsedTime

	log.Printf("Retrieved manufacturer: %+v", manufacturer)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(manufacturer)
}

// UpdateManufacturer updates a manufacturer's details
func (h *ManufacturerHandler) UpdateManufacturer(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing ID parameter", http.StatusBadRequest)
		return
	}

	log.Printf("Updating manufacturer with ID: %s", id)

	var manufacturer Manufacturer
	if err := json.NewDecoder(r.Body).Decode(&manufacturer); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Printf("Payload for update: %+v", manufacturer)

	query := `UPDATE Manufacturers SET Name = ?, ContactInfo = ?, Address = ? WHERE ManufacturerID = ?`
	_, err := h.DB.Exec(query, manufacturer.Name, manufacturer.ContactInfo, manufacturer.Address, id)
	if err != nil {
		log.Printf("Error executing update query: %v", err)
		http.Error(w, "Failed to update manufacturer", http.StatusInternalServerError)
		return
	}

	log.Printf("Manufacturer updated successfully")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Manufacturer updated successfully"))
}

// DeleteManufacturer deletes a manufacturer by ID
func (h *ManufacturerHandler) DeleteManufacturer(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing ID parameter", http.StatusBadRequest)
		return
	}

	log.Printf("Deleting manufacturer with ID: %s", id)

	query := `DELETE FROM Manufacturers WHERE ManufacturerID = ?`
	_, err := h.DB.Exec(query, id)
	if err != nil {
		log.Printf("Error executing delete query: %v", err)
		http.Error(w, "Failed to delete manufacturer", http.StatusInternalServerError)
		return
	}

	log.Printf("Manufacturer deleted successfully")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Manufacturer deleted successfully"))
}

// ListManufacturers retrieves all manufacturers
func (h *ManufacturerHandler) ListManufacturers(w http.ResponseWriter, r *http.Request) {
	log.Printf("Listing all manufacturers")

	rows, err := h.DB.Query(`SELECT ManufacturerID, Name, ContactInfo, Address, CreatedAt FROM Manufacturers`)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		http.Error(w, "Failed to retrieve manufacturers", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var manufacturers []Manufacturer
	for rows.Next() {
		var manufacturer Manufacturer
		var rawCreatedAt interface{} // Capture the raw CreatedAt value

		// Scan into rawCreatedAt for debugging
		if err := rows.Scan(&manufacturer.ManufacturerID, &manufacturer.Name, &manufacturer.ContactInfo, &manufacturer.Address, &rawCreatedAt); err != nil {
			log.Printf("Error scanning row: %v", err) // Log detailed error
			http.Error(w, "Failed to parse manufacturers", http.StatusInternalServerError)
			return
		}

		// Debug rawCreatedAt type
		log.Printf("Raw CreatedAt value: %v, Type: %T", rawCreatedAt, rawCreatedAt)

		// Handle CreatedAt conversion
		switch v := rawCreatedAt.(type) {
		case time.Time:
			manufacturer.CreatedAt = v
		case []uint8:
			parsedTime, err := time.Parse("2006-01-02 15:04:05", string(v))
			if err != nil {
				log.Printf("Error parsing CreatedAt: %v", err)
				http.Error(w, "Failed to parse CreatedAt field", http.StatusInternalServerError)
				return
			}
			manufacturer.CreatedAt = parsedTime
		default:
			log.Printf("Unexpected CreatedAt type: %T", v)
			http.Error(w, "Unsupported CreatedAt field type", http.StatusInternalServerError)
			return
		}

		log.Printf("Scanned manufacturer: %+v", manufacturer)
		manufacturers = append(manufacturers, manufacturer)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		http.Error(w, "Error iterating over rows", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully retrieved manufacturers: %d rows", len(manufacturers))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(manufacturers)
}

// RegisterManufacturerRoutes registers HTTP routes for manufacturers
func (h *ManufacturerHandler) RegisterManufacturerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/manufacturers", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.CreateManufacturer(w, r)
		case http.MethodGet:
			if r.URL.Query().Has("id") {
				h.GetManufacturer(w, r)
			} else {
				h.ListManufacturers(w, r)
			}
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/manufacturers/update", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			h.UpdateManufacturer(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/manufacturers/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			h.DeleteManufacturer(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
