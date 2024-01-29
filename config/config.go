package config

type Config struct {
	Mode       string
	Market     []string
	Timeframe  []string
	Telegram   *Telegram
	BufferSize string
}

type Telegram struct {
	Users    []string
	BotToken string
}
