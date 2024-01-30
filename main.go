package main

import (
	"fmt"
	"github.com/nurtidev/predictor/config"
	"github.com/nurtidev/predictor/core"
	"github.com/nurtidev/predictor/pricer"
	"log"
)

func main() {
	candles, err := pricer.LoadCandlesFromFile("./data/btc_usdt_5m_1d.csv")
	if err != nil {
		log.Fatal(err)
	}

	mng, err := core.New(&config.Config{})
	if err != nil {
		log.Fatal(err)
	}

	for _, candle := range candles {
		if err = core.Analysis(mng, candle); err != nil {
			log.Fatal(err)
		}
	}

	for _, buf := range mng.Buffers {
		fmt.Printf("Buffer size: %d\n", buf.Size)
		fmt.Printf("Template count: %d\n", buf.Stat.TemplateCount)
		fmt.Printf("Motion count: %d\n", buf.Stat.MotionCount)
		fmt.Printf("Breakdown count: %d\n", buf.Stat.BreakdownCount)
	}
}
