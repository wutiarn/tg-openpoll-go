package openpoll

import (
	"golang.org/x/net/context"
	"github.com/bot-api/telegram"
	"github.com/bot-api/telegram/telebot"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"gopkg.in/kyokomi/emoji.v1"
	"github.com/Sirupsen/logrus"
)

func HandleInlineQuery(ctx context.Context, query *telegram.InlineQuery) (err error) {
	voteId := query.ID

	api := telebot.GetAPI(ctx)
	cfg := telegram.AnswerInlineQueryCfg{
		InlineQueryID: query.ID,
		CacheTime: 0,
		IsPersonal: true,
		Results: []telegram.InlineQueryResult{
			telegram.InlineQueryResultArticle{
				BaseInlineQueryResult: telegram.BaseInlineQueryResult{
					Type: "article",
					ID:   "id",
					InputMessageContent: telegram.InputTextMessageContent{
						MessageText: fmt.Sprintf("*%v*", query.Query),
						ParseMode: "Markdown",
					},
					ReplyMarkup: telegram.InlineKeyboardMarkup{
						InlineKeyboard: [][]telegram.InlineKeyboardButton{
							[]telegram.InlineKeyboardButton{
								{
									CallbackData: fmt.Sprintf("v:%v:u", voteId),
									Text: emoji.Sprint(":thumbsup:"),
								}, {
									CallbackData: fmt.Sprintf("v:%v:n", voteId),
									Text: emoji.Sprint(":neutral_face:"),
								},
								{
									CallbackData: fmt.Sprintf("v:%v:d", voteId),
									Text: emoji.Sprint(":thumbsdown:"),
								},
							}, []telegram.InlineKeyboardButton{
								{
									CallbackData: fmt.Sprintf("v:%v:r", voteId),
									Text: "Results",
								},
							},
						},
					},
				},
				Title: fmt.Sprintf("[%v] Start vote: %v", voteId, query.Query),
			},
		},
	}

	_, err = api.AnswerInlineQuery(ctx, cfg)
	if err != nil {
		logrus.Info("Failed to send inline query answer:", spew.Sdump(err))
	}
	return err
}