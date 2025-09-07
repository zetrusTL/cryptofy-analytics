package main

import (
    "fmt"
    "log"
    "cryptofyanalytics/internal/server"
    "cryptofyanalytics/internal/storage/postgres"
)

func main() {
    fmt.Println("🚀 Starting Crypto Analytics Application...")

    // строка подключения формат: postgres://username:password@host:port/database_name
    connStr := "postgres://crypto_user:crypto_password@localhost:5432/crypto_db?sslmode=disable"

    // Создаем репозиторий для работы с базой данных
    repo, err := postgres.NewRepository(connStr)
    if err != nil {
        log.Fatalf("❌ Failed to create repository: %v", err)
    }
    defer repo.Close() 


   srv := server.NewServer(repo)
    log.Println("🚀 Starting HTTP server on :8080")
    if err := srv.Run(":8080"); err != nil {
        log.Fatalf("❌ Failed to start server: %v", err)
    }
}