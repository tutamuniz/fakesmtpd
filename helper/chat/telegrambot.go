package chat

import (
	"log"
	"os"
	"strconv"
	"time"

	tele "gopkg.in/telebot.v3"
)

type TelegramBot struct {
	bot       *tele.Bot
	channelID *tele.Chat
	message   chan string
	Logger    *log.Logger
}

func NewBot(logger *log.Logger) *TelegramBot {
	pref := tele.Settings{
		Token:  os.Getenv("API_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		logger.Fatal(err)
		return nil
	}

	id, _ := strconv.ParseInt(os.Getenv("CHANNEL_ID"), 10, 64)
	channelID, err := b.ChatByID(id)
	if err != nil {
		logger.Fatal(err)
		return nil
	}

	return &TelegramBot{
		bot:       b,
		channelID: channelID,
		message:   make(chan string),
		Logger:    logger,
	}
}

func (tb *TelegramBot) ProcessMessages() {
	go func(m chan string) {
		for {
			select {
			case filename := <-m:
				tb.Logger.Printf("sending message %s\n", filename)

				photo := &tele.Photo{File: tele.FromDisk(filename)} //.FromURL("https://www.mapadomeuceu.com.br/wp-content/uploads/2020/09/Gal%C3%A1xia-de-Andr%C3%B4meda.jpg")}
				_, err := tb.bot.Send(tb.channelID, photo)
				if err != nil {
					tb.Logger.Println(err)
				}
			case <-time.After(time.Second * 5):
			}
		}
	}(tb.message)

	tb.bot.Start()
}

func (tb *TelegramBot) SendMessage(user, msg string) {
	tb.message <- msg
}
