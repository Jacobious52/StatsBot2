package command

import (
	"github.com/Jacobious52/StatsBot2/pkg/storage"
)

// Command is an interface for something that can expose metrics from a datastore
type Command interface {
	// Do will run every call to the registered plugin
	Do(storage.Model, storage.MessageInfo) (interface{}, error)
}
