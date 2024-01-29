package core

import "github.com/nurtidev/predictor/pricer"

func (mng *Manager) checkTemplate() bool {
	switch len(mng.Template.Candles) {
	case 3:
		// Для трех свечей: средняя свеча должна быть другого цвета
		return mng.isDifferentColor(mng.Template.Candles[0], mng.Template.Candles[1]) &&
			mng.isDifferentColor(mng.Template.Candles[2], mng.Template.Candles[1])
	case 4:
		// Для четырех свечей: две средние свечи должны быть одного цвета, отличного от первой и четвертой
		return mng.isDifferentColor(mng.Template.Candles[0], mng.Template.Candles[1]) &&
			mng.isSameColor(mng.Template.Candles[1], mng.Template.Candles[2]) &&
			mng.isDifferentColor(mng.Template.Candles[3], mng.Template.Candles[1])
	case 5:
		// Для пяти свечей: три средние свечи должны быть одного цвета, отличного от первой и пятой
		return mng.isDifferentColor(mng.Template.Candles[0], mng.Template.Candles[1]) &&
			mng.isSameColor(mng.Template.Candles[1], mng.Template.Candles[2]) &&
			mng.isSameColor(mng.Template.Candles[2], mng.Template.Candles[3]) &&
			mng.isDifferentColor(mng.Template.Candles[4], mng.Template.Candles[1])
	default:
		// Не подходит ни под один шаблон
		return false
	}
}

// isSameColor проверяет, одного ли цвета две свечи
func (mng *Manager) isSameColor(c1, c2 *pricer.Candle) bool {
	return c1.Color == c2.Color
}

// isDifferentColor проверяет, разного ли цвета две свечи
func (mng *Manager) isDifferentColor(c1, c2 *pricer.Candle) bool {
	return c1.Color != c2.Color
}
