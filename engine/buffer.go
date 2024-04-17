package engine

import (
	"github.com/nurtidev/predictor/config"
	"github.com/nurtidev/predictor/pricer"
)

func initBuffers(cfg *config.Config) []*Buffer {

	buffers := make([]*Buffer, len(cfg.Trade.BufferSize))

	for _, bufferSize := range cfg.Trade.BufferSize {
		buffer := &Buffer{
			template: &Template{
				Size:    bufferSize,
				Candles: make([]*pricer.Candle, 0),
			},
			motion: &Motion{
				MinSize: cfg.Trade.Motion.MinSize,
				MaxSize: cfg.Trade.Motion.MaxSize,
				Candles: make([]*pricer.Candle, 0),
			},
			breakdown: &Breakdown{
				Percent: cfg.Trade.Breakdown.Percent,
				MinSize: cfg.Trade.Breakdown.MinSize,
				MaxSize: cfg.Trade.Breakdown.MaxSize,
				Candles: make([]*pricer.Candle, 0),
			},
		}

		buffers = append(buffers, buffer)
	}

	return buffers
}
