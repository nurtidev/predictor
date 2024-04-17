package engine

import (
	"github.com/nurtidev/predictor/config"
	"github.com/nurtidev/predictor/pricer"
)

func initPools(cfg *config.Config) []*Pool {
	pools := make([]*Pool, 0)
	for _, v := range cfg.Trade.BufferSize {
		pools = append(pools, &Pool{
			size:    v,
			candles: make([]*pricer.Candle, 0),
		})
	}

	return pools
}

func initBuffer(cfg *config.Config, size int, candle *pricer.Candle, candles []*pricer.Candle) *Buffer {
	return &Buffer{
		status: WaitMotion,
		template: &Template{
			Candle:  candle,
			Size:    size,
			Candles: candles,
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
}

func (buf *Buffer) Scan(candle *pricer.Candle) error {
	switch buf.status {
	case WaitMotion:
		return buf.checkMotion(candle)
	case WaitBreakdown:
		return buf.checkBreakdown(candle)
	}
	return nil
}
