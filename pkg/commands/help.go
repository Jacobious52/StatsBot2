package commands

import (
	"bytes"
	"strings"

	"github.com/Jacobious52/StatsBot2/pkg/storage"
)

type Help struct {
	Commander *CommandManager
}

func (h *Help) Do(storage.Model, storage.MessageInfo) (interface{}, error) {
	var buff bytes.Buffer
	buff.WriteString("[Registered Commands]\n")
	for command := range h.Commander.commands {
		if strings.HasPrefix(command, "/") {
			buff.WriteString(command)
			buff.WriteString("\n")
		}
	}
	buff.WriteString("\n")
	return buff.String(), nil
}
