package format

import (
	"github.com/Jacobious52/StatsBot2/pkg/storage"
)

func Weekly(messages storage.Messages, filter storage.MessageKeyFilter) Table {
	table := make(Table)
	for time, tag := range messages {
		if tag.Matches(filter) {
			_, week := time.ISOWeek()
			table[week]++
		}
	}
	return table
}
