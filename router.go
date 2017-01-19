package openpoll

import (
	"github.com/bot-api/telegram"
	"github.com/bot-api/telegram/telebot"
	"golang.org/x/net/context"
	"github.com/Sirupsen/logrus"
	"github.com/wutiarn/tg-openvote/cmdhandlers"
	"github.com/davecgh/go-spew/spew"
)


func Run(token string, debug bool) {
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
		FullTimestamp:true,
	})

	api := telegram.New(token)
	api.Debug(debug)

	bot := telebot.NewWithAPI(api)

	bot.HandleFunc(defaultHandler)

	bot.Use(telebot.Commands(map[string]telebot.Commander{
		"start": commands.StartCommand{},
	}))

	netCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logrus.Info("Started")

	bot.Serve(netCtx)
}

func defaultHandler(ctx context.Context) error {
	upd := telebot.GetUpdate(ctx)
	api := telebot.GetAPI(ctx)
	msg := upd.Message
	from := upd.From()

	switch {
	case upd.InlineQuery != nil:
		query := upd.InlineQuery
		logrus.Info("Inline query from", from, ":", query.Query)
		HandleInlineQuery(ctx, query)
	case upd.CallbackQuery != nil:
		cb := upd.CallbackQuery
		logrus.Infof("Callback query from %+v: %v", from, cb.Data)
		HandleInlineCallback(ctx, cb)
	case upd.Message != nil:
		logrus.Info("Msg from", from, ":", msg.Text)

		respMsg := telegram.NewMessagef(msg.Chat.ID, "Unsupported command")
		respMsg.ReplyToMessageID = msg.MessageID

		api.Send(ctx, respMsg)
	default:
		logrus.Info("Unsupported update", spew.Sdump(upd))
	}
	return nil
}
