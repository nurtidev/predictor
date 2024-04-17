package engine

import (
	"errors"
	"github.com/nurtidev/predictor/config"
	"github.com/nurtidev/predictor/pricer"
)

func New(cfg *config.Config) *Engine {
	return &Engine{
		buffers:  make([]*Buffer, 0),
		pools:    initPools(cfg),
		fractals: make([]*Fractal, 0),
		metrics:  &Trade{},
	}
}

func (e *Engine) Process(candle *pricer.Candle) error {

	for _, pool := range e.pools {
		pool.candles = append(pool.candles, candle)

		if len(pool.candles) == pool.size && checkTemplate(pool.candles, pool.size) {
			template, success := getTemplateCandle(pool.candles, pool.size)
			if !success {
				return errors.New("can't set template candle")
			}

			buf := initBuffer(e.cfg, pool.size, template, pool.candles)
			e.buffers = append(e.buffers, buf)
		}
		pool.candles = pool.candles[1:]
	}

	for _, buf := range e.buffers {
		if err := buf.Scan(candle); err != nil {
			return err
		}

		if buf.status == WaitAlert && e.isConfirmed(buf) {
			buf.Alert(candle)
		}
	}

	return nil
}

func (e *Engine) isConfirmed(buf *Buffer) bool {
	return true
}

func (e *Engine) isFractalConfirmed(buf *Buffer) bool {
	return false
}
