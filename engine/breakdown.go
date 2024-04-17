package engine

import (
	"errors"
	"github.com/nurtidev/predictor/pricer"
)

// checkBreakdown проверяет условия пробоя
func (buf *Buffer) checkBreakdown(candle *pricer.Candle) error {
	if isDifferentColor(buf.template.Candle, candle) {
		if buf.isValidBreakdown() {
			buf.status = WaitAlert
			return nil
		}
		buf.status = Canceled
		return nil
	}

	buf.breakdown.Candles = append(buf.breakdown.Candles, candle)

	return nil
}

func (buf *Buffer) isValidBreakdown() bool {
	if len(buf.breakdown.Candles) < buf.breakdown.MinSize {
		return false
	}

	if len(buf.breakdown.Candles) > buf.breakdown.MaxSize {
		return false
	}

	switch buf.template.Candle.Color {
	case pricer.ColorRed:
		high, _ := getHighestPrice(buf.motion.Candles)
		low, _ := getLowestPrice(buf.breakdown.Candles)
		isBreakPercent := ((low/high)-1)*100 <= -1*buf.breakdown.Percent
		if isBreakPercent && isBreakTemplate(buf.template.Candle, buf.breakdown.Candles) {
			return true
		}
	case pricer.ColorGreen:
		high, _ := getHighestPrice(buf.breakdown.Candles)
		low, _ := getLowestPrice(buf.motion.Candles)
		isBreakPercent := ((high/low)-1)*100 >= 1.5*buf.breakdown.Percent
		if isBreakPercent && isBreakTemplate(buf.template.Candle, buf.breakdown.Candles) {
			return true
		}
	}
	return false
}

func isBreakTemplate(template *pricer.Candle, candles []*pricer.Candle) bool {
	switch template.Color {
	case pricer.ColorRed:
		for _, candle := range candles {
			if candle.Close < template.Close {
				return true
			}
		}
	case pricer.ColorGreen:
		for _, candle := range candles {
			if candle.Close > template.Close {
				return true
			}
		}
	}
	return false
}

func getHighestPrice(arr []*pricer.Candle) (float64, error) {
	if len(arr) == 0 {
		return 0, errors.New("empty candles arr")
	}
	high := arr[0].High
	for _, v := range arr {
		if high < v.High {
			high = v.High
		}
	}
	return high, nil
}

func getLowestPrice(arr []*pricer.Candle) (float64, error) {
	if len(arr) == 0 {
		return 0, errors.New("empty candles arr")
	}
	low := arr[0].Low
	for _, v := range arr {
		if low > v.Low {
			low = v.Low
		}
	}
	return low, nil
}
