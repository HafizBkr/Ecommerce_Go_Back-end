package search

import (
    "ecommerce-api/models"  // Assurez-vous que ce chemin correspond à votre structure de projet
)

type SearchOptions struct {
    Query    string
    Page     int
    PageSize int
}

type SearchResult struct {
    Products []models.Product `json:"products"`
    Total    int             `json:"total"`
    Page     int             `json:"page"`
    PageSize int             `json:"page_size"`
}
