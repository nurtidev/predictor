package core

import (
	"github.com/nurtidev/predictor/config"
	"github.com/nurtidev/predictor/pricer"
)

type CandleStatus string

var (
	WaitTemplateCandlesStatus  CandleStatus = "WAIT_TEMPLATE"
	WaitMotionCandlesStatus    CandleStatus = "WAIT_MOTION"
	WaitBreakdownCandlesStatus CandleStatus = "WAIT_BREAKDOWN"
	DoneCandlesStatus          CandleStatus = "DONE"
	AlertCandlesStatus         CandleStatus = "ALERT"
)

type Buffer struct {
	cfg       *config.Config
	Status    CandleStatus
	Size      int
	Lifetime  int
	Candle    *pricer.Candle
	Candles   []*pricer.Candle
	Motion    *Motion
	Breakdown *Breakdown
}

type Motion struct {
	MaxSize int
	MinSize int
	Candles []*pricer.Candle
}

type Breakdown struct {
	MaxSize int
	MinSize int
	Percent float64
	Candles []*pricer.Candle
}

func NewBuffer(cfg *config.Config, size int) *Buffer {
	return &Buffer{
		cfg:     cfg,
		Size:    size,
		Candles: make([]*pricer.Candle, 0),
		Status:  WaitTemplateCandlesStatus, // Начальное состояние
	}
}

func (buf *Buffer) getTemplateCandle() (*pricer.Candle, bool) {
	switch buf.Size {
	case 3:
		return buf.Candles[1], true
	case 4:
		candles := make([]*pricer.Candle, 2)
		candles[0] = buf.Candles[1]
		candles[1] = buf.Candles[2]
		return mergeCandles(candles)
	case 5:
		candles := make([]*pricer.Candle, 3)
		candles[0] = buf.Candles[1]
		candles[1] = buf.Candles[2]
		candles[2] = buf.Candles[3]
		return mergeCandles(candles)
	default:
		return &pricer.Candle{}, false
	}
}
