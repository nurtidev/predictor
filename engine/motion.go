package engine

import "github.com/nurtidev/predictor/pricer"

func (buf *Buffer) checkMotion(candle *pricer.Candle) error {
	if len(buf.motion.Candles) == 0 && isSameColor(buf.template.Candle, candle) {
		buf.status = Canceled
		return nil
	}

	if len(buf.motion.Candles) > buf.motion.MaxSize {
		buf.status = Canceled
		return nil
	}

	if isSameColor(buf.template.Candle, candle) {
		if buf.isValidMotion() {
			buf.breakdown.Candles = append(buf.breakdown.Candles, candle)
			buf.status = WaitBreakdown
			return nil
		} else {
			buf.status = Canceled
			return nil
		}
	}

	buf.motion.Candles = append(buf.motion.Candles, candle)
	return nil
}

func (buf *Buffer) isValidMotion() bool {
	if len(buf.motion.Candles) < buf.motion.MinSize {
		return false
	}
	candle := buf.motion.Candles[len(buf.motion.Candles)-1]
	switch buf.template.Candle.Color {
	case pricer.ColorRed:

		return candle.Close > buf.template.Candle.Open
	case pricer.ColorGreen:

		return candle.Close < buf.template.Candle.Open
	}
	return false
}
