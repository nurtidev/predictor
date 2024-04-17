package core

import (
	"github.com/nurtidev/predictor/pricer"
)

func (buf *Buffer) checkMotion(candle *pricer.Candle) error {
	if len(buf.Motion.Candles) == 0 && isSameColor(buf.Candle, candle) {
		buf.Status = DoneCandlesStatus
		return nil
	}

	if len(buf.Motion.Candles) > buf.Motion.MaxSize {
		buf.Status = DoneCandlesStatus
		return nil
	}

	if isSameColor(buf.Candle, candle) {
		if buf.isValidMotion() {
			pricer.IdxMotion++
			buf.Breakdown.Candles = append(buf.Breakdown.Candles, candle)
			buf.Status = WaitBreakdownCandlesStatus
			return nil
		} else {
			buf.Status = DoneCandlesStatus
			return nil
		}
	}

	buf.Motion.Candles = append(buf.Motion.Candles, candle)
	return nil
}

func (buf *Buffer) isValidMotion() bool {
	if len(buf.Motion.Candles) < buf.Motion.MinSize {
		return false
	}
	candle := buf.Motion.Candles[len(buf.Motion.Candles)-1]
	switch buf.Candle.Color {
	case pricer.ColorRed:

		return candle.Close > buf.Candle.Open
	case pricer.ColorGreen:

		return candle.Close < buf.Candle.Open
	}
	return false
}
