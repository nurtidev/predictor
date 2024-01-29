package pricer

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

const (
	ColorRed   = "RED"
	ColorGreen = "GREEN"
)

type Candle struct {
	Color  string
	Time   int64
	Open   float64
	Close  float64
	Low    float64
	High   float64
	Volume float64
}

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
	if openPrice < closePrice {
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
