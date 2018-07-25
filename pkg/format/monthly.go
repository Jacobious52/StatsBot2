package format

import (
	"time"

	"github.com/Jacobious52/StatsBot2/pkg/storage"
)

func Monthly(messages storage.Messages, filter storage.MessageKeyFilter, location *time.Location) Table {
	table := make(Table)
	for time, tag := range messages {
		if tag.Matches(filter) {
			table[int(time.In(location).Month())]++
		}
	}
	return table
}
