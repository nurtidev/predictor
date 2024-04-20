package main

import (
	"github.com/nurtidev/predictor/config"
	"github.com/nurtidev/predictor/core"
	"github.com/nurtidev/predictor/engine"
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

	//mng, _ := core.NewManager(cfg)
	//
	//for _, c := range candles {
	//	err = mng.ProcessCandle(c)
	//	if err != nil {
	//		return err
	//	}
	//}

	eng := engine.New(cfg)

	for _, candle := range candles {
		if err = eng.Process(candle); err != nil {
			return err
		}
	}

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

				mng, err := core.NewManager(cfg)
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
