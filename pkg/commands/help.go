package commands

import "github.com/Jacobious52/StatsBot2/pkg/storage"

type Help struct{}

func (h *Help) Do(storage.Model, storage.MessageInfo) (interface{}, error) {
	return "yay statsbot2.. here are the new commands yea\n/month\n/day\n/dump\nyou just get csv or json.. cbf doing graphs. yay the end", nil
}
