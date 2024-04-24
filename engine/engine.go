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
		Metrics: &Metrics{},
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

	if err := e.collectFractals(candle, e.cfg.Filters.Fractal.Size); err != nil {
		return err
	}

	for _, buf := range e.buffers {
		if err := buf.Scan(candle); err != nil {
			return err
		}

		if buf.status == WaitAlert && e.isConfirmed(buf) {
			buf.Alert(candle)
			buf.status = Done
			e.Metrics.SignalsCount++
		}
	}

	return nil
}

func (e *Engine) isConfirmed(buf *Buffer) bool {
	return e.isFractalConfirmed(buf)
}

func (e *Engine) isFractalConfirmed(buf *Buffer) bool {
	candles := make([]*pricer.Candle, 0)
	candles = append(candles, buf.motion.Candles...)
	candles = append(candles, buf.breakdown.Candles...)

	switch buf.template.Candle.Color {
	case pricer.ColorRed:
		if len(e.fractal.High) < 3 {
			return false
		}
		for i, v := range e.fractal.High {
			if isExistByTime(v, candles) && i > 1 && i < len(e.fractal.High)-1 {
				idx := i
				if e.fractal.High[idx].High > e.fractal.High[idx-1].High && e.fractal.High[idx].High > e.fractal.High[idx+1].High {
					return true
				}
			}
		}

	case pricer.ColorGreen:
		if len(e.fractal.Low) < 3 {
			return false
		}

		for i, v := range e.fractal.Low {
			if isExistByTime(v, candles) && i > 1 && i < len(e.fractal.Low)-1 {
				idx := i
				if e.fractal.Low[idx].Low < e.fractal.Low[idx-1].Low && e.fractal.Low[idx].Low < e.fractal.Low[idx+1].Low {
					return true
				}
			}
		}

	}

	return false
}

func isExistByTime(candle *pricer.Candle, candles []*pricer.Candle) bool {
	for _, v := range candles {
		if candle.Time == v.Time {
			return true
		}
	}
	return false
}
