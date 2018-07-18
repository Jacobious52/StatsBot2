package commands

import (
	"github.com/Jacobious52/StatsBot2/pkg/storage"
)

type OnSticker struct{}

func (h *OnSticker) Do(data storage.Model, info storage.MessageInfo) (interface{}, error) {
	if _, ok := data[info.Chat]; !ok {
		data[info.Chat] = make(storage.Chat)
	}
	if _, ok := data[info.Chat][info.Sender]; !ok {
		data[info.Chat][info.Sender] = make(storage.Messages)
	}
	data[info.Chat][info.Sender][info.Timestamp] = "sticker"
	return nil, nil
}
