package commands

import "github.com/Jacobious52/StatsBot2/pkg/storage"

type Start struct{}

func (h *Start) Do(storage.Model, storage.MessageInfo) (interface{}, error) {
	return "start me?", nil
}
