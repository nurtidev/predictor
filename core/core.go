package core

import (
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

func NewManager(param *Param) (*Manager, error) {
	pool := make([]*Buffer, 0)
	pool = append(pool, NewBuffer(3))
	pool = append(pool, NewBuffer(4))
	pool = append(pool, NewBuffer(5))
	return &Manager{
		Map:              make(map[int]*pricer.Candle),
		Pool:             pool,
		Storage:          make([]*Buffer, 0),
		BreakdownMaxSize: param.BreakdownMaxSize,
		BreakdownMinSize: param.BreakdownMinSize,
		BreakdownPercent: param.BreakdownPercent,
		MotionMaxSize:    param.MotionMaxSize,
		MotionMinSize:    param.MotionMinSize,
	}, nil
}

func (mng *Manager) ProcessCandle(candle *pricer.Candle) error {
	mng.Map[mng.Idx] = candle
	mng.Idx++

	for _, buf := range mng.Pool {
		buf.Candles = append(buf.Candles, candle)
		if len(buf.Candles) == buf.Size {
			if buf.checkTemplate() {
				//buf.Candles[1].Print("Template candle")
				mng.Storage = append(mng.Storage, &Buffer{
					Status:    WaitMotionCandlesStatus,
					Size:      buf.Size,
					Lifetime:  buf.Lifetime,
					Candle:    buf.Candles[1],
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
