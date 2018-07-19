package commands

import (
	"fmt"
	"os"
	"time"

	"github.com/Jacobious52/StatsBot2/pkg/format"
	"github.com/Jacobious52/StatsBot2/pkg/storage"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
)

type Day struct {
	CSVDir string
}

func (m *Day) Do(data storage.Model, info storage.MessageInfo) (interface{}, error) {
	filename := fmt.Sprint(m.CSVDir, "/", info.Chat, "-", time.Now(), ".csv")
	log.Debugln("creating csv at ", filename)
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	format.CSV(data[info.Chat], format.Dayly, file)

	return &tb.Document{File: tb.FromDisk(filename), FileName: "day.csv"}, nil
}
