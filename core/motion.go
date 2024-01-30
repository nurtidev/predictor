package core

import (
	"github.com/nurtidev/predictor/pricer"
	"math"
)

// checkMotionSize проверяет размеры свечей в Motion
func (buf *Buffer) checkMotionSize() bool {
	templateSize := calculateSize(buf.Candles)
	motionSize := calculateSize(buf.Motion)

	// Условие размера зависит от количества шаблонных свечей
	multiplier := 1.5
	if buf.Size == 3 {
		multiplier = 2
	}

	return motionSize >= templateSize*multiplier
}

// calculateSize вычисляет суммарный размер свечей
func calculateSize(candles []*pricer.Candle) float64 {
	var size float64
	for _, candle := range candles {
		size += math.Abs(candle.Close - candle.Open)
	}
	return size
}

// initMotionFromTemplate инициализирует Motion последними свечами из Template одного цвета
func (buf *Buffer) initMotionFromTemplate() {
	if len(buf.Candles) == 0 {
		return
	}
	lastColor := buf.Candles[len(buf.Candles)-1].Color
	for i := len(buf.Candles) - 1; i >= 0; i-- {
		if buf.Candles[i].Color != lastColor {
			break
		}
		buf.Motion = append([]*pricer.Candle{buf.Candles[i]}, buf.Motion...)
	}
}
