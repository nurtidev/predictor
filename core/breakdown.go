package core

import "github.com/nurtidev/predictor/pricer"

// checkBreakdown проверяет условия пробоя
func (buf *Buffer) checkBreakdown() bool {
	if len(buf.Breakdown) < 2 {
		return false // Пробой требует минимум двух свечей
	}

	lastCandle := buf.Breakdown[len(buf.Breakdown)-1]
	direction := buf.getBreakdownDirection()

	if direction == "up" {
		return lastCandle.Close > buf.getTemplateOpenPrice()
	} else if direction == "down" {
		return lastCandle.Close < buf.getTemplateOpenPrice()
	}

	return false
}

// getBreakdownDirection определяет направление пробоя
func (buf *Buffer) getBreakdownDirection() string {
	// Используем вторую свечу в массиве Breakdown для определения направления
	if len(buf.Breakdown) > 1 && buf.Breakdown[1].Color == pricer.ColorRed {
		return "down" // Красные свечи указывают на пробой вниз
	}
	return "up" // В противном случае пробой вверх (по умолчанию для зеленых свечей)
}

// getTemplateOpenPrice возвращает цену открытия шаблонной свечи
func (buf *Buffer) getTemplateOpenPrice() float64 {
	// Используем вторую свечу в массиве Template как шаблонную
	if len(buf.Candles) > 1 {
		return buf.Candles[1].Open
	}
	return 0 // Возвращаем 0, если шаблонных свечей нет
}
