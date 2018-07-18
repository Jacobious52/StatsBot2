package main

import (
	"time"

	"github.com/Jacobious52/StatsBot2/pkg/commands"
	"github.com/Jacobious52/StatsBot2/pkg/storage"
	log "github.com/sirupsen/logrus"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	tb "gopkg.in/tucnak/telebot.v2"
)

var telegramToken = kingpin.Flag("token", "telegram bot token").Envar("TBTOKEN").Required().String()
var dataStorePath = kingpin.Flag("store", "path to save and read store file").Default("/usr/share/store/store.json").String()
var logLevel = kingpin.Flag("level", "logging level to use").Default("info").String()

func main() {
	kingpin.Parse()
	level, _ := log.ParseLevel(*logLevel)
	log.SetLevel(level)

	// create the bot
	bot, err := tb.NewBot(tb.Settings{
		Token: *telegramToken,
		Poller: &tb.LongPoller{
			Timeout: 10 * time.Second,
		},
	})
	if err != nil {
		log.Fatalf("failed to created bot: %v", err)
	}
	log.Infoln("created bot!")

	// load data and start syncing to disk
	dataStore := storage.NewDataStore(*dataStorePath)
	err = dataStore.Load()
	if err != nil {
		log.Warningf("data not loaded. %v", err)
	}
	dataStore.StartSync(nil)

	// load plugins
	log.Infoln("loading plugins")
	commandManager := commands.NewCommandManager(bot, dataStore)
	commandManager.RegisterCommand("/help", new(commands.Help))
	commandManager.RegisterCommand("/start", new(commands.Start))
	commandManager.RegisterCommand(tb.OnText, new(commands.OnText))
	commandManager.RegisterCommand(tb.OnSticker, new(commands.OnSticker))
	commandManager.RegisterCommand(tb.OnPhoto, new(commands.OnPhoto))
	commandManager.RegisterCommand(tb.OnEdited, new(commands.OnEdited))
	commandManager.RegisterCommand("/dump", &commands.Dump{DBPath: *dataStorePath})

	log.Infoln("starting bot")
	bot.Start()
}
