package main

import (
	"log"

	"github.com/boltdb/bolt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/xokcton/saver-tg-bot/pkg/config"
	"github.com/xokcton/saver-tg-bot/pkg/pocket"
	"github.com/xokcton/saver-tg-bot/pkg/server"
	"github.com/xokcton/saver-tg-bot/pkg/storage"
	"github.com/xokcton/saver-tg-bot/pkg/storage/boltdb"
	"github.com/xokcton/saver-tg-bot/pkg/telegram"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	botApi, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatal(err)
	}
	botApi.Debug = true

	pocketClient, err := pocket.NewClient(cfg.PocketConsumerKey)
	if err != nil {
		log.Fatal(err)
	}

	db, err := initBolt(cfg.BoltDBFile)
	if err != nil {
		log.Fatal(err)
	}
	storage := boltdb.NewTokenStorage(db)

	bot := telegram.NewBot(botApi, pocketClient, cfg.AuthServerURL, storage, cfg.Messages)

	redirectServer := server.NewAuthServer(cfg.BotURL, storage, pocketClient)

	go func() {
		if err := redirectServer.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	if err := bot.Start(); err != nil {
		log.Fatal(err)
	}
}

func initBolt(db_file string) (*bolt.DB, error) {
	db, err := bolt.Open(db_file, 0600, nil)
	if err != nil {
		return nil, err
	}

	if err := db.Batch(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(storage.AccessTokens))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte(storage.RequestTokens))
		return err
	}); err != nil {
		return nil, err
	}

	return db, nil
}
