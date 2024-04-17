package engine

import (
	"fmt"
	"github.com/nurtidev/predictor/pricer"
	"time"
)

func (buf *Buffer) Alert(candle *pricer.Candle) {
	color := ""
	switch buf.template.Candle.Color {
	case pricer.ColorRed:
		color = pricer.Red
	case pricer.ColorGreen:
		color = pricer.Green
	default:
		color = pricer.Reset // No color or default terminal color
	}

	//if buf.cfg.Telegram.Enable {
	//	str := fmt.Sprintf("💣 BOOOOOM 💣 \n Market: %s\n Timeframe: %s\n Breakdown found!\n Time: %s\n Template time: %s\n \n",
	//		candle.Market,
	//		candle.Timeframe,
	//		time.UnixMilli(candle.Time).UTC().Format(time.DateTime),
	//		time.UnixMilli(buf.Candle.Time).UTC().Format(time.DateTime),
	//	)
	//
	//	buf.sendMsgTelegram(str)
	//}

	fmt.Printf("%sBreakdown found!\t Time: %s\t Template time: %s\t %s\n",
		color,
		time.UnixMilli(candle.Time).UTC().Format(time.DateTime),
		time.UnixMilli(buf.template.Candle.Time).UTC().Format(time.DateTime),
		pricer.Reset,
	)
}
