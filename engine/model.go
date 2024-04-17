package engine

import (
	"github.com/nurtidev/predictor/config"
	"github.com/nurtidev/predictor/pricer"
)

type bufferStatus string

const (
	WaitMotion    bufferStatus = "wait_motion"
	WaitBreakdown bufferStatus = "wait_breakdown"
	WaitAlert     bufferStatus = "wait_alert"
	Canceled      bufferStatus = "canceled" // todo: возможно у нас будут разные типы cancel для статистики
	Done          bufferStatus = "done"
)

type Engine struct {
	cfg     *config.Config
	buffers []*Buffer
	pools   []*Pool
	fractal *Fractal
	metrics *Trade
}

type Pool struct {
	size    int
	candles []*pricer.Candle
}

type Trade struct {
	Balance      float64
	SignalsCount int
}

type Fractal struct {
	HighIdx int
	LowIdx  int
	High    map[int]*pricer.Candle
	Low     map[int]*pricer.Candle
	Candles []*pricer.Candle
}

type Buffer struct {
	status    bufferStatus
	template  *Template
	motion    *Motion
	breakdown *Breakdown
}

type Template struct {
	Size    int
	Candle  *pricer.Candle
	Candles []*pricer.Candle
}

type Breakdown struct {
	Percent float64
	MinSize int
	MaxSize int
	Candles []*pricer.Candle
}

type Motion struct {
	MinSize int
	MaxSize int
	Candles []*pricer.Candle
}
