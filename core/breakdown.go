package core

import (
	"errors"
	"github.com/nurtidev/predictor/pricer"
)

// checkBreakdown проверяет условия пробоя
func (buf *Buffer) checkBreakdown(candle *pricer.Candle) error {
	if len(buf.Breakdown.Candles) > buf.Breakdown.MaxSize {
		buf.Status = DoneCandlesStatus
		return nil
	}

	if isDifferentColor(buf.Candle, candle) {
		if buf.isValidBreakdown() {
			buf.Alert(candle)
			buf.Status = AlertCandlesStatus
			return nil
		}
		buf.Status = DoneCandlesStatus
		return nil
	}

	buf.Breakdown.Candles = append(buf.Breakdown.Candles, candle)

	return nil
}

func (buf *Buffer) isValidBreakdown() bool {
	if len(buf.Breakdown.Candles) < buf.Breakdown.MinSize {
		return false
	}

	switch buf.Candle.Color {
	case pricer.ColorRed:
		high, _ := getHighestPrice(buf.Motion.Candles)
		low, _ := getLowestPrice(buf.Breakdown.Candles)
		isBreakPercent := ((low/high)-1)*100 <= -1*buf.Breakdown.Percent
		if isBreakPercent && isBreakTemplate(buf.Candle, buf.Breakdown.Candles) {
			return true
		}
	case pricer.ColorGreen:
		high, _ := getHighestPrice(buf.Breakdown.Candles)
		low, _ := getLowestPrice(buf.Motion.Candles)
		isBreakPercent := ((high/low)-1)*100 >= 1.5*buf.Breakdown.Percent
		if isBreakPercent && isBreakTemplate(buf.Candle, buf.Breakdown.Candles) {
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
