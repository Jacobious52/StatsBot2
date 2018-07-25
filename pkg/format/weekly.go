package format

import (
	"time"

	"github.com/Jacobious52/StatsBot2/pkg/storage"
)

func Weekly(messages storage.Messages, filter storage.MessageKeyFilter, location *time.Location) Table {
	table := make(Table)
	for time, tag := range messages {
		if tag.Matches(filter) {
			_, week := time.In(location).ISOWeek()
			table[week]++
		}
	}
	return table
}
