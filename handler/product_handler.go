package handler

import (
	"belajar-go-pg-es/models"
	"belajar-go-pg-es/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	pgRepo *repository.PostgresRepo
	esRepo *repository.ElasticRepo
}

func NewProductHandler(pg *repository.PostgresRepo, es *repository.ElasticRepo) *ProductHandler {
	return &ProductHandler{pgRepo: pg, esRepo: es}
}

func (h *ProductHandler) Create(c *gin.Context) {
	var p models.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format JSON salah"})
		return
	}

	if err := h.pgRepo.Create(c.Request.Context(), &p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.esRepo.Upsert(p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Produk dibuat", "data": p})
}

func (h *ProductHandler) Search(c *gin.Context) {
	keyword := c.Query("q")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter ?q= wajib diisi"})
		return
	}

	// Ambil parameter harga (jika kosong otomatis bernilai 0)
	minPrice, _ := strconv.Atoi(c.Query("min_price"))
	maxPrice, _ := strconv.Atoi(c.Query("max_price"))

	products, err := h.esRepo.Search(c.Request.Context(), keyword, minPrice, maxPrice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"keyword":   keyword,
		"min_price": minPrice,
		"max_price": maxPrice,
		"total":     len(products),
		"results":   products,
	})
}

func (h *ProductHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID harus angka"})
		return
	}

	var p models.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format JSON salah"})
		return
	}
	p.ID = id

	rows, err := h.pgRepo.Update(c.Request.Context(), p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
		return
	}

	if err := h.esRepo.Upsert(p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Produk diperbarui", "data": p})
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID harus angka"})
		return
	}

	rows, err := h.pgRepo.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
		return
	}

	if err := h.esRepo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Produk dihapus permanen"})
}