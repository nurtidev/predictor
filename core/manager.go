package core

import (
	"fmt"
	"github.com/nurtidev/predictor/pricer"
	"time"
)

func (mng *Manager) reset() {
	mng.Status = WaitTemplateCandlesStatus
	for _, v := range mng.Buffers {
		v.reset()
	}
}

func (buf *Buffer) reset() {
	buf.isMain = false
	buf.Status = WaitTemplateCandlesStatus
	buf.Candle = &pricer.Candle{}
	buf.Candles = make([]*pricer.Candle, 0)
	buf.Motion = make([]*pricer.Candle, 0)
	buf.Breakdown = make([]*pricer.Candle, 0)
}

func (mng *Manager) SetMain(buf *Buffer) {
	for _, b := range mng.Buffers {
		if b.Size == buf.Size {
			b.isMain = true
		} else {
			b.isMain = false
		}
	}
}

func (buf *Buffer) WaitTemplateCandles(candle *pricer.Candle) error {
	buf.Candles = append(buf.Candles, candle)
	if len(buf.Candles) == buf.Size {
		if buf.checkTemplate() {
			fmt.Printf("TEMPLATE! Candle time: %s\n", time.Unix(buf.Candles[1].Time, 0).UTC().Format(time.DateTime))
			buf.Stat.TemplateCount++
			buf.Status = WaitMotionCandlesStatus
			buf.Candle = buf.Candles[1]
		} else {
			buf.Candles = buf.Candles[1:]
		}
	}
	return nil
}

func (buf *Buffer) WaitMotionCandles(candle *pricer.Candle) error {
	// Если Motion пуст, начинаем с последних свечей одного цвета из Template
	if len(buf.Motion) == 0 {
		buf.initMotionFromTemplate()
	}

	// Проверяем цвет текущей свечи с первой свечой в Motion
	if candle.Color != buf.Motion[0].Color {
		// Если цвета разные, проверяем размер и очищаем Motion
		if buf.checkMotionSize() {
			buf.Stat.MotionCount++
			buf.Status = WaitBreakdownCandlesStatus
			return nil
		}
		buf.reset()
		return nil
	}
	// Добавляем свечу в Motion, если цвет совпадает
	buf.Motion = append(buf.Motion, candle)

	return nil

}

func (buf *Buffer) WaitBreakdownCandles(candle *pricer.Candle) error {
	// Если Breakdown пуст, начинаем с текущей свечи
	if len(buf.Breakdown) == 0 {
		buf.Breakdown = append(buf.Breakdown, candle)
		return nil
	}

	// Проверяем цвет текущей свечи с первой свечой в Breakdown
	if candle.Color != buf.Breakdown[0].Color {
		// Если цвета разные, проверяем условия пробоя и очищаем Breakdown
		if buf.checkBreakdown() {
			buf.Stat.BreakdownCount++
			fmt.Printf("BREAKDOWN! Candle time: %s\n", time.Unix(candle.Time, 0).UTC().Format(time.DateTime))
			fmt.Printf("Candle: %v\n", candle)
		}
		buf.reset()
		return nil
	}

	// Добавляем свечу в Breakdown, если цвет совпадает
	buf.Breakdown = append(buf.Breakdown, candle)
	return nil

}
