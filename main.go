package main

import (
	"fmt"
	"github.com/nurtidev/predictor/core"
	"github.com/nurtidev/predictor/pricer"
	"log"
)

func main() {
	candles, err := pricer.LoadCandlesFromFile("./data/btc_usdt_5m_3d.csv")
	if err != nil {
		log.Fatal(err)
	}

	mng, err := core.NewManager(&core.Param{
		BreakdownPercent: 0.2,
		BreakdownMinSize: 3,
		BreakdownMaxSize: 10,
		MotionMinSize:    2,
		MotionMaxSize:    10,
	})
	if err != nil {
		log.Fatal(err)
	}

	for i, candle := range candles {
		candle.Idx = i
		if err = mng.ProcessCandle(candle); err != nil {
			log.Fatal(err)
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

}
