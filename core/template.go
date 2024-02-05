package core

import "github.com/nurtidev/predictor/pricer"

func (buf *Buffer) isTrending() bool {
	if buf.Size < 2 {
		// Недостаточно данных для определения тренда
		return false
	}

	trendColor := buf.Candles[buf.Size-2].Color // Предполагаем, что цвет последней свечи тренда совпадает с предпоследней
	trendCount := 0

	// Перебираем свечи в обратном порядке, начиная со второй с конца, и считаем количество свечей тренда подряд
	for i := buf.Size - 2; i >= 0; i-- {
		if buf.Candles[i].Color == trendColor {
			trendCount++
		} else {
			break
		}
	}

	// Проверяем, достаточно ли свечей тренда для установления факта наличия тренда
	// Здесь вы можете определить свои параметры, например, тренд считается если есть 3 свечи подряд одного цвета
	return trendCount >= 3
}

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
