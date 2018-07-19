package format

import (
	"encoding/csv"
	"fmt"
	"io"

	"github.com/Jacobious52/StatsBot2/pkg/storage"
	log "github.com/sirupsen/logrus"
)

type Table map[int]int

type Formatter func(storage.Messages) Table

func CSV(chat storage.Chat, formatter Formatter, w io.Writer) {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	header := []string{"time"}
	tables := make(map[storage.UserKey]Table)
	for user, messages := range chat {
		tables[user] = formatter(messages)
		header = append(header, string(user))
	}

	writer.Write(header)

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

	for i := start; i < end+1; i++ {
		row := make([]string, len(tables)+1)
		row[0] = fmt.Sprint(i)
		for user, table := range tables {
			j := indexOf(user, header)
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
