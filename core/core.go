package core

import (
	"errors"
	"github.com/nurtidev/predictor/config"
	"github.com/nurtidev/predictor/pricer"
)

type Manager struct {
	Pool             []*Buffer
	Storage          []*Buffer
	Map              map[int]*pricer.Candle
	BreakdownPercent float64
	BreakdownMinSize int
	BreakdownMaxSize int
	MotionMinSize    int
	MotionMaxSize    int
	Idx              int
}

type Param struct {
	BreakdownPercent float64
	BreakdownMinSize int
	BreakdownMaxSize int
	MotionMinSize    int
	MotionMaxSize    int
}

func NewManager(cfg *config.Config) (*Manager, error) {
	pool := make([]*Buffer, 0)
	pool = append(pool, NewBuffer(cfg, 3))
	pool = append(pool, NewBuffer(cfg, 4))
	pool = append(pool, NewBuffer(cfg, 5))
	return &Manager{
		Map:              make(map[int]*pricer.Candle),
		Pool:             pool,
		Storage:          make([]*Buffer, 0),
		BreakdownMaxSize: cfg.Trade.Breakdown.MaxSize,
		BreakdownMinSize: cfg.Trade.Breakdown.MinSize,
		BreakdownPercent: cfg.Trade.Breakdown.Percent,
		MotionMaxSize:    cfg.Trade.Motion.MaxSize,
		MotionMinSize:    cfg.Trade.Motion.MinSize,
	}, nil
}

func (mng *Manager) ProcessCandle(candle *pricer.Candle) error {
	mng.Map[mng.Idx] = candle
	mng.Idx++

	for _, buf := range mng.Pool {
		buf.Candles = append(buf.Candles, candle)
		if len(buf.Candles) == buf.Size {
			if buf.checkTemplate() {
				template, success := buf.getTemplateCandle()
				if !success {
					return errors.New("can't set template candle")
				}
				mng.Storage = append(mng.Storage, &Buffer{
					cfg:       buf.cfg,
					Status:    WaitMotionCandlesStatus,
					Size:      buf.Size,
					Lifetime:  buf.Lifetime,
					Candle:    template,
					Candles:   buf.Candles,
					Motion:    &Motion{MinSize: mng.MotionMinSize, MaxSize: mng.MotionMaxSize, Candles: make([]*pricer.Candle, 0)},
					Breakdown: &Breakdown{MinSize: mng.BreakdownMinSize, MaxSize: mng.BreakdownMaxSize, Percent: mng.BreakdownPercent, Candles: make([]*pricer.Candle, 0)},
				})
			}
			buf.Candles = buf.Candles[1:]
		}
	}

	for _, buf := range mng.Storage {
		if err := mng.Scan(buf, candle); err != nil {
			return err
		}
	}

	return nil
}

func (mng *Manager) Scan(buf *Buffer, candle *pricer.Candle) error {
	switch buf.Status {
	case WaitMotionCandlesStatus:
		return buf.checkMotion(candle)
	case WaitBreakdownCandlesStatus:
		return buf.checkBreakdown(candle)
	}
	return nil
}

func mergeCandles(candles []*pricer.Candle) (*pricer.Candle, bool) {
	if len(candles) < 2 {
		return &pricer.Candle{}, false // Не достаточно свечей для объединения
	}

	minLow := candles[0].Low
	maxHigh := candles[0].High
	totalVolume := candles[0].Volume

	for i, candle := range candles[1:] {
		if candle.Low < minLow {
			minLow = candle.Low
		}
		if candle.High > maxHigh {
			maxHigh = candle.High
		}
		totalVolume += candle.Volume
		if i == len(candles)-2 { // Проверяем, что все свечи одного цвета
			return &pricer.Candle{
				Market:    candles[0].Market,
				Timeframe: candles[0].Timeframe,
				Color:     candles[0].Color,
				Time:      candles[0].Time,
				Open:      candles[0].Open,
				Close:     candles[len(candles)-1].Close,
				Low:       minLow,
				High:      maxHigh,
				Volume:    totalVolume,
			}, true
		}
	}

	return &pricer.Candle{}, false
}
