package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// Customer represents the Customers table structure
type Customer struct {
	CustomerID   int       `json:"customerID"`
	Name         string    `json:"name"`
	PhoneNumber  string    `json:"phoneNumber"`
	ContactInfo  string    `json:"contactInfo"`
	CustomerType string    `json:"customerType"`
	Address      string    `json:"address"`
	CreatedAt    time.Time `json:"createdAt"`
}

// CustomerHandler structure
type CustomerHandler struct {
	DB *sql.DB
}

// NewCustomerHandler creates a new CustomerHandler
func NewCustomerHandler(db *sql.DB) *CustomerHandler {
	return &CustomerHandler{DB: db}
}

// CreateCustomer creates a new customer
func (h *CustomerHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var customer Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Printf("Decoded customer: %+v", customer)

	query := `INSERT INTO Customers (Name, PhoneNumber, ContactInfo, CustomerType, Address) VALUES (?, ?, ?, ?, ?)`
	result, err := h.DB.Exec(query, customer.Name, customer.PhoneNumber, customer.ContactInfo, customer.CustomerType, customer.Address)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		http.Error(w, "Failed to create customer", http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	customer.CustomerID = int(id)
	customer.CreatedAt = time.Now()

	log.Printf("Created customer: %+v", customer)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customer)
}

// GetCustomer retrieves a single customer by ID
func (h *CustomerHandler) GetCustomer(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing ID parameter", http.StatusBadRequest)
		return
	}

	log.Printf("Retrieving customer with ID: %s", id)

	var customer Customer
	var createdAt string

	query := `SELECT CustomerID, Name, PhoneNumber, ContactInfo, CustomerType, Address, CreatedAt FROM Customers WHERE CustomerID = ?`
	row := h.DB.QueryRow(query, id)

	if err := row.Scan(&customer.CustomerID, &customer.Name, &customer.PhoneNumber, &customer.ContactInfo, &customer.CustomerType, &customer.Address, &createdAt); err != nil {
		log.Printf("Error scanning row: %v", err)
		if err == sql.ErrNoRows {
			http.Error(w, "Customer not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve customer", http.StatusInternalServerError)
		}
		return
	}

	parsedTime, err := time.Parse("2006-01-02 15:04:05", createdAt)
	if err != nil {
		log.Printf("Error parsing CreatedAt: %v", err)
		http.Error(w, "Failed to parse CreatedAt field", http.StatusInternalServerError)
		return
	}

	customer.CreatedAt = parsedTime

	log.Printf("Retrieved customer: %+v", customer)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customer)
}

// ListCustomers retrieves all customers
func (h *CustomerHandler) ListCustomers(w http.ResponseWriter, r *http.Request) {
	log.Println("Listing all customers")

	query := `SELECT CustomerID, Name, PhoneNumber, ContactInfo, CustomerType, Address, CreatedAt FROM Customers`
	rows, err := h.DB.Query(query)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		http.Error(w, "Failed to retrieve customers", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	customers := []Customer{}
	for rows.Next() {
		var customer Customer
		var createdAt string
		if err := rows.Scan(&customer.CustomerID, &customer.Name, &customer.PhoneNumber, &customer.ContactInfo, &customer.CustomerType, &customer.Address, &createdAt); err != nil {
			log.Printf("Error scanning row: %v", err)
			http.Error(w, "Failed to parse customer data", http.StatusInternalServerError)
			return
		}

		parsedTime, err := time.Parse("2006-01-02 15:04:05", createdAt)
		if err != nil {
			log.Printf("Error parsing CreatedAt: %v", err)
			http.Error(w, "Failed to parse CreatedAt field", http.StatusInternalServerError)
			return
		}

		customer.CreatedAt = parsedTime
		customers = append(customers, customer)
	}

	log.Printf("Retrieved customers: %+v", customers)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
}

// UpdateCustomer updates a customer's details
func (h *CustomerHandler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing ID parameter", http.StatusBadRequest)
		return
	}

	log.Printf("Updating customer with ID: %s", id)

	var customer Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Printf("Payload for update: %+v", customer)

	query := `UPDATE Customers SET Name = ?, PhoneNumber = ?, ContactInfo = ?, CustomerType = ?, Address = ? WHERE CustomerID = ?`
	_, err := h.DB.Exec(query, customer.Name, customer.PhoneNumber, customer.ContactInfo, customer.CustomerType, customer.Address, id)
	if err != nil {
		log.Printf("Error executing update query: %v", err)
		http.Error(w, "Failed to update customer", http.StatusInternalServerError)
		return
	}

	log.Printf("Customer updated successfully")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Customer updated successfully"))
}

// DeleteCustomer deletes a customer by ID
func (h *CustomerHandler) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing ID parameter", http.StatusBadRequest)
		return
	}

	log.Printf("Deleting customer with ID: %s", id)

	query := `DELETE FROM Customers WHERE CustomerID = ?`
	_, err := h.DB.Exec(query, id)
	if err != nil {
		log.Printf("Error executing delete query: %v", err)
		http.Error(w, "Failed to delete customer", http.StatusInternalServerError)
		return
	}

	log.Printf("Customer deleted successfully")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Customer deleted successfully"))
}

// RegisterCustomerRoutes registers HTTP routes for customers
func (h *CustomerHandler) RegisterCustomerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/customers", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.CreateCustomer(w, r)
		case http.MethodGet:
			if r.URL.Query().Has("id") {
				h.GetCustomer(w, r)
			} else {
				h.ListCustomers(w, r)
			}
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/customers/update", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			h.UpdateCustomer(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/customers/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			h.DeleteCustomer(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
