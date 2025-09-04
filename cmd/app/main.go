package main

import (
    "context"
    "fmt"
    "log"
    "time"
    "cryptofyanalytics/internal/models"
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

    fmt.Println("✅ Database connection established successfully!")

    // создаем тестовые данные
    testPrice := &models.CryptoPrice{
        CoinID:       "bitcoin",
        Symbol:       "btc",
        PriceUSD:     50000.12345678,
        MarketCapUSD: 1.1e12,  // 1.1 трн
        Volume24hUSD: 2.5e10,  // 25 млрд
        LastUpdated:  time.Now(),
    }

    // сохраняем данные в базу
    ctx := context.Background()
    err = repo.InsertPrice(ctx, testPrice)
    if err != nil {
        log.Printf("⚠️ Failed to insert test data: %v", err)
    } else {
        fmt.Println("✅ Test data inserted successfully!")
    }

    // читаем 
    prices, err := repo.GetPrices(ctx, "bitcoin", 5)
    if err != nil {
        log.Printf("⚠️ Failed to get prices: %v", err)
    } else {
        fmt.Printf("✅ Retrieved %d prices:\n", len(prices))
        for i, price := range prices {
            fmt.Printf("  %d. %s: $%.2f (updated at %s)\n", 
                i+1, price.Symbol, price.PriceUSD, 
                price.LastUpdated.Format("2006-01-02 15:04:05"))
        }
    }

    fmt.Println("🎉 Application finished successfully!")
}