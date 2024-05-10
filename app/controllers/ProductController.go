package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	model "github.com/Project-Sprint-Golang/EniQilo-Store/app/models"
	"github.com/Project-Sprint-Golang/EniQilo-Store/config"
	"github.com/Project-Sprint-Golang/EniQilo-Store/helper"
	"github.com/gin-gonic/gin"
)

func AddProduct(c *gin.Context) {
	var product model.ProductRequest
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	isValid := helper.ValidateURL(product.ImageURL)
	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	_, err := config.DB.Query("INSERT INTO products (name, sku, category, imageUrl, notes, price, stock, location, isAvailable)VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)", product.Name, product.SKU, product.Category, product.ImageURL, product.Notes, product.Price, product.Stock, product.Location, product.IsAvailable)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error when Add Product"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Product added successfully"})

}

func GetAllProduct(c *gin.Context) {
	var products []model.ProductResponse
	query := "SELECT id, name, sku, category, imageUrl, notes, price, stock, location, isAvailable, createdAt FROM products WHERE 1=1 AND deletedAt IS NULL"
	args := []interface{}{}

	var params model.GetProductParams

	if err := c.BindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}

	if params.ID != "" {
		query += " AND id = $" + strconv.Itoa(len(args)+1)
		args = append(args, params.ID)
	}
	if params.Name != "" {
		query += " AND lower(name) LIKE $" + strconv.Itoa(len(args)+1)
		args = append(args, "%"+strings.ToLower(params.Name)+"%")
	}
	if params.IsAvailable == "true" || params.IsAvailable == "1" {
		query += " AND isAvailable = $" + strconv.Itoa(len(args)+1)
		args = append(args, params.IsAvailable)
	} else if params.IsAvailable == "false" || params.IsAvailable == "0" {
		query += " AND isAvailable = $" + strconv.Itoa(len(args)+1)
		args = append(args, params.IsAvailable)
	}

	if params.Category != "" {
		switch params.Category {
		case "Clothing", "Accessories", "Footwear", "Beverages":
			query += " AND category = $" + strconv.Itoa(len(args)+1)
			args = append(args, params.Category)
		default:

		}
	}
	if params.SKU != "" {
		query += " AND sku =$" + strconv.Itoa(len(args)+1)
		args = append(args, params.SKU)
	}
	if params.InStock != "" {
		if params.InStock == "true" || params.InStock == "1" {
			query += " AND stock > 0"
		}
		if params.InStock == "false" || params.InStock == "0" {
			query += " AND stock = 0"
		}
	}
	if params.PriceSort == "asc" {
		query += " ORDER BY price ASC"
	} else if params.PriceSort == "desc" {
		query += " ORDER BY price DESC"
	}
	if params.CreatedAt == "asc" {
		query += " ORDER BY createdAt ASC"
	} else if params.CreatedAt == "desc" {
		query += " ORDER BY createdAt DESC"
	}
	query += " LIMIT $" + strconv.Itoa(len(args)+1)
	args = append(args, params.Limit)
	query += " OFFSET $" + strconv.Itoa(len(args)+1)
	args = append(args, params.Offset)
	fmt.Println("limit ", params.Limit)
	fmt.Println(query)
	rows, err := config.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error when Retrieve 1"})
		return
	}
	// defer rows.Close()
	for rows.Next() {
		var p model.ProductResponse
		err := rows.Scan(&p.ID, &p.Name, &p.SKU, &p.Category, &p.ImageURL, &p.Notes, &p.Price, &p.Stock, &p.Location, &p.IsAvailable, &p.CreatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error when Retrieve"})
			return
		}
		products = append(products, p)
	}
	if len(products) == 0 {
		// Return response with empty array
		c.JSON(http.StatusOK, gin.H{
			"message": "Success",
			"data":    []gin.H{}, // Empty array
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Message": "Success",
		"data":    products,
	})
}

func UpdateProduct(c *gin.Context) {
	var product model.ProductRequest
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	productID := c.Param("id")

	isValid := helper.ValidateURL(product.ImageURL)
	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var exists bool
	err := config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM products WHERE id = $1 AND deletedAt IS NULL)", productID).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check product existence"})
		return
	}
	if !exists {
		// Product does not exist
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	query := `
    UPDATE products 
    SET 
        name = $1,
        sku = $2,
        category = $3,
        imageUrl = $4,
        notes = $5,
        price = $6,
        stock = $7,
        location = $8,
        isAvailable = $9
    WHERE 
        id = $10
        AND deletedAt IS NULL
`
	_, err = config.DB.Query(query, product.Name, product.SKU, product.Category, product.ImageURL, product.Notes, product.Price, product.Stock, product.Location, product.IsAvailable, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error when Add Product"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully Update Product"})
}

func DeleteProduct(c *gin.Context) {
	productID := c.Param("id")
	var exists bool
	err := config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM products WHERE id = $1 AND deletedAt IS NULL)", productID).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product Not Found"})
		return
	}
	if !exists {
		// Product does not exist
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	query := `
    UPDATE products 
    SET 
        deletedAt = NOW()
    WHERE 
        id = $1
        AND deletedAt IS NULL
`
	_, err = config.DB.Query(query, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error when Add Product"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully Delete Product"})
}

// Search SKU

func GetSKUProduct(c *gin.Context) {
	var products []model.ProductResponse
	query := "SELECT id, name, sku, category, imageUrl, notes, price, stock, location, isAvailable, createdAt FROM products WHERE 1=1 AND isAvailable = true AND deletedAt IS NULL"
	args := []interface{}{}

	var params model.GetProductParams

	if err := c.BindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}

	if params.Name != "" {
		query += " AND lower(name) LIKE $" + strconv.Itoa(len(args)+1)
		args = append(args, "%"+strings.ToLower(params.Name)+"%")
	}

	if params.Category != "" {
		switch params.Category {
		case "Clothing", "Accessories", "Footwear", "Beverages":
			query += " AND category = $" + strconv.Itoa(len(args)+1)
			args = append(args, params.Category)
		default:

		}
	}
	if params.SKU != "" {
		query += " AND sku =$" + strconv.Itoa(len(args)+1)
		args = append(args, params.SKU)
	}
	if params.InStock != "" {
		if params.InStock == "true" || params.InStock == "1" {
			query += " AND stock > 0"
		}
		if params.InStock == "false" || params.InStock == "0" {
			query += " AND stock = 0"
		}
	}
	if params.PriceSort == "asc" {
		query += " ORDER BY price ASC"
	} else if params.PriceSort == "desc" {
		query += " ORDER BY price DESC"
	}

	query += " LIMIT $" + strconv.Itoa(len(args)+1)
	args = append(args, params.Limit)
	query += " OFFSET $" + strconv.Itoa(len(args)+1)
	args = append(args, params.Offset)
	fmt.Println("limit ", params.Limit)
	fmt.Println(query)
	rows, err := config.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error when Retrieve 1"})
		return
	}
	// defer rows.Close()
	for rows.Next() {
		var p model.ProductResponse
		err := rows.Scan(&p.ID, &p.Name, &p.SKU, &p.Category, &p.ImageURL, &p.Notes, &p.Price, &p.Stock, &p.Location, &p.IsAvailable, &p.CreatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error when Retrieve"})
			return
		}
		products = append(products, p)
	}
	if len(products) == 0 {
		// Return response with empty array
		c.JSON(http.StatusOK, gin.H{
			"message": "Success",
			"data":    []gin.H{}, // Empty array
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Message": "Success",
		"data":    products,
	})
}
