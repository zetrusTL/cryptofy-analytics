package fetcher

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

// CoinGeckoClient представляет клиент для работы с CoinGecko API
type CoinGeckoClient struct {
    BaseURL    string
    HTTPClient *http.Client
}

// создает и возвращает новый экземпляр клиента
func NewCoinGeckoClient() *CoinGeckoClient {
    return &CoinGeckoClient{
        BaseURL: "https://api.coingecko.com/api/v3",
        HTTPClient: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
}

// представляет структуру ответа от API
type PriceResponse struct {
    Bitcoin struct {
        USD float64 `json:"usd"`
    } `json:"bitcoin"`
    Ethereum struct {
        USD float64 `json:"usd"`
    } `json:"ethereum"`
}

// получает текущую цену для указанной монеты
func (cg *CoinGeckoClient) FetchCurrentPrice(coinID string) (float64, error) {
    //  URL запрос
    url := fmt.Sprintf("%s/simple/price?ids=%s&vs_currencies=usd", 
        cg.BaseURL, coinID)

    // HTTP-запрос
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return 0, fmt.Errorf("failed to create request: %w", err)
    }

    req.Header.Set("Accept", "application/json")

    // выполняем запрос
    resp, err := cg.HTTPClient.Do(req)
    if err != nil {
        return 0, fmt.Errorf("failed to make request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return 0, fmt.Errorf("API returned non-200 status: %d", resp.StatusCode)
    }

    // читаем
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return 0, fmt.Errorf("failed to read response body: %w", err)
    }

    // парсим JSON ответ
    var data map[string]map[string]float64
    if err := json.Unmarshal(body, &data); err != nil {
        return 0, fmt.Errorf("failed to parse JSON: %w", err)
    }

    // извлекаем цену из ответа
    if coinData, exists := data[coinID]; exists {
        if price, exists := coinData["usd"]; exists {
            return price, nil
        }
    }

    return 0, fmt.Errorf("price not found for coin: %s", coinID)
}

// получает цены для нескольких монет сразу
func (cg *CoinGeckoClient) FetchMultiplePrices(coinIDs []string) (map[string]float64, error) {

    var ids string
    for i, id := range coinIDs {
        if i > 0 {
            ids += ","
        }
        ids += id
    }

    url := fmt.Sprintf("%s/simple/price?ids=%s&vs_currencies=usd", 
        cg.BaseURL, ids)

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }

    req.Header.Set("Accept", "application/json")

    resp, err := cg.HTTPClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to make request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API returned non-200 status: %d", resp.StatusCode)
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response body: %w", err)
    }

    var data map[string]map[string]float64
    if err := json.Unmarshal(body, &data); err != nil {
        return nil, fmt.Errorf("failed to parse JSON: %w", err)
    }

    result := make(map[string]float64)
    for coinID, coinData := range data {
        if price, exists := coinData["usd"]; exists {
            result[coinID] = price
        }
    }

    return result, nil
}