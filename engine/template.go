package engine

import "github.com/nurtidev/predictor/pricer"

func checkTemplate(candles []*pricer.Candle, size int) bool {
	switch size {
	case 3:
		// Для трех свечей: средняя свеча должна быть другого цвета
		return isDifferentColor(candles[0], candles[1]) &&
			isDifferentColor(candles[2], candles[1])
	case 4:
		// Для четырех свечей: две средние свечи должны быть одного цвета, отличного от первой и четвертой
		return isDifferentColor(candles[0], candles[1]) &&
			isSameColor(candles[1], candles[2]) &&
			isDifferentColor(candles[3], candles[1])
	case 5:
		// Для пяти свечей: три средние свечи должны быть одного цвета, отличного от первой и пятой
		return isDifferentColor(candles[0], candles[1]) &&
			isSameColor(candles[1], candles[2]) &&
			isSameColor(candles[2], candles[3]) &&
			isDifferentColor(candles[4], candles[1])
	default:
		// Не подходит ни под один шаблон
		return false
	}
}

// isSameColor проверяет, одного ли цвета две свечи
func isSameColor(c1, c2 *pricer.Candle) bool {
	return c1.Color == c2.Color
}

// isDifferentColor проверяет, разного ли цвета две свечи
func isDifferentColor(c1, c2 *pricer.Candle) bool {
	return c1.Color != c2.Color
}

func getTemplateCandle(candles []*pricer.Candle, size int) (*pricer.Candle, bool) {
	switch size {
	case 3:
		return candles[1], true
	case 4:
		result := make([]*pricer.Candle, 2)
		result[0] = candles[1]
		result[1] = candles[2]
		return mergeCandles(result)
	case 5:
		result := make([]*pricer.Candle, 3)
		result[0] = candles[1]
		result[1] = candles[2]
		result[2] = candles[3]
		return mergeCandles(result)
	default:
		return &pricer.Candle{}, false
	}
}

func mergeCandles(candles []*pricer.Candle) (*pricer.Candle, bool) {
	if len(candles) < 2 {
		return &pricer.Candle{}, false // Не достаточно свечей для объединения
	}

	minLow := candles[0].Low
	maxHigh := candles[0].High
	totalVolume := candles[0].Volume

	for i, candle := range candles[1:] {
		if candle.Low < minLow {
			minLow = candle.Low
		}
		if candle.High > maxHigh {
			maxHigh = candle.High
		}
		totalVolume += candle.Volume
		if i == len(candles)-2 { // Проверяем, что все свечи одного цвета
			return &pricer.Candle{
				Market:    candles[0].Market,
				Timeframe: candles[0].Timeframe,
				Color:     candles[0].Color,
				Time:      candles[0].Time,
				Open:      candles[0].Open,
				Close:     candles[len(candles)-1].Close,
				Low:       minLow,
				High:      maxHigh,
				Volume:    totalVolume,
			}, true
		}
	}

	return &pricer.Candle{}, false
}
