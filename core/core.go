package core

import (
	"errors"
	"github.com/nurtidev/predictor/config"
	"github.com/nurtidev/predictor/pricer"
)

const (
	WaitTemplateCandlesStatus  = "WAIT_TEMPLATE"
	WaitMotionCandlesStatus    = "WAIT_MOTION"
	WaitBreakdownCandlesStatus = "WAIT_BREAKDOWN"
)

type Manager struct {
	Status  string
	Pool    []*pricer.Candle // место для общего хранения свеч
	Buffers []*Buffer
}

type Statistic struct {
	TemplateCount  int
	MotionCount    int
	BreakdownCount int
}

type Buffer struct {
	isMain           bool
	Status           string
	Size             int
	Lifetime         int
	Candle           *pricer.Candle   // здесь мы будем хранить шаблонную свечу
	Candles          []*pricer.Candle // здесь все свечи включая шаблонную
	Motion           []*pricer.Candle
	Breakdown        []*pricer.Candle
	Stat             *Statistic
	BreakdownPercent float64
}

func New(cfg *config.Config) (*Manager, error) {
	return &Manager{
		Status: WaitTemplateCandlesStatus,
		Buffers: []*Buffer{
			{
				Status:    WaitTemplateCandlesStatus,
				Size:      3,
				Lifetime:  10,
				Candle:    &pricer.Candle{},
				Candles:   make([]*pricer.Candle, 0),
				Motion:    make([]*pricer.Candle, 0),
				Breakdown: make([]*pricer.Candle, 0),
				Stat: &Statistic{
					TemplateCount:  0,
					MotionCount:    0,
					BreakdownCount: 0,
				},
			},
			{
				Size:      4,
				Candle:    &pricer.Candle{},
				Candles:   make([]*pricer.Candle, 0),
				Motion:    make([]*pricer.Candle, 0),
				Breakdown: make([]*pricer.Candle, 0),
				Stat: &Statistic{
					TemplateCount:  0,
					MotionCount:    0,
					BreakdownCount: 0,
				},
			},
			{
				Size:      5,
				Candle:    &pricer.Candle{},
				Candles:   make([]*pricer.Candle, 0),
				Motion:    make([]*pricer.Candle, 0),
				Breakdown: make([]*pricer.Candle, 0),
				Stat: &Statistic{
					TemplateCount:  0,
					MotionCount:    0,
					BreakdownCount: 0,
				},
			},
		},
	}, nil
}

func Analysis(mng *Manager, candle *pricer.Candle) error {
	//mng.Pool = append(mng.Pool, candle)
	for _, buf := range mng.Buffers {
		switch buf.Status {
		case WaitTemplateCandlesStatus:
			return buf.WaitTemplateCandles(candle)
		case WaitMotionCandlesStatus:
			return buf.WaitMotionCandles(candle)
		case WaitBreakdownCandlesStatus:
			return buf.WaitBreakdownCandles(candle)
		default:
			return errors.New("unknown manager status")
		}
	}
	return errors.New("empty buffers")
}
