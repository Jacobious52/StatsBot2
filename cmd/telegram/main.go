package main

import (
	"time"

	"github.com/Jacobious52/StatsBot2/pkg/commands"
	"github.com/Jacobious52/StatsBot2/pkg/format"
	"github.com/Jacobious52/StatsBot2/pkg/storage"
	log "github.com/sirupsen/logrus"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	tb "gopkg.in/tucnak/telebot.v2"
)

var telegramToken = kingpin.Flag("token", "telegram bot token").Envar("TBTOKEN").Required().String()
var dataStorePath = kingpin.Flag("store", "path to save and read store file").Default("/usr/share/store/store.json").String()
var csvStoreDir = kingpin.Flag("csv", "dir to save and read csv files").Default("/usr/share/store/csv").String()
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
	commandManager := commands.NewCommandManager(bot, dataStore)
	log.Infoln("loading plugins")

	// security risk: disable unless debugging
	// commandManager.RegisterCommand("/dump", &commands.Dump{DBPath: *dataStorePath})

	// help triggers
	commandManager.RegisterCommand("/help", &commands.Help{Commander: commandManager})

	// event triggers
	commandManager.RegisterCommand(tb.OnText, &commands.Event{Type: "text"})
	commandManager.RegisterCommand(tb.OnSticker, &commands.Event{Type: "sticker"})
	commandManager.RegisterCommand(tb.OnPhoto, &commands.Event{Type: "photo"})
	commandManager.RegisterCommand(tb.OnEdited, &commands.Event{Type: "edited"})

	// stats triggers
	commandManager.RegisterCommand("/month",
		&commands.Stats{
			Name:   "month",
			CSVDir: *csvStoreDir,
			Format: format.Monthly,
		},
	)
	commandManager.RegisterCommand("/day",
		&commands.Stats{
			Name:   "day",
			CSVDir: *csvStoreDir,
			Format: format.Dayly,
		},
	)
	commandManager.RegisterCommand("/month_edited",
		&commands.Stats{
			Name:   "month_edited",
			CSVDir: *csvStoreDir,
			Format: format.Monthly,
			Filter: "edited",
		},
	)
	commandManager.RegisterCommand("/day_edited",
		&commands.Stats{
			Name:   "day_edited",
			CSVDir: *csvStoreDir,
			Format: format.Dayly,
			Filter: "edited",
		},
	)
	commandManager.RegisterCommand("/month_sticker",
		&commands.Stats{
			Name:   "month_sticker",
			CSVDir: *csvStoreDir,
			Format: format.Monthly,
			Filter: "sticker",
		},
	)
	commandManager.RegisterCommand("/day_sticker",
		&commands.Stats{
			Name:   "day_sticker",
			CSVDir: *csvStoreDir,
			Format: format.Dayly,
			Filter: "sticker",
		},
	)
	commandManager.RegisterCommand("/month_photo",
		&commands.Stats{
			Name:   "month_photo",
			CSVDir: *csvStoreDir,
			Format: format.Monthly,
			Filter: "photo",
		},
	)
	commandManager.RegisterCommand("/day_photo",
		&commands.Stats{
			Name:   "day_photo",
			CSVDir: *csvStoreDir,
			Format: format.Dayly,
			Filter: "photo",
		},
	)
	commandManager.RegisterCommand("/month_text",
		&commands.Stats{
			Name:   "month_text",
			CSVDir: *csvStoreDir,
			Format: format.Monthly,
			Filter: "text",
		},
	)
	commandManager.RegisterCommand("/day_text",
		&commands.Stats{
			Name:   "day_text",
			CSVDir: *csvStoreDir,
			Format: format.Dayly,
			Filter: "text",
		},
	)

	log.Infoln("starting bot")
	bot.Start()
}
