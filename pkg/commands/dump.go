package commands

import (
	"bytes"
	"encoding/json"

	"github.com/Jacobious52/StatsBot2/pkg/storage"
	tb "gopkg.in/tucnak/telebot.v2"
)

type Dump struct {
	DBPath string
}

func (h *Dump) Do(data storage.Model, info storage.MessageInfo) (interface{}, error) {
	w := bytes.NewBufferString("")
	err := json.NewEncoder(w).Encode(&data)
	if err != nil {
		return "failed to dump database", err
	}
	return tb.Document{File: tb.FromDisk(h.DBPath), FileName: "db.json"}, nil
}
