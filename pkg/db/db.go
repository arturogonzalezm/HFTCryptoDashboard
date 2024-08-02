package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type TickerData struct {
	EventType           string `json:"e"`
	EventTime           int64  `json:"E"`
	Symbol              string `json:"s"`
	PriceChange         string `json:"p"`
	PriceChangePercent  string `json:"P"`
	WeightedAvgPrice    string `json:"w"`
	FirstTradePrice     string `json:"x"`
	LastPrice           string `json:"c"`
	LastQuantity        string `json:"Q"`
	BestBidPrice        string `json:"b"`
	BestBidQuantity     string `json:"B"`
	BestAskPrice        string `json:"a"`
	BestAskQuantity     string `json:"A"`
	OpenPrice           string `json:"o"`
	HighPrice           string `json:"h"`
	LowPrice            string `json:"l"`
	Volume              string `json:"v"`
	QuoteVolume         string `json:"q"`
	StatisticsOpenTime  int64  `json:"O"`
	StatisticsCloseTime int64  `json:"C"`
	FirstTradeID        int64  `json:"F"`
	LastTradeID         int64  `json:"L"`
	TotalTrades         int64  `json:"n"`
}

var db *sql.DB

func InitDB(dataSourceName string) {
	var err error
	db, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	}
}

func InsertTickerData(ticker TickerData) error {
	if ticker.EventType == "" || ticker.EventTime == 0 || ticker.Symbol == "" {
		log.Printf("Invalid ticker data: %+v", ticker)
		return nil
	}

	query := `INSERT INTO ticker_data (
		event_type, event_time, symbol, price_change, price_change_percent,
		weighted_avg_price, first_trade_price, last_price, last_quantity,
		best_bid_price, best_bid_quantity, best_ask_price, best_ask_quantity,
		open_price, high_price, low_price, volume, quote_volume,
		statistics_open_time, statistics_close_time, first_trade_id,
		last_trade_id, total_trades
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
		$15, $16, $17, $18, $19, $20, $21, $22, $23)`

	_, err := db.Exec(query,
		ticker.EventType, ticker.EventTime, ticker.Symbol, ticker.PriceChange,
		ticker.PriceChangePercent, ticker.WeightedAvgPrice, ticker.FirstTradePrice,
		ticker.LastPrice, ticker.LastQuantity, ticker.BestBidPrice, ticker.BestBidQuantity,
		ticker.BestAskPrice, ticker.BestAskQuantity, ticker.OpenPrice, ticker.HighPrice,
		ticker.LowPrice, ticker.Volume, ticker.QuoteVolume, ticker.StatisticsOpenTime,
		ticker.StatisticsCloseTime, ticker.FirstTradeID, ticker.LastTradeID, ticker.TotalTrades)
	return err
}
