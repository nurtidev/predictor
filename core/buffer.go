package core

import (
	"fmt"
	"github.com/nurtidev/predictor/pricer"
	"time"
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

func NewBuffer(size int) *Buffer {
	return &Buffer{
		Size:    size,
		Candles: make([]*pricer.Candle, 0),
		Status:  WaitTemplateCandlesStatus, // Начальное состояние
	}
}

func (buf *Buffer) Alert(candle *pricer.Candle) {
	color := ""
	switch buf.Candle.Color {
	case pricer.ColorRed:
		color = pricer.Red
	case pricer.ColorGreen:
		color = pricer.Green
	default:
		color = pricer.Reset // No color or default terminal color
	}

	fmt.Printf("%sBreakdown found!\t Time: %s\t Template time: %s\t %s\n",
		color,
		time.Unix(candle.Time, 0).UTC().Format("2006-01-02 15:04:05"),
		time.Unix(buf.Candle.Time, 0).UTC().Format("2006-01-02 15:04:05"),
		pricer.Reset,
	)
}
