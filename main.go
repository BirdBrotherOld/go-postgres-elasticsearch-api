package main

import (
	"belajar-go-pg-es/handler"
	"belajar-go-pg-es/repository"
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Muat variabel dari file .env
	if err := godotenv.Load(); err != nil {
		log.Println("File .env tidak ditemukan, menggunakan environment sistem")
	}

	// 2. Koneksi PostgreSQL dari .env
	dsn := os.Getenv("DB_DSN")
	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal("Gagal koneksi ke PostgreSQL: ", err)
	}
	defer db.Close()

	// 3. Koneksi Elasticsearch dari .env
	esCfg := elasticsearch.Config{
		Addresses: []string{os.Getenv("ELASTIC_URL")},
	}
	es, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		log.Fatal("Gagal membuat client Elasticsearch: ", err)
	}

	// 4. Rakit Repository & Handler
	pgRepo := repository.NewPostgresRepo(db)
	esRepo := repository.NewElasticRepo(es)
	productHandler := handler.NewProductHandler(pgRepo, esRepo)

	// 5. Daftarkan Rute & Jalankan Server
	r := gin.Default()
	r.POST("/products", productHandler.Create)
	r.GET("/search", productHandler.Search)
	r.PUT("/products/:id", productHandler.Update)
	r.DELETE("/products/:id", productHandler.Delete)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = ":8080"
	}
	r.Run(port)
}