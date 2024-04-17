package engine

import (
	"errors"
	"github.com/nurtidev/predictor/config"
	"github.com/nurtidev/predictor/pricer"
)

func New(cfg *config.Config) *Engine {
	return &Engine{
		cfg:     cfg,
		buffers: make([]*Buffer, 0),
		pools:   initPools(cfg),
		fractal: initFractal(),
		metrics: &Trade{},
	}
}

func (e *Engine) Process(candle *pricer.Candle) error {

	for _, pool := range e.pools {
		pool.candles = append(pool.candles, candle)

		if len(pool.candles) == pool.size {
			if checkTemplate(pool.candles, pool.size) {
				template, success := getTemplateCandle(pool.candles, pool.size)
				if !success {
					return errors.New("can't set template candle")
				}

				buf := initBuffer(e.cfg, pool.size, template, pool.candles)

				e.buffers = append(e.buffers, buf)
			}

			pool.candles = pool.candles[1:]
		}
	}

	if err := e.collectFractals(candle); err != nil {
		return err
	}

	for _, buf := range e.buffers {
		if err := buf.Scan(candle); err != nil {
			return err
		}

		if buf.status == WaitAlert && e.isConfirmed(buf) {
			buf.Alert(candle)
			buf.status = Done
		}
	}

	return nil
}

func (e *Engine) isConfirmed(buf *Buffer) bool {
	return true
}

func (e *Engine) isFractalConfirmed(buf *Buffer) bool {
	switch buf.template.Candle.Color {
	case pricer.ColorRed:
		if len(e.fractal.High) < 2 {
			return false
		}

		for k, v := range e.fractal.High {
			if v.Time == buf.template.Candle.Time && k > 1 {
				idx := k
				if e.fractal.High[idx].High > e.fractal.High[idx-1].High {
					return true
				}
			}
		}

	case pricer.ColorGreen:
		if len(e.fractal.Low) < 2 {
			return false
		}

		for k, v := range e.fractal.High {
			if v.Time == buf.template.Candle.Time && k > 1 {
				idx := k
				if e.fractal.Low[idx].Low > e.fractal.Low[idx-1].Low {
					return true
				}
			}
		}

	}

	return false
}
