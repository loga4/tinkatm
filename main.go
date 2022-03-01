package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-redis/redis"
	tgapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron/v3"

	atm "tinkoffbot/pkg/atm"
	"tinkoffbot/pkg/config"
	"tinkoffbot/pkg/logger/zaplog"
)

var ctx = context.Background()

var atms = make(map[string]map[string]atm.ATM)

func main() {
	cfgPath := flag.String("config", "./config.yml", "")
	flag.Parse()

	cfg, err := config.New(*cfgPath)
	if err != nil {
		log.Panic(err)
	}

	rds := redis.NewClient(&redis.Options{
		Addr: cfg.Redis,
		DB:   0,
	})
	if err := rds.Ping().Err(); err != nil {
		log.Fatalln("failed connect to redis server", "dsn", cfg.Redis, "err", err)
		return
	}

	logger, err := zaplog.New("[traderbot]", "", "", "")
	if err != nil {
		log.Fatalln("failed create logger", "err", err)
		return
	}

	var client = &http.Client{}

	tgbotapi, err := tgapi.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		log.Panic(err)
	}
	tgbotapi.Debug = false

	tinkClient := atm.NewTinkoffClient(cfg, client, logger, tgbotapi, rds, atms)

	// tib := tgbot.NewTgInfoBot(
	// 	logger,
	// 	tgbotapi,
	// )
	// go tib.Listen()

	c := cron.New(cron.WithSeconds(), cron.WithLocation(time.UTC))

	tinkClient.SendRequest("USD")
	tinkClient.SendRequest("EUR")

	c.AddFunc("0 * * * * *", func() {
		tinkClient.SendRequest("USD")
		tinkClient.SendRequest("EUR")
	})
	go c.Start()

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
}
