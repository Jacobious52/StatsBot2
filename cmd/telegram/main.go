package main

import (
	"fmt"
	"time"

	"github.com/Jacobious52/StatsBot2/pkg/commands"
	"github.com/Jacobious52/StatsBot2/pkg/format"
	"github.com/Jacobious52/StatsBot2/pkg/storage"
	log "github.com/sirupsen/logrus"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	tb "gopkg.in/tucnak/telebot.v2"
)

var telegramToken = kingpin.Flag("token", "telegram bot token").Envar("TBTOKEN").Required().String()
var dataStorePath = kingpin.Flag("db", "path to save and read store file").Default("/var/lib/statsbot/db.json").String()
var csvStoreDir = kingpin.Flag("csv", "dir to save and read csv files").Default("/tmp").String()
var logLevel = kingpin.Flag("log", "logging level to use").Default("info").String()

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

	// list of formmaters and filters to combine
	formatters := map[string]format.Formatter{
		"day":   format.Dayly,
		"month": format.Monthly,
		"hour":  format.Hourly,
		"week":  format.Weekly,
	}
	filters := []storage.MessageKeyFilter{
		"text",
		"edited",
		"photo",
		"sticker",
	}

	// combine them all making commands
	for formatterName, formatter := range formatters {
		// add one formatter without a filter
		commandManager.RegisterCommand(
			fmt.Sprintf("/%s", formatterName),
			&commands.Stats{
				Name:   formatterName,
				CSVDir: *csvStoreDir,
				Format: formatter,
			},
		)

		// add the formatters with all the filters
		for _, filter := range filters {
			key := fmt.Sprintf("%s_%s", formatterName, filter)
			commandManager.RegisterCommand(
				fmt.Sprintf("/%s", key),
				&commands.Stats{
					Name:   key,
					CSVDir: *csvStoreDir,
					Format: formatter,
					Filter: filter,
				},
			)
		}
	}

	log.Infoln("starting bot")
	bot.Start()
}
