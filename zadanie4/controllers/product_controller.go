package controllers

import (
	"net/http"
	"strconv"

	"shop/config"
	"shop/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func CreateProduct(c echo.Context) error {
	product := new(models.Product)
	if err := c.Bind(product); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if err := config.DB.Create(product).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, product)
}

func GetProducts(c echo.Context) error {
	var products []models.Product
	if err := config.DB.Preload("Category").Find(&products).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, products)
}

func GetProductByID(c echo.Context) error {
	id := c.Param("id")
	var product models.Product

	if err := config.DB.Preload("Category").First(&product, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Product not found"})
	}

	return c.JSON(http.StatusOK, product)
}

func UpdateProduct(c echo.Context) error {
	id := c.Param("id")

	var product models.Product
	if err := config.DB.First(&product, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Product not found"})
	}

	updateData := new(models.Product)
	if err := c.Bind(updateData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	product.Name = updateData.Name
	product.Description = updateData.Description
	product.Price = updateData.Price
	product.CategoryID = updateData.CategoryID

	if err := config.DB.Save(&product).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, product)
}

func DeleteProduct(c echo.Context) error {
	id := c.Param("id")
	var product models.Product

	if err := config.DB.First(&product, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Product not found"})
	}

	if err := config.DB.Delete(&product).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Product deleted"})
}

func ScopeMinPrice(price float64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("price >= ?", price)
	}
}

func GetProductsWithScopes(c echo.Context) error {
	priceStr := c.QueryParam("min_price")
	if priceStr == "" {
		priceStr = "0"
	}
	minPrice, _ := strconv.ParseFloat(priceStr, 64)

	var products []models.Product
	if err := config.DB.Scopes(ScopeMinPrice(minPrice)).Find(&products).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, products)
}
