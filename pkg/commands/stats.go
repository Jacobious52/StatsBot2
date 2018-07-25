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

type Stats struct {
	Name   string
	CSVDir string
	Format format.Formatter
	Filter storage.MessageKeyFilter
}

func (m *Stats) Do(data storage.Model, info storage.MessageInfo) (interface{}, error) {
	filename := fmt.Sprint(m.CSVDir, "/", m.Name, "-", info.Chat, "-", time.Now(), ".csv")
	log.Debugln("creating csv at ", filename)
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	format.CSV(data[info.Chat], m.Name, m.Format, m.Filter, file)

	return &tb.Document{File: tb.FromDisk(filename), FileName: fmt.Sprint(m.Name, ".csv")}, nil
}
