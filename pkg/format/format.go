package format

import (
	"encoding/csv"
	"fmt"
	"io"

	"github.com/Jacobious52/StatsBot2/pkg/storage"
	log "github.com/sirupsen/logrus"
)

type Table map[int]int

type Formatter func(storage.Messages, storage.MessageKeyFilter) Table

func CSV(chat storage.Chat, timeFrame string, formatter Formatter, filter storage.MessageKeyFilter, w io.Writer) {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	// create the header of names
	header := []string{timeFrame}
	tables := make(map[storage.UserKey]Table)
	for user, messages := range chat {
		tables[user] = formatter(messages, filter)
		header = append(header, string(user))
	}
	writer.Write(header)

	// get the min and max ranges recorded
	start := (1 << 31) - 1
	end := 0
	for _, table := range tables {
		for t := range table {
			if t < start {
				start = t
			}
			if t > end {
				end = t
			}
		}
	}

	log.Debugf("start=%d end=%d len(tables)=%d", start, end, len(tables))

	// write the rows in the correct ordering of the names in headers
	for i := start; i < end+1; i++ {
		row := make([]string, len(tables)+1)
		row[0] = fmt.Sprint(i)
		for user, table := range tables {
			j := indexOf(user, header)
			if j == -1 {
				continue
			}
			row[j] = fmt.Sprint(table[i])
		}
		writer.Write(row)
	}
}

func indexOf(user storage.UserKey, header []string) int {
	for i, v := range header {
		if string(user) == v {
			return i
		}
	}
	return -1
}
