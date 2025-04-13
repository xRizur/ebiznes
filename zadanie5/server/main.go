package main

import (
	"net/http"

	"github.com/glebarez/sqlite" // Zmiana z gorm.io/driver/sqlite na implementację w czystym Go
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
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
	db, err := gorm.Open(sqlite.Open("shop.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Product{}, &CartItem{}, &Payment{})

	var count int64
	db.Model(&Product{}).Count(&count)
	if count == 0 {
		products := []Product{
			{Name: "Laptop", Description: "Wydajny laptop dla programistów", Price: 3999.99, ImageURL: "https://via.placeholder.com/150"},
			{Name: "Smartfon", Description: "Smartfon z najnowszym systemem", Price: 1999.99, ImageURL: "https://via.placeholder.com/150"},
			{Name: "Słuchawki", Description: "Słuchawki z redukcją szumów", Price: 399.99, ImageURL: "https://via.placeholder.com/150"},
			{Name: "Mysz komputerowa", Description: "Bezprzewodowa mysz ergonomiczna", Price: 149.99, ImageURL: "https://via.placeholder.com/150"},
		}
		db.Create(&products)
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/products", func(c echo.Context) error {
		var products []Product
		result := db.Find(&products)
		if result.Error != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Nie udało się pobrać produktów"})
		}
		return c.JSON(http.StatusOK, products)
	})

	e.GET("/cart", func(c echo.Context) error {
		var cartItems []CartItem
		result := db.Preload("Product").Find(&cartItems)
		if result.Error != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Nie udało się pobrać koszyka"})
		}
		return c.JSON(http.StatusOK, cartItems)
	})

	e.POST("/cart", func(c echo.Context) error {
		cartItem := new(CartItem)
		if err := c.Bind(cartItem); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Nieprawidłowe dane"})
		}

		result := db.Create(&cartItem)
		if result.Error != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Nie udało się dodać do koszyka"})
		}

		return c.JSON(http.StatusCreated, cartItem)
	})

	e.GET("/payments", func(c echo.Context) error {
		var payments []Payment
		result := db.Find(&payments)
		if result.Error != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Nie udało się pobrać historii płatności"})
		}
		return c.JSON(http.StatusOK, payments)
	})

	e.POST("/payments", func(c echo.Context) error {
		payment := new(Payment)
		if err := c.Bind(payment); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Nieprawidłowe dane płatności"})
		}

		payment.Status = "completed"
		result := db.Create(&payment)
		if result.Error != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Błąd przetwarzania płatności"})
		}

		db.Exec("DELETE FROM cart_items")

		return c.JSON(http.StatusCreated, payment)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
