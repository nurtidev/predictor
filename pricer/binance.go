package pricer

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Binance struct {
	mu sync.Mutex
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

func (c *CandleBinance) Parse(response []interface{}, market string, timeFrame string) {
	c.Market = market
	c.TickInterval = timeFrame
	c.OpenTime = response[0].(float64)
	c.Open, _ = strconv.ParseFloat(response[1].(string), 64)
	c.High, _ = strconv.ParseFloat(response[2].(string), 64)
	c.Low, _ = strconv.ParseFloat(response[3].(string), 64)
	c.Close, _ = strconv.ParseFloat(response[4].(string), 64)
	c.Volume, _ = strconv.ParseFloat(response[5].(string), 64)
	if c.Close > c.Open {
		c.Color = ColorGreen
	}
	if c.Close < c.Open {
		c.Color = ColorRed
	}
}

type Candlestick [][]interface{}

func (b *Binance) GetCandleFromBinance(pair, timeframe string) (*Candle, error) {
	start, end, err := syncTime(timeframe)
	if err != nil {
		return nil, err
	}
	limit := 1
	url := fmt.Sprintf("https://api.binance.com/api/v3/klines?symbol=%s&interval=%s&limit=%d&startTime=%d&endTime=%d", pair, timeframe, limit, start, end)
	var response Candlestick
	b.mu.Lock()
	time.Sleep(1 * time.Second)
	if err := SendHTTPRequest(http.MethodGet, url, nil, nil, &response); err != nil {
		return nil, err
	}
	b.mu.Unlock()

	for _, v := range response {
		var candle CandleBinance
		candle.Parse(v, pair, timeframe)
		return &Candle{
			Market:    candle.Market,
			Timeframe: candle.TickInterval,
			Color:     candle.Color,
			Time:      int64(candle.OpenTime),
			Open:      candle.Open,
			Close:     candle.Close,
			Low:       candle.Low,
			High:      candle.High,
			Volume:    candle.Volume,
		}, nil
	}

	return nil, errors.New("candle not found")
}

func syncTime(timeframe string) (int64, int64, error) {
	duration, err := GetDuration(timeframe)
	if err != nil {
		return 0, 0, err
	}

	now := time.Now()

	end := now.Truncate(duration)
	start := end.Add(-1 * duration)

	return start.UnixMilli(), end.UnixMilli(), nil
}

func GetDuration(timeframe string) (time.Duration, error) {
	switch timeframe {
	case "1m":
		return 1 * time.Minute, nil
	case "5m":
		return 5 * time.Minute, nil
	case "10m":
		return 10 * time.Minute, nil
	case "15m":
		return 15 * time.Minute, nil
	default:
		return 0, errors.New("unknown timeframe")
	}
}
