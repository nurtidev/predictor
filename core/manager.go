package core

import (
	"fmt"
	"github.com/nurtidev/predictor/pricer"
	"time"
)

func (mng *Manager) reset() {
	mng.Status = WaitTemplateCandlesStatus
	mng.Template.Candles = nil
	mng.Motion.Candles = nil
	mng.Breakdown.Candles = nil
}

func (mng *Manager) WaitTemplateCandles(candle *pricer.Candle) error {
	mng.Template.Candles = append(mng.Template.Candles, candle)
	if len(mng.Template.Candles) == mng.Template.Size {
		if mng.checkTemplate() {
			mng.TemplateCount++
			//fmt.Printf("TEMPLATE! Candle time: %s\n", time.Unix(candle.Time, 0).Format(time.DateTime))
			//fmt.Println(mng.Template.Candles[1])
			mng.Status = WaitMotionCandlesStatus
			return nil
		} else {
			mng.Template.Candles = mng.Template.Candles[1:]
		}

	}
	return nil
}

func (mng *Manager) WaitMotionCandles(candle *pricer.Candle) error {
	// Если Motion пуст, начинаем с последних свечей одного цвета из Template
	if len(mng.Motion.Candles) == 0 {
		mng.initMotionFromTemplate()
	}

	// Проверяем цвет текущей свечи с первой свечой в Motion
	if candle.Color != mng.Motion.Candles[0].Color {
		// Если цвета разные, проверяем размер и очищаем Motion
		if mng.checkMotionSize() {
			mng.MotionCount++
			mng.Status = WaitBreakdownCandlesStatus
			return nil
		}
		mng.reset()
		return nil
	}

	// Добавляем свечу в Motion, если цвет совпадает
	mng.Motion.Candles = append(mng.Motion.Candles, candle)

	return nil
}

func (mng *Manager) WaitBreakdownCandles(candle *pricer.Candle) error {
	// Если Breakdown пуст, начинаем с текущей свечи
	if len(mng.Breakdown.Candles) == 0 {
		mng.Breakdown.Candles = append(mng.Breakdown.Candles, candle)
		return nil
	}

	// Проверяем цвет текущей свечи с первой свечой в Breakdown
	if candle.Color != mng.Breakdown.Candles[0].Color {
		// Если цвета разные, проверяем условия пробоя и очищаем Breakdown
		if mng.checkBreakdown() {
			mng.BreakdownCount++
			fmt.Printf("BREAKDOWN! Candle time: %s\n", time.Unix(candle.Time, 0).UTC().Format(time.DateTime))
			fmt.Printf("Candle: %v\n", candle)
		}
		mng.reset()
		return nil
	}

	// Добавляем свечу в Breakdown, если цвет совпадает
	mng.Breakdown.Candles = append(mng.Breakdown.Candles, candle)

	return nil
}
