package engine

import (
	"fmt"
	"github.com/nurtidev/predictor/pricer"
)

func initFractal() *Fractal {
	return &Fractal{
		HighIdx: 0,
		LowIdx:  0,
		High:    make(map[int]*pricer.Candle),
		Low:     make(map[int]*pricer.Candle),
		Candles: make([]*pricer.Candle, 0),
	}
}

func (e *Engine) collectFractals(candle *pricer.Candle) error {
	e.fractal.Candles = append(e.fractal.Candles, candle)

	if len(e.fractal.Candles) == 5 {
		middleCandle := e.fractal.Candles[2]

		isFractalHigh := middleCandle.High > e.fractal.Candles[1].High && middleCandle.High > e.fractal.Candles[0].High &&
			middleCandle.High > e.fractal.Candles[3].High && middleCandle.High > e.fractal.Candles[4].High
		isFractalLow := middleCandle.Low < e.fractal.Candles[1].Low && middleCandle.Low < e.fractal.Candles[0].Low &&
			middleCandle.Low < e.fractal.Candles[3].Low && middleCandle.Low < e.fractal.Candles[4].Low

		if isFractalHigh {
			fmt.Println("found high fractal")
			e.fractal.High[e.fractal.HighIdx] = middleCandle
			e.fractal.HighIdx++
		}

		if isFractalLow {
			fmt.Println("found low fractal")
			e.fractal.Low[e.fractal.LowIdx] = middleCandle
			e.fractal.LowIdx++
		}

	}

	return nil

}
