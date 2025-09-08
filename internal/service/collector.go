package service

import (
    "context"
    "log"
    "time"

    "cryptofyanalytics/internal/fetcher"
    "cryptofyanalytics/internal/models"
    "cryptofyanalytics/internal/storage"
)

// структура для сбора данных о криптовалютах
type CryptoCollector struct {
    fetcher  *fetcher.CoinGeckoClient
    storage  storage.Repository
    interval time.Duration
}

//  создает новый экземпляр сборщика данных
func NewCryptoCollector(fetcher *fetcher.CoinGeckoClient, storage storage.Repository, interval time.Duration) *CryptoCollector {
    return &CryptoCollector{
        fetcher:  fetcher,
        storage:  storage,
        interval: interval,
    }
}

// сбор данных по расписанию
func (cc *CryptoCollector) Start(ctx context.Context) {
    ticker := time.NewTicker(cc.interval)
    defer ticker.Stop()

    //  делаем первый сбор
    cc.CollectPrices(ctx)

    for {
        select {
        case <-ticker.C:
            cc.CollectPrices(ctx)
        case <-ctx.Done():
            log.Println("Collector stopped")
            return
        }
    }
}

// получает и сохраняет текущие цены
func (cc *CryptoCollector) CollectPrices(ctx context.Context) {
    log.Println("Starting data collection...")

    coins := []string{"bitcoin", "ethereum", "cardano", "solana"}
    prices, err := cc.fetcher.FetchMultiplePrices(coins)
    if err != nil {
        log.Printf("Error fetching prices: %v", err)
        return
    }

    // сохраняем каждую цену в базу данных
    for coinID, price := range prices {
        cryptoPrice := &models.CryptoPrice{
            CoinID:      coinID,
            Symbol:      cc.getSymbol(coinID),
            PriceUSD:    price,
            LastUpdated: time.Now(),
        }

        if err := cc.storage.InsertPrice(ctx, cryptoPrice); err != nil {
            log.Printf("Error saving %s price: %v", coinID, err)
        } else {
            log.Printf("Saved %s price: $%.2f", coinID, price)
        }
    }

    log.Println("Data collection completed")
}

// возвращает символ монеты по её ID
func (cc *CryptoCollector) getSymbol(coinID string) string {
    symbols := map[string]string{
        "bitcoin":  "btc",
        "ethereum": "eth",
        "cardano":  "ada",
        "solana":   "sol",
    }
    return symbols[coinID]
}