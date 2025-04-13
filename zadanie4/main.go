package main

import (
	"log"
	"net/http"
	"shop/config"
	"shop/controllers"
	"shop/models"

	"github.com/labstack/echo/v4"
)

func main() {
	config.ConnectDB()

	err := config.DB.AutoMigrate(
		&models.Product{},
		&models.Category{},
		&models.Cart{},
	)
	if err != nil {
		log.Fatalf("Błąd migracji: %v", err)
	}
	categories := []models.Category{
		{Name: "Electronics"},
		{Name: "Books"},
		{Name: "Clothing"},
		{Name: "Toys"},
		{Name: "Groceries"},
	}

	for _, category := range categories {
		if err := config.DB.FirstOrCreate(&category, models.Category{Name: category.Name}).Error; err != nil {
			log.Printf("Error seeding category %v: %v", category.Name, err)
		}
	}
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Witaj w Go Echo Shop!")
	})

	initRoutes(e)

	e.Logger.Fatal(e.Start(":8080"))
}

func initRoutes(e *echo.Echo) {
	p := e.Group("/products")
	p.POST("", controllers.CreateProduct)
	p.GET("", controllers.GetProducts)
	p.GET("/:id", controllers.GetProductByID)
	p.PUT("/:id", controllers.UpdateProduct)
	p.DELETE("/:id", controllers.DeleteProduct)

	p.GET("/scopes", controllers.GetProductsWithScopes)

	cart := e.Group("/carts")
	cart.POST("", controllers.CreateCart)
	cart.GET("/:id", controllers.GetCartByID)
	cart.POST("/:cart_id/add-product/:product_id", controllers.AddProductToCart)
	cart.DELETE("/:cart_id/remove-product/:product_id", controllers.RemoveProductFromCart)
}
