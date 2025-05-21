package main

import (
	"net/http"
	"strings"

	"github.com/glebarez/sqlite"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

const (
	placeholderImageURL = "https://microless.com/cdn/products/f026b0f0fb6302d095eda73e25215408-hi.jpg"
)

type Product struct {
	ID          uint    `json:"id" gorm:"primaryKey"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"imageUrl"`
}

type CartItem struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	ProductID uint    `json:"productId"`
	Product   Product `json:"product" gorm:"foreignKey:ProductID"`
	Quantity  int     `json:"quantity"`
}

type Payment struct {
	ID         uint    `json:"id" gorm:"primaryKey"`
	Amount     float64 `json:"amount"`
	CardNumber string  `json:"cardNumber"`
	CardHolder string  `json:"cardHolder"`
	ExpiryDate string  `json:"expiryDate"`
	CVV        string  `json:"cvv"`
	Status     string  `json:"status"`
}

func main() {
	db := initializeDatabase()

	e := setupServer()

	registerRoutes(e, db)

	e.Logger.Fatal(e.Start(":8080"))
}

func initializeDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("shop.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Product{}, &CartItem{}, &Payment{})
	seedDatabaseIfEmpty(db)

	return db
}

func seedDatabaseIfEmpty(db *gorm.DB) {
	var count int64
	db.Model(&Product{}).Count(&count)

	if count == 0 {
		products := []Product{
			{Name: "Laptop", Description: "Wydajny laptop dla programistów", Price: 3999.99, ImageURL: placeholderImageURL},
			{Name: "Smartfon", Description: "Smartfon z najnowszym systemem", Price: 1999.99, ImageURL: placeholderImageURL},
			{Name: "Słuchawki", Description: "Słuchawki z redukcją szumów", Price: 399.99, ImageURL: placeholderImageURL},
			{Name: "Mysz komputerowa", Description: "Bezprzewodowa mysz ergonomiczna", Price: 149.99, ImageURL: placeholderImageURL},
		}
		db.Create(&products)
	}
}

// Set up Echo server with middleware
func setupServer() *echo.Echo {
	e := echo.New()

	// Add middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	return e
}

// Register all API routes
func registerRoutes(e *echo.Echo, db *gorm.DB) {
	// Products endpoints
	e.GET("/products", func(c echo.Context) error {
		return handleGetProducts(c, db)
	})

	// Cart endpoints
	e.GET("/cart", func(c echo.Context) error {
		return handleGetCart(c, db)
	})

	e.POST("/cart", func(c echo.Context) error {
		return handleAddToCart(c, db)
	})

	// Payment endpoints
	e.GET("/payments", func(c echo.Context) error {
		return handleGetPayments(c, db)
	})

	e.POST("/payments", func(c echo.Context) error {
		return handleCreatePayment(c, db)
	})
}

// Handle get products request
func handleGetProducts(c echo.Context, db *gorm.DB) error {
	var products []Product
	result := db.Find(&products)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Nie udało się pobrać produktów"})
	}
	return c.JSON(http.StatusOK, products)
}

// Handle get cart request
func handleGetCart(c echo.Context, db *gorm.DB) error {
	var cartItems []CartItem
	result := db.Preload("Product").Find(&cartItems)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Nie udało się pobrać koszyka"})
	}
	return c.JSON(http.StatusOK, cartItems)
}

// Handle add to cart request
func handleAddToCart(c echo.Context, db *gorm.DB) error {
	cartItem := new(CartItem)
	if err := c.Bind(cartItem); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Nieprawidłowe dane"})
	}

	if !productExists(db, cartItem.ProductID) {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Produkt nie istnieje"})
	}

	var product Product
	db.First(&product, cartItem.ProductID)
	cartItem.Product = product

	result := db.Create(&cartItem)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Nie udało się dodać do koszyka"})
	}

	return c.JSON(http.StatusCreated, cartItem)
}

// Handle get payments request
func handleGetPayments(c echo.Context, db *gorm.DB) error {
	var payments []Payment
	result := db.Find(&payments)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Nie udało się pobrać historii płatności"})
	}
	return c.JSON(http.StatusOK, payments)
}

// Handle create payment request
func handleCreatePayment(c echo.Context, db *gorm.DB) error {
	payment := new(Payment)
	if err := c.Bind(payment); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Nieprawidłowe dane płatności"})
	}

	if err := validatePayment(payment); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	payment.Status = "completed"
	result := db.Create(&payment)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Błąd przetwarzania płatności"})
	}

	db.Exec("DELETE FROM cart_items")

	return c.JSON(http.StatusCreated, payment)
}

// Validate payment details
func validatePayment(payment *Payment) error {
	if payment.CardNumber == "" {
		return errorStrings("Numer karty jest wymagany")
	}

	cardNumber := strings.ReplaceAll(payment.CardNumber, " ", "")
	if len(cardNumber) != 16 || !isNumeric(cardNumber) {
		return errorStrings("Nieprawidłowy format numeru karty")
	}

	if payment.CardHolder == "" {
		return errorStrings("Imię i nazwisko jest wymagane")
	}

	if payment.ExpiryDate == "" {
		return errorStrings("Data ważności jest wymagana")
	} else if !isValidExpiryDate(payment.ExpiryDate) {
		return errorStrings("Nieprawidłowy format daty")
	}

	if payment.CVV == "" {
		return errorStrings("Kod CVV jest wymagany")
	} else if len(payment.CVV) < 3 || len(payment.CVV) > 4 || !isNumeric(payment.CVV) {
		return errorStrings("Nieprawidłowy format CVV")
	}

	if payment.Amount <= 0 {
		return errorStrings("Kwota płatności jest nieprawidłowa")
	}

	return nil
}

// Helper type for string errors
type errorStrings string

// Error returns the error string
func (e errorStrings) Error() string {
	return string(e)
}

// Check if a product exists in the database
func productExists(db *gorm.DB, productID uint) bool {
	var product Product
	result := db.First(&product, productID)
	return result.Error == nil
}

// Check if a string consists only of numeric characters
func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// Validate the format of an expiry date (MM/YY)
func isValidExpiryDate(date string) bool {
	if len(date) != 5 || date[2] != '/' {
		return false
	}

	month := date[:2]
	year := date[3:]

	if !isNumeric(month) || !isNumeric(year) {
		return false
	}

	monthNum := 0
	for _, c := range month {
		monthNum = monthNum*10 + int(c-'0')
	}

	return monthNum >= 1 && monthNum <= 12
}
