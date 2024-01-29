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
	Status         string
	Pool           []*pricer.Candle // место для общего хранения свеч
	Template       *Buffer          // место для хранения свеч среди которых будем искать шаблонную свечу
	Motion         *Buffer          // место для хранения свеч которые будут подтверждать движение
	Breakdown      *Buffer          // место для хранения свеч которые будут подтверждать пробой уровня
	Lifetime       int
	TemplateCount  int
	MotionCount    int
	BreakdownCount int
}

type Buffer struct {
	Candles      []*pricer.Candle
	Size         int
	ResetCounter int
}

func New(cfg *config.Config) (*Manager, error) {
	return &Manager{
		Status: WaitTemplateCandlesStatus,
		Template: &Buffer{
			Candles: make([]*pricer.Candle, 0),
			Size:    4,
		},
		Motion: &Buffer{
			Candles: make([]*pricer.Candle, 0),
			Size:    3,
		},
		Breakdown: &Buffer{
			Candles: make([]*pricer.Candle, 0),
			Size:    3,
		},
	}, nil
}

func Analysis(mng *Manager, candle *pricer.Candle) error {
	switch mng.Status {
	case WaitTemplateCandlesStatus:
		return mng.WaitTemplateCandles(candle)
	case WaitMotionCandlesStatus:
		return mng.WaitMotionCandles(candle)
	case WaitBreakdownCandlesStatus:
		return mng.WaitBreakdownCandles(candle)
	default:
		return errors.New("unknown manager status")
	}
}
