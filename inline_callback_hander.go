package openpoll

import (
	"github.com/Sirupsen/logrus"
	"github.com/bot-api/telegram"
	"github.com/bot-api/telegram/telebot"
	"github.com/davecgh/go-spew/spew"
	"golang.org/x/net/context"
	"strings"
)

func HandleInlineCallback(ctx context.Context, cb *telegram.CallbackQuery) (err error) {
	api := telebot.GetAPI(ctx)

	//redis, err := GetRedisConnection(); if err != nil {
	//	logrus.WithError(err).Error("Redis is unavailable, skipping request")
	//	return err
	//}

	cbDataParts := strings.Split(cb.Data, ":")

	switch {
	case len(cbDataParts) != 3:
		logrus.Warnf("Too many arguments in cb data: %v", cbDataParts)
		return nil
	case cbDataParts[0] != "v":
		logrus.Warnf("Unsupported action: %v", cbDataParts[0])
		return nil
	}

	//pollId, err := strconv.ParseInt(cbDataParts[1], 10, 64); if err != nil {
	//	logrus.Warn("Cannot parse poll id: ", cbDataParts[1])
	//	return nil
	//}

	//voteResult := cbDataParts[2]

	cfg := telegram.AnswerCallbackCfg{
		CallbackQueryID: cb.ID,
		Text:            "Vote confirmed",
		ShowAlert:       false,
	}

	_, err = api.AnswerCallbackQuery(ctx, cfg)

	if err != nil {
		logrus.Info("Failed to send inline query answer:", spew.Sdump(err))
	}

	return err
}
