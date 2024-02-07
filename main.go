package main

import (
	"fmt"
	"github.com/nurtidev/predictor/config"
	"github.com/nurtidev/predictor/core"
	"github.com/nurtidev/predictor/pricer"
	"log"
	"sync"
	"time"
)

// chat_id: 1811775131 Nurtilek
// chat_id: 1025535666 Dima

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	//client := &pricer.Binance{}
	//
	//fmt.Println(cfg)
	//
	//candle, err := client.GetCandleFromBinance("BTCUSDT", "10m")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//candle.Print("test")

	switch cfg.Mode {
	case "simulation":
		if err = simulation(cfg); err != nil {
			log.Fatal(err)
		}
	case "realtime":
		if err = realtime(cfg); err != nil {
			log.Fatal(err)
		}
	}
}

func simulation(cfg *config.Config) error {
	candles, err := pricer.LoadCandlesFromFile(cfg.Filepath)
	if err != nil {
		log.Fatal(err)
	}

	mng, err := core.NewManager(&core.Param{
		BreakdownPercent: cfg.Trade.Breakdown.Percent,
		BreakdownMinSize: cfg.Trade.Breakdown.MinSize,
		BreakdownMaxSize: cfg.Trade.Breakdown.MaxSize,
		MotionMinSize:    cfg.Trade.Breakdown.MinSize,
		MotionMaxSize:    cfg.Trade.Breakdown.MaxSize,
	})
	if err != nil {
		return err
	}

	for i, candle := range candles {
		candle.Idx = i
		if err = mng.ProcessCandle(candle); err != nil {
			return err
		}
	}

	templateCount, motionCount, breakdownCount, alertCount := 0, 0, 0, 0

	for _, v := range mng.Storage {
		if len(v.Motion.Candles) > 0 {
			motionCount++
		}
		if len(v.Breakdown.Candles) > 0 {
			breakdownCount++
		}
		if v.Status == core.AlertCandlesStatus {
			alertCount++
		}
	}

	templateCount = len(mng.Storage)

	fmt.Printf("templateCount: %d\n", templateCount)
	fmt.Printf("motionCount: %d\n", motionCount)
	fmt.Printf("breakdownCount: %d\n", breakdownCount)
	fmt.Printf("alertCount: %d\n", alertCount)

	return nil
}

func realtime(cfg *config.Config) error {
	var wg sync.WaitGroup

	client := pricer.Binance{}
	for _, pair := range cfg.Market {
		for _, timeframe := range cfg.Timeframe {
			wg.Add(1)
			go func(pair, timeframe string) {
				defer wg.Done()

				mng, err := core.NewManager(&core.Param{
					BreakdownPercent: cfg.Trade.Breakdown.Percent,
					BreakdownMinSize: cfg.Trade.Breakdown.MinSize,
					BreakdownMaxSize: cfg.Trade.Breakdown.MaxSize,
					MotionMinSize:    cfg.Trade.Breakdown.MinSize,
					MotionMaxSize:    cfg.Trade.Breakdown.MaxSize,
				})
				if err != nil {
					log.Fatal(err)
				}

				duration, err := pricer.GetDuration(timeframe)
				if err != nil {
					log.Fatal(err)
				}

				idx := 0
				candle, err := client.GetCandleFromBinance(pair, timeframe)
				if err != nil {
					log.Fatal(err)
				}
				candle.Print("START CANDLE! ")
				candle.Idx = idx
				if err = mng.ProcessCandle(candle); err != nil {
					log.Fatal(err)
				}
				idx++

				ticker := time.NewTicker(duration)
				defer ticker.Stop()

				for range ticker.C {
					candle, err = client.GetCandleFromBinance(pair, timeframe)
					if err != nil {
						log.Fatal(err)
					}

					candle.Idx = idx
					if err = mng.ProcessCandle(candle); err != nil {
						log.Fatal(err)
					}

					idx++
				}

			}(pair, timeframe)
		}
	}

	wg.Wait()

	return nil
}
