package main

import (
    "fmt"
    "log"

    "cryptofyanalytics/internal/fetcher"
)

func main() {
    client := fetcher.NewCoinGeckoClient()

    // Тестируем получение цены Bitcoin
    price, err := client.FetchCurrentPrice("bitcoin")
    if err != nil {
        log.Fatalf("Error fetching Bitcoin price: %v", err)
    }
    fmt.Printf("💰 Bitcoin price: $%.2f\n", price)

    // Тестируем получение цен нескольких монет
    coins := []string{"bitcoin", "ethereum", "cardano"}
    prices, err := client.FetchMultiplePrices(coins)
    if err != nil {
        log.Fatalf("Error fetching multiple prices: %v", err)
    }

    fmt.Println("\n📊 Multiple prices:")
    for coin, price := range prices {
        fmt.Printf("  %s: $%.2f\n", coin, price)
    }
}