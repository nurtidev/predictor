package engine

import (
	"fmt"
	"github.com/nurtidev/predictor/pricer"
)

func initFractal() *Fractal {
	return &Fractal{
		High:    make([]*pricer.Candle, 0),
		Low:     make([]*pricer.Candle, 0),
		Candles: make([]*pricer.Candle, 0),
	}
}

func (e *Engine) collectFractals(candle *pricer.Candle, windowSize int) error {
	if windowSize%2 == 0 || windowSize < 5 {
		return fmt.Errorf("windowSize must be an odd number and at least 5")
	}

	e.fractal.Candles = append(e.fractal.Candles, candle)

	if len(e.fractal.Candles) == windowSize {
		middleIndex := windowSize / 2
		middleCandle := e.fractal.Candles[middleIndex]

		isFractalHigh := true
		isFractalLow := true

		for i := 0; i < windowSize; i++ {
			if i != middleIndex {
				if middleCandle.High <= e.fractal.Candles[i].High {
					isFractalHigh = false
				}
				if middleCandle.Low >= e.fractal.Candles[i].Low {
					isFractalLow = false
				}
			}
		}

		if isFractalHigh {
			e.fractal.High = append(e.fractal.High, middleCandle)
		}

		if isFractalLow {
			e.fractal.Low = append(e.fractal.Low, middleCandle)
		}

		e.fractal.Candles = e.fractal.Candles[1:]
	}

	return nil
}
