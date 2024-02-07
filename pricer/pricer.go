package pricer

import (
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
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

	fmt.Printf("%s %s\t Time: %s\t Market: %s\t Timeframe: %s\t Open: %.2f\t Close: %.2f\t Low: %.2f\t High: %.2f\t%s\n",
		color,
		header,
		time.UnixMilli(c.Time).UTC().Format(time.DateTime),
		c.Market,
		c.Timeframe,
		c.Open,
		c.Close,
		c.Low,
		c.High,
		Reset,
	)
}

type Candle struct {
	Idx       int
	Market    string
	Timeframe string
	Color     string
	Time      int64
	Open      float64
	Close     float64
	Low       float64
	High      float64
	Volume    float64
}

func parseCurrencyPairAndTimeframe(filename string) (string, string, error) {
	// Регулярное выражение для парсинга
	re := regexp.MustCompile(`(\w+)_(\w+)_(\w+)_(\w+)\.csv$`)
	matches := re.FindStringSubmatch(filename)

	if len(matches) < 5 {
		return "", "", fmt.Errorf("failed to parse filename: %s", filename)
	}

	// Собираем пару валют и таймфрейм
	currencyPair := matches[1] + "_" + matches[2] // BTC_USDT
	timeframe := matches[3]                       // 15m

	return currencyPair, timeframe, nil
}

func LoadCandlesFromFile(path string) ([]*Candle, error) {
	currencyPair, timeframe, err := parseCurrencyPairAndTimeframe(path)
	if err != nil {
		return nil, err
	}
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
		candle.Market = currencyPair
		candle.Timeframe = timeframe
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
