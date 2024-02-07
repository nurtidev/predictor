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
func sendMsgTelegram(str string) {
	bot, err := tgbotapi.NewBotAPI("6809455254:AAFEZlurUlCq8YBC7mtnc-a_TvgSF1mCs8U")
	if err != nil {
		log.Panic(err)
	}

	arr := make([]int64, 0)
	arr = append(arr, 1811775131)
	arr = append(arr, 1025535666)

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

	//str := fmt.Sprintf("💣 BOOOOOM 💣 \n Market: %s\n Timeframe: %s\n Breakdown found!\n Time: %s\n Template time: %s\n \n",
	//	candle.Market,
	//	candle.Timeframe,
	//	time.Unix(candle.Time, 0).UTC().Format("2006-01-02 15:04:05"),
	//	time.Unix(buf.Candle.Time, 0).UTC().Format("2006-01-02 15:04:05"),
	//)
	//
	//sendMsgTelegram(str)

	fmt.Printf("%sBreakdown found!\t Time: %s\t Template time: %s\t %s\n",
		color,
		time.Unix(candle.Time, 0).UTC().Format("2006-01-02 15:04:05"),
		time.Unix(buf.Candle.Time, 0).UTC().Format("2006-01-02 15:04:05"),
		pricer.Reset,
	)
}
