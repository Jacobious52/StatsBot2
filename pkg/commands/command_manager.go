package commands

import (
	"github.com/Jacobious52/StatsBot2/pkg/command"
	"github.com/Jacobious52/StatsBot2/pkg/storage"
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
)

// CommandManager manages all the commands with the datastore and bot
type CommandManager struct {
	// Commands are all the exposers plugged in
	commands map[string]command.Command
	bot      *tb.Bot
	store    *storage.DataStore
}

// NewCommandManager creates a new command manager
func NewCommandManager(bot *tb.Bot, store *storage.DataStore) *CommandManager {
	return &CommandManager{
		commands: make(map[string]command.Command),
		bot:      bot,
		store:    store,
	}
}

// RegisterCommand adds the command to the commands list
func (p *CommandManager) RegisterCommand(name string, c command.Command) {
	p.commands[name] = c

	// create the endpoint
	p.bot.Handle(name, func(m *tb.Message) {
		logger := log.WithFields(log.Fields{
			"command": name,
			"chat":    m.Chat.ID,
			"sender":  m.Sender.ID,
			"time":    m.Time(),
		})

		logger.Debugln("locking store")
		p.store.Lock()
		response, err := c.Do(p.store.Data, storage.MessageInfo{
			Chat:      storage.ChatKey(m.Chat.ID),
			Sender:    storage.UserKey(m.Sender.Username),
			Timestamp: m.Time(),
		})
		p.store.Dirty = true
		logger.Debugln("unlocking store")
		p.store.Unlock()
		if err != nil {
			logger.Errorf("failed to run command: %v", name, err)
			return
		}
		if response != nil {
			logger.Infoln("sending response")
			p.bot.Send(m.Chat, response)
		} else {
			logger.Infoln("no response needed")
		}
	})
}
