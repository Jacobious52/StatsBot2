package format

import (
	"github.com/Jacobious52/StatsBot2/pkg/storage"
)

func Dayly(messages storage.Messages) Table {
	table := make(Table)
	for time := range messages {
		table[time.YearDay()]++
	}
	return table
}
