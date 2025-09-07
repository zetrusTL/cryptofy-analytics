package storage 

import( "context"
		"cryptofyanalytics/internal/models"
	)


	//Интерфейс хранилища 
type Repository interface {
	//сохраняет цену крипты в БД
	InsertPrice(ctx context.Context, price *models.CryptoPrice) error

	//возвращает последние цуены для указанной монеты
	GetPrices(ctx context.Context, coinID string, limit int) ([]*models.CryptoPrice, error)

	Close() error
}

// соединение с БД.