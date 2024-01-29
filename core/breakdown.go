package core

import "github.com/nurtidev/predictor/pricer"

// checkBreakdown проверяет условия пробоя
func (mng *Manager) checkBreakdown() bool {
	if len(mng.Breakdown.Candles) < 2 {
		return false // Пробой требует минимум двух свечей
	}

	lastCandle := mng.Breakdown.Candles[len(mng.Breakdown.Candles)-1]
	direction := mng.getBreakdownDirection()

	if direction == "up" {
		return lastCandle.Close > mng.getTemplateOpenPrice()
	} else if direction == "down" {
		return lastCandle.Close < mng.getTemplateOpenPrice()
	}

	return false
}

// getBreakdownDirection определяет направление пробоя
func (mng *Manager) getBreakdownDirection() string {
	// Используем вторую свечу в массиве Breakdown для определения направления
	if len(mng.Breakdown.Candles) > 1 && mng.Breakdown.Candles[1].Color == pricer.ColorRed {
		return "down" // Красные свечи указывают на пробой вниз
	}
	return "up" // В противном случае пробой вверх (по умолчанию для зеленых свечей)
}

// getTemplateOpenPrice возвращает цену открытия шаблонной свечи
func (mng *Manager) getTemplateOpenPrice() float64 {
	// Используем вторую свечу в массиве Template как шаблонную
	if len(mng.Template.Candles) > 1 {
		return mng.Template.Candles[1].Open
	}
	return 0 // Возвращаем 0, если шаблонных свечей нет
}
