package commands

import "github.com/Jacobious52/StatsBot2/pkg/storage"

type Help struct{}

func (h *Help) Do(storage.Model, storage.MessageInfo) (interface{}, error) {
	return "help me!", nil
}
