package core

import "github.com/nurtidev/predictor/pricer"

func (buf *Buffer) checkTemplate() bool {
	switch buf.Size {
	case 3:
		// Для трех свечей: средняя свеча должна быть другого цвета
		return isDifferentColor(buf.Candles[0], buf.Candles[1]) &&
			isDifferentColor(buf.Candles[2], buf.Candles[1])
	case 4:
		// Для четырех свечей: две средние свечи должны быть одного цвета, отличного от первой и четвертой
		return isDifferentColor(buf.Candles[0], buf.Candles[1]) &&
			isSameColor(buf.Candles[1], buf.Candles[2]) &&
			isDifferentColor(buf.Candles[3], buf.Candles[1])
	case 5:
		// Для пяти свечей: три средние свечи должны быть одного цвета, отличного от первой и пятой
		return isDifferentColor(buf.Candles[0], buf.Candles[1]) &&
			isSameColor(buf.Candles[1], buf.Candles[2]) &&
			isSameColor(buf.Candles[2], buf.Candles[3]) &&
			isDifferentColor(buf.Candles[4], buf.Candles[1])
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
