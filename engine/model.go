package engine

import "github.com/nurtidev/predictor/pricer"

type Engine struct {
	buffers  []*Buffer
	pool     []*pricer.Candle
	fractals []*Fractal
	metrics  *Trade
}

type Trade struct {
	Balance      float64
	SignalsCount int
}

type Fractal struct {
	Candles []*pricer.Candle
}

type Buffer struct {
	template  *Template
	motion    *Motion
	breakdown *Breakdown
}

type Template struct {
	Size    int
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
