package models

import (
    "time"
)

//Мродель данных

type CryptoPrice struct {
    ID           int64     `json:"id"`           
    CoinID       string    `json:"coin_id"`       
    Symbol       string    `json:"symbol"`        
    PriceUSD     float64   `json:"price_usd"`  
    MarketCapUSD float64   `json:"market_cap_usd"` 
    Volume24hUSD float64   `json:"volume_24h_usd"` 
    LastUpdated  time.Time `json:"last_updated"` 
    CreatedAt    time.Time `json:"created_at"`    
}