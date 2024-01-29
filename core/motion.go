package core

import (
	"github.com/nurtidev/predictor/pricer"
	"math"
)

// checkMotionSize проверяет размеры свечей в Motion
func (mng *Manager) checkMotionSize() bool {
	templateSize := mng.calculateSize(mng.Template.Candles)
	motionSize := mng.calculateSize(mng.Motion.Candles)

	// Условие размера зависит от количества шаблонных свечей
	multiplier := 1.5
	if len(mng.Template.Candles) == 3 {
		multiplier = 2
	}

	return motionSize >= templateSize*multiplier
}

// calculateSize вычисляет суммарный размер свечей
func (mng *Manager) calculateSize(candles []*pricer.Candle) float64 {
	var size float64
	for _, candle := range candles {
		size += math.Abs(candle.Close - candle.Open)
	}
	return size
}

// initMotionFromTemplate инициализирует Motion последними свечами из Template одного цвета
func (mng *Manager) initMotionFromTemplate() {
	if len(mng.Template.Candles) == 0 {
		return
	}

	lastColor := mng.Template.Candles[len(mng.Template.Candles)-1].Color
	for i := len(mng.Template.Candles) - 1; i >= 0; i-- {
		if mng.Template.Candles[i].Color != lastColor {
			break
		}
		mng.Motion.Candles = append([]*pricer.Candle{mng.Template.Candles[i]}, mng.Motion.Candles...)
	}
}
