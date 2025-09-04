package storage 

import( "context"
		"cryptofy-analytics/internal/models"
	)

type Repository interface {
	//сохраняет цену крипты в БД
	InsertPrice(ctx context.Context, price *models>CryptoPrice) error

	//возвращает последние цуены для указанной монеты
	GetPrice(ctx context.Context, coinID string, limit int) ([]*models.CryptoPrice, error)

	Close() error
}

// соединение с БД.