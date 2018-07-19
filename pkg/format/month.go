package format

import (
	"github.com/Jacobious52/StatsBot2/pkg/storage"
)

func Monthly(messages storage.Messages) Table {
	table := make(Table)
	for time := range messages {
		table[int(time.Month())]++
	}
	return table
}
