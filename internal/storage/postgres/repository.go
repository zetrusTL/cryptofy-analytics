package postgres

import(
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
	"cryptofy-analytics/internal/models"
//дрова SQL
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Repository struct {
//подключение к БД
	db *sql.DB
}

func NewRepository(connectionString string) (*Repository, error) {
	//открываем соездинение с БД
	db, error := sql.Open("pgx", connectionString)
	if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
	}

	//проверка
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }

	log.Println("✅ Successfully connected to PostgreSQL database")
    return &Repository{db: db}, nil
}

//сохраняет цену в базу данных
func (r *Repository) InsertPrice(ctx context.Context, price *models.CryptoPrice) error {
    // защитa от SQL-инъекций
    query := `
        INSERT INTO crypto_prices 
        (coin_id, symbol, price_usd, market_cap_usd, volume_24h_usd, last_updated) 
        VALUES ($1, $2, $3, $4, $5, $6)
    `

    // идет запрос с параметрами
    _, err := r.db.ExecContext(ctx, query,
        price.CoinID,
        price.Symbol,
        price.PriceUSD,
        price.MarketCapUSD,
        price.Volume24hUSD,
        price.LastUpdated,
    )

    if err != nil {
        return fmt.Errorf("failed to insert price: %w", err)
    }

    return nil
}

// возвращает последние цены для монеты
func (r *Repository) GetPrices(ctx context.Context, coinID string, limit int) ([]*models.CryptoPrice, error) {
    query := `
        SELECT id, coin_id, symbol, price_usd, market_cap_usd, volume_24h_usd, last_updated, created_at
        FROM crypto_prices 
        WHERE coin_id = $1 
        ORDER BY last_updated DESC 
        LIMIT $2
    `

    // получаем строки
    rows, err := r.db.QueryContext(ctx, query, coinID, limit)
    if err != nil {
        return nil, fmt.Errorf("failed to query prices: %w", err)
    }
    defer rows.Close() 

    var prices []*models.CryptoPrice


    for rows.Next() {
        var price models.CryptoPrice
        // сканируем данные из строки в структуру
        err := rows.Scan(
            &price.ID,
            &price.CoinID,
            &price.Symbol,
            &price.PriceUSD,
            &price.MarketCapUSD,
            &price.Volume24hUSD,
            &price.LastUpdated,
            &price.CreatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan price: %w", err)
        }
        prices = append(prices, &price)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating rows: %w", err)
    }

    return prices, nil
}

// закрывает соединение с базой данных
func (r *Repository) Close() error {
    return r.db.Close()
}

