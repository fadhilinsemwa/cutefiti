package handlers

import (
	"database/sql"
	"encoding/json"
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

	query := `INSERT INTO Manufacturers (Name, ContactInfo, Address) VALUES (?, ?, ?)`
	result, err := h.DB.Exec(query, manufacturer.Name, manufacturer.ContactInfo, manufacturer.Address)
	if err != nil {
		http.Error(w, "Failed to create manufacturer", http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	manufacturer.ManufacturerID = int(id)
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

	var manufacturer Manufacturer
	query := `SELECT ManufacturerID, Name, ContactInfo, Address, CreatedAt FROM Manufacturers WHERE ManufacturerID = ?`
	row := h.DB.QueryRow(query, id)
	if err := row.Scan(&manufacturer.ManufacturerID, &manufacturer.Name, &manufacturer.ContactInfo, &manufacturer.Address, &manufacturer.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Manufacturer not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve manufacturer", http.StatusInternalServerError)
		}
		return
	}

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

	var manufacturer Manufacturer
	if err := json.NewDecoder(r.Body).Decode(&manufacturer); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	query := `UPDATE Manufacturers SET Name = ?, ContactInfo = ?, Address = ? WHERE ManufacturerID = ?`
	_, err := h.DB.Exec(query, manufacturer.Name, manufacturer.ContactInfo, manufacturer.Address, id)
	if err != nil {
		http.Error(w, "Failed to update manufacturer", http.StatusInternalServerError)
		return
	}

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

	query := `DELETE FROM Manufacturers WHERE ManufacturerID = ?`
	_, err := h.DB.Exec(query, id)
	if err != nil {
		http.Error(w, "Failed to delete manufacturer", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Manufacturer deleted successfully"))
}

// ListManufacturers retrieves all manufacturers
func (h *ManufacturerHandler) ListManufacturers(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(`SELECT ManufacturerID, Name, ContactInfo, Address, CreatedAt FROM Manufacturers`)
	if err != nil {
		http.Error(w, "Failed to retrieve manufacturers", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var manufacturers []Manufacturer
	for rows.Next() {
		var manufacturer Manufacturer
		if err := rows.Scan(&manufacturer.ManufacturerID, &manufacturer.Name, &manufacturer.ContactInfo, &manufacturer.Address, &manufacturer.CreatedAt); err != nil {
			http.Error(w, "Failed to parse manufacturers", http.StatusInternalServerError)
			return
		}
		manufacturers = append(manufacturers, manufacturer)
	}

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
