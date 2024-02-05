package pricer

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	ColorRed   = "RED"
	ColorGreen = "GREEN"
)

// ANSI escape codes for colors
const (
	Red   = "\033[31m"
	Green = "\033[32m"
	Reset = "\033[0m"
)

func (c *Candle) Print(header string) {
	color := ""
	switch c.Color {
	case ColorRed:
		color = Red
	case ColorGreen:
		color = Green
	default:
		color = Reset // No color or default terminal color
	}

	fmt.Printf("%s %s\t Time: %s\t Open: %.2f\t Close: %.2f\t Low: %.2f\t High: %.2f\t%s\n",
		color,
		header,
		time.Unix(c.Time, 0).UTC().Format(time.DateTime),
		c.Open,
		c.Close,
		c.Low,
		c.High,
		Reset,
	)
}

type Candle struct {
	Idx    int
	Color  string
	Time   int64
	Open   float64
	Close  float64
	Low    float64
	High   float64
	Volume float64
}

type CandleBinance struct {
	Market                   string
	TickInterval             string
	Color                    string
	OpenTime                 float64
	Open                     float64
	High                     float64
	Low                      float64
	Close                    float64
	Volume                   float64
	CloseTime                float64
	QuoteAssetVolume         float64
	NumberOfTrades           int
	TakerBuyBaseAssetVolume  float64
	TakerBuyQuoteAssetVolume float64
	Ignore                   float64
}

type Candlestick [][]interface{}

func LoadCandlesFromFile(path string) ([]*Candle, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var candles []*Candle
	for _, record := range records[1:] { // Пропускаем заголовок
		candle, err := parseRecord(record)
		if err != nil {
			return nil, err
		}
		candles = append(candles, candle)
	}

	return candles, nil
}

func parseRecord(record []string) (*Candle, error) {
	if len(record) != 6 {
		return nil, fmt.Errorf("invalid record length: %v", record)
	}

	time, err := strconv.ParseInt(record[0], 10, 64)
	if err != nil {
		return nil, err
	}
	openPrice, err := strconv.ParseFloat(record[1], 64)
	if err != nil {
		return nil, err
	}
	closePrice, err := strconv.ParseFloat(record[2], 64)
	if err != nil {
		return nil, err
	}
	low, err := strconv.ParseFloat(record[3], 64)
	if err != nil {
		return nil, err
	}
	high, err := strconv.ParseFloat(record[4], 64)
	if err != nil {
		return nil, err
	}
	volume, err := strconv.ParseFloat(record[5], 64)
	if err != nil {
		return nil, err
	}

	color := ColorGreen
	if openPrice > closePrice {
		color = ColorRed
	}

	return &Candle{
		Color:  color,
		Time:   time,
		Open:   openPrice,
		Close:  closePrice,
		Low:    low,
		High:   high,
		Volume: volume,
	}, nil
}
