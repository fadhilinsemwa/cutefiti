package handlers

import (
	"database/sql"
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

// CreateAllTables creates all tables in the schema
func (handler *TableHandler) CreateAllTables() error {
	tableCreationFuncs := []func() error{
		handler.CreateManufacturersTable,
		handler.CreateProductsTable,
		handler.CreateUnitsTable,
		handler.CreatePurchasesTable,
		handler.CreateSalesTable,
	}

	for _, createFunc := range tableCreationFuncs {
		if err := createFunc(); err != nil {
			return err
		}
	}

	log.Println("All tables created successfully.")
	return nil
}

// CreateManufacturersTable creates the Manufacturers table.
func (handler *TableHandler) CreateManufacturersTable() error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS Manufacturers (
		ManufacturerID INT AUTO_INCREMENT PRIMARY KEY,
		Name VARCHAR(255) NOT NULL,
		ContactInfo TEXT,
		Address TEXT,
		CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := handler.DB.Exec(createTableSQL)
	if err != nil {
		return err
	}

	log.Println("Manufacturers table created successfully.")
	return nil
}

// CreateProductsTable creates the Products table.
func (handler *TableHandler) CreateProductsTable() error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS Products (
		ProductID INT AUTO_INCREMENT PRIMARY KEY,
		Name VARCHAR(255) NOT NULL,
		ManufacturerID INT NOT NULL,
		Category VARCHAR(255),
		UnitPrice DECIMAL(10, 2) NOT NULL,
		QuantityInStockCartons INT DEFAULT 0,
		QuantityInStockBags INT DEFAULT 0,
		QuantityInStockPacks INT DEFAULT 0,
		QuantityInStockPieces INT DEFAULT 0,
		CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (ManufacturerID) REFERENCES Manufacturers(ManufacturerID) ON DELETE CASCADE
	);
	`

	_, err := handler.DB.Exec(createTableSQL)
	if err != nil {
		return err
	}

	log.Println("Products table created successfully.")
	return nil
}

// CreateUnitsTable creates the Units table.
func (handler *TableHandler) CreateUnitsTable() error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS Units (
		UnitID INT AUTO_INCREMENT PRIMARY KEY,
		ProductID INT NOT NULL,
		CartonToBagRatio INT NOT NULL,
		BagToPackRatio INT,
		PackToPieceRatio INT,
		BagToPieceRatio INT,
		CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (ProductID) REFERENCES Products(ProductID) ON DELETE CASCADE
	);
	`

	_, err := handler.DB.Exec(createTableSQL)
	if err != nil {
		return err
	}

	log.Println("Units table created successfully.")
	return nil
}

// CreatePurchasesTable creates the Purchases table.
func (handler *TableHandler) CreatePurchasesTable() error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS Purchases (
		PurchaseID INT AUTO_INCREMENT PRIMARY KEY,
		ProductID INT NOT NULL,
		ManufacturerID INT NOT NULL,
		QuantityCartons INT DEFAULT 0,
		QuantityBags INT DEFAULT 0,
		QuantityPacks INT DEFAULT 0,
		QuantityPieces INT DEFAULT 0,
		TotalCost DECIMAL(10, 2) NOT NULL,
		PurchaseDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (ProductID) REFERENCES Products(ProductID) ON DELETE CASCADE,
		FOREIGN KEY (ManufacturerID) REFERENCES Manufacturers(ManufacturerID) ON DELETE CASCADE
	);
	`

	_, err := handler.DB.Exec(createTableSQL)
	if err != nil {
		return err
	}

	log.Println("Purchases table created successfully.")
	return nil
}

// CreateSalesTable creates the Sales table.
func (handler *TableHandler) CreateSalesTable() error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS Sales (
		SaleID INT AUTO_INCREMENT PRIMARY KEY,
		ProductID INT NOT NULL,
		CustomerID INT NOT NULL,
		QuantityCartons INT DEFAULT 0,
		QuantityBags INT DEFAULT 0,
		QuantityPacks INT DEFAULT 0,
		QuantityPieces INT DEFAULT 0,
		TotalPrice DECIMAL(10, 2) NOT NULL,
		PaymentReceived DECIMAL(10, 2) DEFAULT NULL,
		SaleDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (ProductID) REFERENCES Products(ProductID) ON DELETE CASCADE,
		FOREIGN KEY (CustomerID) REFERENCES Customers(CustomerID) ON DELETE CASCADE
	);
	`

	_, err := handler.DB.Exec(createTableSQL)
	if err != nil {
		return err
	}

	log.Println("Sales table created successfully.")
	return nil
}
