CREATE TABLE IF NOT EXISTS crypto_prices (
    id SERIAL PRIMARY KEY,
    coin_id VARCHAR(20) NOT NULL,
    symbol VARCHAR(10) NOT NULL,
    price_usd DECIMAL(20, 8) NOT NULL,
    market_cap_usd DECIMAL(25, 8),
    volume_24h_usd DECIMAL(25, 8),
    last_updated TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE crypto_prices IS 'Хранит исторические данные о ценах криптовалют';
COMMENT ON COLUMN crypto_prices.coin_id IS 'Внутренний ID монеты на CoinGecko';
COMMENT ON COLUMN crypto_prices.symbol IS 'Тикер монеты (например, btc)';
COMMENT ON COLUMN crypto_prices.price_usd IS 'Текущая цена в USD';
COMMENT ON COLUMN crypto_prices.last_updated IS 'Время обновления данных от API';