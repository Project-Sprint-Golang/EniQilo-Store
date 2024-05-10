package model

import "time"

type ProductRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=30"`
	SKU         string  `json:"sku" binding:"required,min=1,max=30"`
	Category    string  `json:"category" binding:"required,oneof=Clothing Accessories Footwear Beverages"`
	ImageURL    string  `json:"imageUrl" binding:"required,url"`
	Notes       string  `json:"notes" binding:"required,min=1,max=200"`
	Price       float64 `json:"price" binding:"required,min=1"`
	Stock       int     `json:"stock" binding:"required,min=0,max=100000"`
	Location    string  `json:"location" binding:"required,min=1,max=200"`
	IsAvailable bool    `json:"isAvailable" binding:"required"`
}

type ProductResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	SKU         string    `json:"sku"`
	Category    string    `json:"category"`
	ImageURL    string    `json:"imageUrl"`
	Notes       string    `json:"notes"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	Location    string    `json:"location"`
	IsAvailable bool      `json:"isAvailable"`
	CreatedAt   time.Time `json:"createdAt"`
}

type GetProductParams struct {
	ID          string `form:"id"`
	Limit       int    `form:"limit,default=5"`
	Offset      int    `form:"offset,default=0"`
	Name        string `form:"name"`
	IsAvailable string `form:"isAvailable"`
	Category    string `form:"category"`
	SKU         string `form:"sku"`
	PriceSort   string `form:"price"`
	InStock     string `form:"inStock"`
	CreatedAt   string `form:"createdAt"`
}
