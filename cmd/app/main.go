package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "cryptofyanalytics/internal/fetcher"
    "cryptofyanalytics/internal/server"
    "cryptofyanalytics/internal/service"
    "cryptofyanalytics/internal/storage/postgres"
)

func main() {
    fmt.Println("🚀 Starting Crypto Analytics Application...")

    // строка подключения формат: postgres://username:password@host:port/database_name
    connStr := "postgres://crypto_user:crypto_password@localhost:5432/crypto_db?sslmode=disable"

    // cоздаем репозиторий для работы с базой данных
    repo, err := postgres.NewRepository(connStr)
    if err != nil {
        log.Fatalf("❌ Failed to create repository: %v", err)
    }
    defer repo.Close() 

    fmt.Println("✅ Database connection established successfully!")

   // создаем клиент для CoinGecko API
    geckoClient := fetcher.NewCoinGeckoClient()

    // создаем сервис для сбора данных (каждые 5 минут)
    collector := service.NewCryptoCollector(geckoClient, repo, 5*time.Minute)

    // создаем контекст для graceful shutdown
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Запускаем сбор данных в фоновой горутине
    go collector.Start(ctx)

    // создаем и запускаем HTTP-сервер
    srv := server.NewServer(repo)
    
    // запускаем сервер в отдельной горутине
    go func() {
        log.Println("🚀 Starting HTTP server on :8080")
        if err := srv.Run(":8080"); err != nil {
            log.Fatalf("❌ Failed to start server: %v", err)
        }
    }()

    // ожидание сигналов завершения
    waitForShutdown(cancel)
}

//  ожидает сигналов OS для graceful shutdown
func waitForShutdown(cancel context.CancelFunc) {
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    
    <-sigChan
    fmt.Println("\n🛑 Shutdown signal received...")
    cancel() // Останавливаем сборщик данных
    time.Sleep(1 * time.Second) // Даем время на завершение
    fmt.Println("✅ Application stopped gracefully")
}
