# Product Catalog & Search API (Go + PostgreSQL + Elasticsearch)

REST API layanan katalog produk yang memisahkan penyimpanan data transaksional (**PostgreSQL**) dan mesin pencarian teks cepat (**Elasticsearch**) menggunakan **Go (Gin)**.

## Fitur Utama
- **CRUD Lengkap (`POST`, `PUT`, `DELETE`)**: Setiap perubahan data di PostgreSQL otomatis tersinkronisasi ke Elasticsearch.
- **Full-Text Search & Typo Tolerance**: Pencarian multi-kolom (`name` dan `description`) menggunakan `fuzziness: AUTO`.
- **Filter Rentang Harga**: Mendukung pencarian dengan parameter `min_price` dan `max_price` menggunakan Bool & Range Query.
- **Clean Architecture**: Pemisahan lapisan `handler`, `repository`, dan `models`.

## Cara Menjalankan Proyek
1. Jalankan PostgreSQL dan Elasticsearch dengan Docker:
   ```bash
   docker compose up -d
   ```
2. Salin konfigurasi `.env` dan sesuaikan kredensial database.
3. Jalankan aplikasi Go:
   ```bash
   go mod tidy
   go run main.go
   ```