package core

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nurtidev/predictor/pricer"
	"log"
	"time"
)

// chat_id: 1811775131 Nurtilek
// chat_id: 1025535666 Dima
func (buf *Buffer) sendMsgTelegram(str string) {
	bot, err := tgbotapi.NewBotAPI(buf.cfg.Telegram.Token)
	if err != nil {
		log.Panic(err)
	}

	arr := make([]int64, len(buf.cfg.Telegram.Users))
	for i := range arr {
		arr[i] = int64(buf.cfg.Telegram.Users[i])
	}

	for _, chatId := range arr {
		msg := tgbotapi.NewMessage(chatId, str)
		bot.Send(msg)
	}
}

func (buf *Buffer) Alert(candle *pricer.Candle) {
	color := ""
	switch buf.Candle.Color {
	case pricer.ColorRed:
		color = pricer.Red
	case pricer.ColorGreen:
		color = pricer.Green
	default:
		color = pricer.Reset // No color or default terminal color
	}

	if buf.cfg.Telegram.Enable {
		str := fmt.Sprintf("💣 BOOOOOM 💣 \n Market: %s\n Timeframe: %s\n Breakdown found!\n Time: %s\n Template time: %s\n \n",
			candle.Market,
			candle.Timeframe,
			time.UnixMilli(candle.Time).UTC().Format(time.DateTime),
			time.UnixMilli(buf.Candle.Time).UTC().Format(time.DateTime),
		)

		buf.sendMsgTelegram(str)
	}

	fmt.Printf("%sBreakdown found!\t Time: %s\t Template time: %s\t %s\n",
		color,
		time.UnixMilli(candle.Time).UTC().Format(time.DateTime),
		time.UnixMilli(buf.Candle.Time).UTC().Format(time.DateTime),
		pricer.Reset,
	)
}
