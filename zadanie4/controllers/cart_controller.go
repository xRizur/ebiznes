package controllers

import (
	"net/http"

	"shop/config"
	"shop/models"

	"github.com/labstack/echo/v4"
)

func CreateCart(c echo.Context) error {
	cart := new(models.Cart)
	if err := c.Bind(cart); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if err := config.DB.Create(cart).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, cart)
}

func GetCartByID(c echo.Context) error {
	id := c.Param("id")
	var cart models.Cart

	if err := config.DB.Preload("Products").First(&cart, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Cart not found"})
	}

	return c.JSON(http.StatusOK, cart)
}

func AddProductToCart(c echo.Context) error {
	cartID := c.Param("cart_id")
	productID := c.Param("product_id")

	var cart models.Cart
	if err := config.DB.First(&cart, cartID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Cart not found"})
	}

	var product models.Product
	if err := config.DB.First(&product, productID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Product not found"})
	}

	if err := config.DB.Model(&cart).Association("Products").Append(&product); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, cart)
}

func RemoveProductFromCart(c echo.Context) error {
	cartID := c.Param("cart_id")
	productID := c.Param("product_id")

	var cart models.Cart
	if err := config.DB.Preload("Products").First(&cart, cartID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Cart not found"})
	}

	var product models.Product
	if err := config.DB.First(&product, productID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Product not found"})
	}

	if err := config.DB.Model(&cart).Association("Products").Delete(&product); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, cart)
}
