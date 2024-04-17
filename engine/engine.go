package engine

import (
	"github.com/nurtidev/predictor/config"
	"github.com/nurtidev/predictor/pricer"
)

func New(cfg *config.Config) *Engine {
	return &Engine{
		buffers:  initBuffers(cfg),
		pool:     make([]*pricer.Candle, 0),
		fractals: make([]*Fractal, 0),
		metrics:  &Trade{},
	}
}

func (e *Engine) Process(candle *pricer.Candle) error {
	for _, b := range e.buffers {
		b.template.Candles = append(b.template.Candles, candle)
		if len(b.template.Candles) == b.template.Size {
			// todo: make
		}

		b.template.Candles = b.template.Candles[1:]
	}

	return nil
}
